package main

import (
	"app/core/doctor"
	"app/core/events"
	"app/core/processing"
	"app/core/settings"
	"context"
	"fmt"
	"image"

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
