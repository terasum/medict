//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"embed"
	"os"

	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

// Version 是应用版本号，CI 通过 ldflags 注入：
//
//	wails build -ldflags "-X main.Version=<tag>"
//
// 之前 main 包没有该变量，链接器静默丢弃注入，版本号始终为空。
// 本地构建未注入时显示 "dev"。
var Version = "dev"

const cssEditorWindowFlag = "--medict-css-editor-window"

func parseCSSWindowArgs(args []string) (dictID, dictName, seedPath string, ok bool) {
	if len(args) < 4 || args[0] != cssEditorWindowFlag || !safeDictID.MatchString(args[1]) {
		return "", "", "", false
	}
	return args[1], args[2], args[3], true
}

func main() {
	if dictID, dictName, seedPath, ok := parseCSSWindowArgs(os.Args[1:]); ok {
		runCSSWindow(dictID, dictName, seedPath)
		return
	}
	// Create an instance of the app structure
	app := NewApp()

	appOptions := newAppOptions(app, false)
	appOptions.OnStartup = app.startup
	appOptions.OnDomReady = app.domReady
	appOptions.OnShutdown = app.shutdown

	mainApp := application.NewWithOptions(appOptions)
	if err := mainApp.Run(); err != nil {
		panic(err)
	}
}

func runCSSWindow(dictID, dictName, seedPath string) {
	app := newCSSWindowApp(dictID, dictName, seedPath)
	appOptions := newAppOptions(app, true)
	appOptions.OnStartup = func(ctx context.Context) { app.ctx = ctx }
	appOptions.OnBeforeClose = app.requestCSSWindowClose
	mainApp := application.NewWithOptions(appOptions)
	if err := mainApp.Run(); err != nil {
		panic(err)
	}
}

func newAppOptions(app *App, cssEditor bool) *options.App {
	title, width, height := "Medict", 800, 600
	minWidth, minHeight := 720, 570
	maxWidth, maxHeight := 1280, 740
	hideOnClose := true
	background := &options.RGBA{R: 33, G: 37, B: 43, A: 255}
	if cssEditor {
		title, width, height = "词典 CSS 编辑器", 760, 560
		minWidth, minHeight = 600, 420
		maxWidth, maxHeight = 1920, 1200
		hideOnClose = false
		background = &options.RGBA{R: 248, G: 249, B: 250, A: 255}
	}
	return &options.App{
		Title:             title,
		Width:             width,
		Height:            height,
		MinWidth:          minWidth,
		MinHeight:         minHeight,
		MaxWidth:          maxWidth,
		MaxHeight:         maxHeight,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: hideOnClose,
		BackgroundColour:  background,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Logger:             logger.NewDefaultLogger(),
		LogLevel:           logger.INFO,
		LogLevelProduction: logger.INFO,
		CSSDragProperty:    "--wails-draggable",
		CSSDragValue:       "drag",
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHidden(),
			WebviewIsTransparent: !cssEditor,
			WindowIsTranslucent:  !cssEditor,
			About: &mac.AboutInfo{
				Title:   "Medict APP",
				Message: "Light weight dictionary app for professional",
				Icon:    icon,
			},
		},
	}
}
