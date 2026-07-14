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
	"github.com/terasum/medict/pkg/service"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
)

var log = logging.MustGetLogger("app")

// App struct
type App struct {
	ctx context.Context

	errorChannel chan error
	stopChannel  chan int
	bs           *backserver.BackServer
	dictSvc      *service.DictService
	bookmarks    *service.BookmarkStore
	// initErr 捕获 appInit 同步阶段的错误，由 errorChanListen 在 Wails
	// startup() 生命周期里确定性地产出，避免向无缓冲 channel 塞值带来的时序赌博。
	initErr error
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		errorChannel: make(chan error),
		stopChannel:  make(chan int),
		bs:           &backserver.BackServer{Ready: false},
	}

	app.initErr = app.appInit()
	return app
}

func (b *App) appInit() error {
	conf, err := entry.LoadApp()
	if err != nil {
		return err
	}

	dictsSvc, err := service.NewDictService(conf)
	if err != nil {
		return err
	}
	b.dictSvc = dictsSvc

	// Bookmark store (#643): persisted in app config dir.
	configDir, err := utils.AppConfigDir()
	if err != nil {
		return err
	}
	b.bookmarks = service.NewBookmarkStore(configDir)

	bs, err := backserver.NewStaticServer(conf)
	if err != nil {
		return err
	}

	if err := bs.SetUp(dictsSvc); err != nil {
		return err
	}
	// assign backend server
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
	// Release per-dictionary resources (leveldb handles, etc.).
	if b.dictSvc != nil {
		if err := b.dictSvc.Close(); err != nil {
			log.Errorf("shutdown: close dictionaries failed: %s", err.Error())
		}
	}
}

// Typed IPC handlers (issue #729): each frontend call maps to a typed App
// method, replacing the old string-dispatched Dispatch / handlerMap.
func (b *App) InitDicts() *model.Resp { return b.bs.Controller.InitDicts() }

func (b *App) GetAllDicts() *model.Resp { return b.bs.Controller.GetAllDicts() }

func (b *App) SearchWord(dictId, word string) *model.Resp {
	return b.bs.Controller.SearchWord(dictId, word)
}

func (b *App) BuildIndexByDictId(dictid string) *model.Resp {
	return b.bs.Controller.BuildIndexByDictId(dictid)
}

// Bookmark IPC (#643)
func (b *App) AddBookmark(word, dictId, dictName string) *model.Resp {
	b.bookmarks.Add(word, dictId, dictName)
	return model.BuildSuccess(nil)
}

func (b *App) RemoveBookmark(word, dictId string) *model.Resp {
	b.bookmarks.Remove(word, dictId)
	return model.BuildSuccess(nil)
}

func (b *App) GetBookmarks() *model.Resp {
	return model.BuildSuccess(b.bookmarks.All())
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
