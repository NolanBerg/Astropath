package main

import (
	"app/core/doctor"
	"app/core/events"
	"app/core/processing"
	"app/core/settings"
	"app/core/workflow"
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"

	rt "runtime"

	"github.com/sanity-io/litter"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	eventManager *events.EventManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	eventManager := events.NewEventManager()
	return &App{
		eventManager: eventManager,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.eventManager.SetContext(ctx)
}

func (a *App) GetDoctorResults() doctor.DoctorReport {
	doctorReport := doctor.GetDoctorReport(a.eventManager)
	litter.Dump(doctorReport)
	return doctorReport
}

func (a *App) LoadApplicationSettings() settings.AppSettings {
	appSettings, err := settings.LoadAppSettings()
	if err != nil {
		panic(err)
	}

	litter.Dump(appSettings)
	return appSettings
}

func (a *App) ResetAppSettings() settings.AppSettings {
	appSettings, err := settings.ResetAppSettings()
	if err != nil {
		panic(err)
	}

	return appSettings
}

type ImportResult struct {
	ErrorMessage string
	FilePaths    []string
	Bounds       image.Rectangle
}

func (a *App) ImportImages() ImportResult {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})

	if err != nil {
		return ImportResult{
			ErrorMessage: err.Error(),
			FilePaths:    nil,
		}
	}

	if len(files) < 2 {
		return ImportResult{
			ErrorMessage: fmt.Sprintf("You must select two images with the same dimensions. Selected: %d", len(files)),
			FilePaths:    nil,
		}
	}

	bounds, paths, err := processing.ValidateFrames(files, a.eventManager)
	if err != nil {
		r := ImportResult{
			ErrorMessage: err.Error(),
			FilePaths:    nil,
			Bounds:       bounds,
		}

		litter.Dump(r)
		return r
	}

	r := ImportResult{
		ErrorMessage: "",
		FilePaths:    paths,
	}

	litter.Dump(r)
	return r
}

// UpdateTimelapseSetting updates the timelapse generation setting
func (a *App) UpdateTimelapseSetting(enabled bool) error {
	appSettings, err := settings.LoadAppSettings()
	if err != nil {
		return err
	}

	// Validate FFMPEG installation if enabling timelapse
	if enabled {
		doctorReport := doctor.GetDoctorReport(a.eventManager)
		if !doctorReport.SystemHasFFMPEG {
			return fmt.Errorf("FFMPEG is required for timelapse generation")
		}
	}

	appSettings.EnableTimelapseGeneration = enabled
	return settings.StoreAppSettings(appSettings)
}

// UpdateARWSetting updates the ARW conversion setting
func (a *App) UpdateARWSetting(enabled bool) error {
	appSettings, err := settings.LoadAppSettings()
	if err != nil {
		return err
	}

	// Validate dcraw installation if enabling ARW conversion
	if enabled {
		doctorReport := doctor.GetDoctorReport(a.eventManager)
		if !doctorReport.SystemHasARWConversion {
			return fmt.Errorf("dcraw is required for ARW conversion")
		}
	}

	appSettings.EnableARWConversion = enabled
	return settings.StoreAppSettings(appSettings)
}

// SelectARWDirectory opens a directory dialog for selecting ARW conversion output location
func (a *App) SelectARWDirectory() (string, error) {
	dialogOptions := runtime.OpenDialogOptions{
		Title: "Select ARW Conversion Directory",
	}

	path, err := runtime.OpenDirectoryDialog(a.ctx, dialogOptions)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist")
	}

	testFile := path + "/.astropath_test"
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return "", fmt.Errorf("directory is not writable")
	}
	os.Remove(testFile)

	return path, nil
}

// UpdateARWTempFilePath updates the ARW conversion temporary file path in the AppSettings
func (a *App) UpdateARWTempFilePath(newPath string) error {
	appSettings, err := settings.LoadAppSettings()
	if err != nil {
		return err
	}

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist")
	}

	appSettings.ARWTempFilePath = newPath
	return settings.StoreAppSettings(appSettings)
}

// UpdateSettings updates all application settings at once
func (a *App) UpdateSettings(s settings.AppSettings) error {
	if s.EnableARWConversion {
		if s.ARWTempFilePath == "" {
			return fmt.Errorf("output directory must be selected for ARW conversion")
		}

		if os.MkdirAll(s.ARWTempFilePath, 0777) != nil {
			return fmt.Errorf("selected directory does not exist")
		}
	}

	return settings.StoreAppSettings(s)
}

func (a *App) StartProcessingWorkflow(workflow workflow.Workflow, batch processing.ImageBatch) {
	fmt.Printf("Starting processing %v\n", batch)
	processor := processing.CreateWorkflowProcessor(&batch, a.eventManager, &workflow)
	processor.Process()
}

func (a *App) ShowPathInFinder(path string) bool {
	os := rt.GOOS
	var os_cmd string
	switch os {
	case "windows":
		os_cmd = "start"
	case "darwin":
		os_cmd = "open"
	case "linux":
		os_cmd = "xdg-open"
	default:
		return false
	}

	cmd := exec.Command(os_cmd, filepath.Dir(path))
	return cmd.Run() != nil
}
