package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// Needed to handle loading fontawesome fonts and image paths for thumbnails. POG!
func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := req.URL.Path
	println("[FileLoader] ", requestedFilename)
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("[FileLoader:Error] Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:           "Astropath",
		Width:           1024,
		Height:          768,
		DisableResize:   true,
		Fullscreen:      false,
		CSSDragProperty: "--wails-drop-target",
		CSSDragValue:    "drop",
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app, // Ensure the App instance is included here
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
