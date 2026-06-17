package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"mkd-epub-exporters/pkg/converter"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// Locate/Init external Pandoc engine
	_ = converter.InitPandoc()

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "MarkItDown Converter",
		Width:             900,
		Height:            560,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 13, G: 11, B: 24, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			BackdropType:                      windows.Mica,
			Theme:                             windows.Dark,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarDefault(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "MarkItDown Converter",
				Message: "High-performance document to Markdown & EPUB converter",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
