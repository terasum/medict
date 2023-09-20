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
	"errors"
	"github.com/op/go-logging"
	"github.com/skratchdot/open-golang/open"
	"github.com/terasum/medict/internal/entry"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/backserver"
	"github.com/terasum/medict/pkg/model"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
)

var log = logging.MustGetLogger("default")

// App struct
type App struct {
	ctx context.Context

	errorChannel chan error
	stopChannel  chan int
	bs           *backserver.BackServer
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		errorChannel: make(chan error),
		stopChannel:  make(chan int),
		bs:           &backserver.BackServer{Ready: false},
	}

	err := app.appInit()
	if err != nil {
		go func() {
			app.errorChannel <- err
		}()
	}
	return app
}

func (b *App) appInit() error {
	conf, err := entry.DefaultConfig()
	if err != nil {
		return err
	}

	bs, err := backserver.NewStaticServer(conf)
	if err != nil {
		return err
	}

	err = bs.SetUp()
	if err != nil {
		return err
	}
	b.bs = bs
	// running bs, this is not blocking
	bs.Start()
	return nil
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	go b.stopChanListen(ctx)
	go b.errorChanListen(ctx)
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
	b.bs.GracefulStop()
}

func (b *App) Dispatch(apiName string, args map[string]interface{}) *model.Resp {
	log.Infof("[wails] IPC request dispatch [%s] | args: %v\n", apiName, args)
	return b.bs.DispatchIPCReq(apiName, args)
}

func (b *App) ResourceServerAddr() string {
	return b.bs.StaticServerBaseUrl()
}

func (b *App) OpenFinder(filepath string) error {
	if !fileutil.Exist(filepath) {
		return errors.New("file path not exist, cannot open")
	}
	err := open.Run(filepath)
	if err != nil {
		return err
	}
	return nil
}

func (b *App) BaseDictDir() string {
	if b.bs == nil {
		return "internal error"
	}
	f, e := utils.ReplaceHome(b.bs.Config.BaseDictDir)
	if e != nil {
		return "internal error"
	}
	return f
}
