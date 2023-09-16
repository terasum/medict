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
	"fmt"
	"github.com/terasum/medict/internal/entry"
	"github.com/terasum/medict/pkg/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

// App struct
type App struct {
	ctx context.Context

	sserver *service.StaticService
	ready   bool

	errorChannel chan error
	stopChannel  chan int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	errorChannel := make(chan error)
	stopChannel := make(chan int)
	b.errorChannel = errorChannel
	b.stopChannel = stopChannel

	go b.startErrorSignalListen()
	go b.startStopSignalListen()

	config, err := entry.DefaultConfig()
	if err != nil {
		panicWithErrorMessageDialog(ctx, err)
		os.Exit(-1)
		return
	}

	dictService, err := service.NewDictService(config)
	if err != nil {
		panicWithErrorMessageDialog(ctx, err)
		os.Exit(-1)
		return
	}

	staticService := service.NewStaticServer(ctx, config, dictService, errorChannel, stopChannel)

	// 开始服务
	go staticService.Start()

	b.sserver = staticService
	b.ready = true
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	close(b.stopChannel)
	close(b.errorChannel)
}

func (b *App) StaticServerURL() string {
	if b.ready {
		return b.sserver.StaticServerBaseUrl()
	} else {
		return ""
	}
}

func panicWithErrorMessageDialog(ctx context.Context, err error) {
	_, dialogErr := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Title:         "错误",
		Type:          runtime.ErrorDialog,
		Message:       err.Error(),
		Buttons:       []string{"OK"},
		DefaultButton: "OK",
	})
	if dialogErr != nil {
		fmt.Printf("[CRITIC] open dialog error %s\n", dialogErr.Error())
	}
}

func (b *App) startErrorSignalListen() {
	for {
		select {
		case err := <-b.errorChannel:
			panicWithErrorMessageDialog(b.ctx, err)
		case <-b.stopChannel:
			return
		}
	}
}

func (b *App) startStopSignalListen() {
	for {
		select {
		case <-b.stopChannel:
			return
		}
	}
}
