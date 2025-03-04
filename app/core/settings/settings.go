package settings

import (
	"app/core/workflow"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const APP_NAME = "Astropath"
const APP_SETTINGS_FILENAME = "appsettings.json"
const DEFAULT_OUTPUT_LOCATION = "/output"

type AppSettings struct {
	UserFirstTime             bool                // whether this is the user's first time experiencing the application or not. set to true after leaving the main-menu page for first time
	Workflows                 []workflow.Workflow // Slice of all Workflow's the user has saved.
	EnableTimelapseGeneration bool                // Whether timelapse generation support is allowed. If true, user MUST have ffmpeg
	EnableARWConversion       bool                // Whether arw conversion is enabled. If true the user must have dcraw installed
	ARWTempFilePath           string              // Where to store the files for conversion. When the user enables this in the application, the user MUST enter a valid location for this. It will take up alot of space.
}

func LoadAppSettings() (AppSettings, error) {
	var settings AppSettings
	appFolder, err := GetAppFolder()
	if err != nil {
		return settings, err
	}

	settingsPath := filepath.Join(appFolder, APP_SETTINGS_FILENAME)
	info, err := os.Stat(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return createAndStoreDefaultSettings()
		}

		return settings, err
	}

	if info.IsDir() {
		os.RemoveAll(settingsPath)
		return createAndStoreDefaultSettings()
	}

	bytes, err := os.ReadFile(settingsPath)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func createAndStoreDefaultSettings() (AppSettings, error) {
	appFolder, _ := GetAppFolder()

	defaultSettings := AppSettings{
		UserFirstTime:             true,
		EnableARWConversion:       false,
		ARWTempFilePath:           filepath.Join(appFolder, DEFAULT_OUTPUT_LOCATION, "temp_arws"),
		EnableTimelapseGeneration: false,

		Workflows: []workflow.Workflow{
			*workflow.CreateDefaultWorkflow(filepath.Join(appFolder, DEFAULT_OUTPUT_LOCATION)),
		},
	}

	if err := StoreAppSettings(defaultSettings); err != nil {
		return AppSettings{}, err
	}

	return defaultSettings, nil
}

func ResetAppSettings() (AppSettings, error) {
	return createAndStoreDefaultSettings()
}

func StoreAppSettings(settings AppSettings) error {
	appFolder, err := GetAppFolder()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(appFolder, APP_SETTINGS_FILENAME)
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}

// Returns the Settings folder. Need to check to make sure it works on all operating sytems.
// On Windows:  %APPDATA%\Astropath
// On macOS: ~/Library/Application Support/Astropath

func IsDevBuild() bool {
	return true
}

func GetAppFolder() (string, error) {
	if IsDevBuild() {
		return os.Getwd()
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not get user config directory: %w", err)
	}

	appFolder := filepath.Join(configDir, APP_NAME)

	_, err = os.Stat(appFolder)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(appFolder, 0755); err != nil {
			return "", fmt.Errorf("could not create app config directory: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("error checking app directory: %w", err)
	}

	testFile := filepath.Join(appFolder, ".test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return "", fmt.Errorf("cannot write to app directory: %w", err)
	}

	os.Remove(testFile)
	return appFolder, nil
}
