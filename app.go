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
	"os"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/entry"
	"github.com/terasum/medict/internal/static/handler"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/backserver"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
	"github.com/terasum/medict/pkg/service/ankiexport"
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
	conf         *config.Config // app config (medict.toml);用于偏好回写
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
	b.conf = conf

	dictsSvc, err := service.NewDictService(conf)
	if err != nil {
		return err
	}
	b.dictSvc = dictsSvc

	// Bookmark store (#643): persisted in app config dir (SQLite).
	configDir, err := utils.AppConfigDir()
	if err != nil {
		return err
	}
	b.bookmarks, err = service.NewBookmarkStore(configDir)
	if err != nil {
		return err
	}

	// #783: load per-dictionary CSS overrides into the handler's in-memory map.
	loadUserCSSOverrides()

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
	b.ctx = ctx
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
	// Close the bookmark SQLite store.
	if b.bookmarks != nil {
		if err := b.bookmarks.Close(); err != nil {
			log.Errorf("shutdown: close bookmark store failed: %s", err.Error())
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

// GetPreferences returns all persisted settings (medict.toml + runtime overrides)
// as a flat map, for the frontend to read back. b.conf may be nil if appInit
// failed — return an empty map in that case.
//
// Note: viper normalizes keys to lowercase, so returned keys are lowercase
// regardless of the casing used when saving.
func (b *App) GetPreferences() *model.Resp {
	if b.conf == nil {
		return model.BuildSuccess(map[string]any{})
	}
	return model.BuildSuccess(b.conf.Preferences())
}

// SavePreferences merges the given key/value pairs into medict.toml and writes
// it back to disk; existing keys are preserved. Generic write-back path for user
// preferences (multi-dict active set, theme, font size, …). Keys are
// case-insensitive (viper lowercases them).
func (b *App) SavePreferences(prefs map[string]any) *model.Resp {
	if b.conf == nil {
		return model.BuildError(errors.New("config not initialized"), model.InnerSysErrCode)
	}
	if err := b.conf.WritePreferences(prefs); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(nil)
}

// Bookmark / notebook IPC (#643)
func (b *App) AddBookmark(word, dictId, notebookId string) *model.Resp {
	dictName := ""
	if d, ok := b.dictSvc.GetDictPlain(dictId); ok {
		dictName = d.Name
	}
	// Render a self-contained HTML snapshot so the word is still viewable if the
	// dictionary is later unloaded. Snapshot failure is non-fatal: the word is
	// saved without a snapshot and falls back to a live lookup.
	html := ""
	if h, err := b.bs.Controller.RenderSnapshot(dictId, word); err != nil {
		log.Errorf("AddBookmark: snapshot %q failed: %s", word, err)
	} else {
		html = h
	}
	if err := b.bookmarks.Add(word, dictId, dictName, notebookId, html); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(nil)
}

func (b *App) RemoveBookmark(word, dictId, notebookId string) *model.Resp {
	if err := b.bookmarks.Remove(word, dictId, notebookId); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(nil)
}

func (b *App) GetBookmarks() *model.Resp {
	items, err := b.bookmarks.All()
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(items)
}

// GetBookmarkSnapshot returns the stored HTML snapshot for a bookmark ("" if none).
func (b *App) GetBookmarkSnapshot(word, dictId, notebookId string) *model.Resp {
	html, err := b.bookmarks.GetSnapshot(word, dictId, notebookId)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(html)
}

// Notebook management: list / create / rename / delete / set-default.
func (b *App) GetNotebooks() *model.Resp {
	items, err := b.bookmarks.Notebooks()
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(items)
}

func (b *App) CreateNotebook(name string) *model.Resp {
	nb, err := b.bookmarks.AddNotebook(name)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(nb)
}

func (b *App) RenameNotebook(id, name string) *model.Resp {
	if err := b.bookmarks.RenameNotebook(id, name); err != nil {
		return model.BuildError(err, model.BadParamErrCode)
	}
	return model.BuildSuccess(nil)
}

func (b *App) DeleteNotebook(id string) *model.Resp {
	if err := b.bookmarks.RemoveNotebook(id); err != nil {
		return model.BuildError(err, model.BadParamErrCode)
	}
	return model.BuildSuccess(nil)
}

func (b *App) SetDefaultNotebook(id string) *model.Resp {
	if err := b.bookmarks.SetDefaultNotebook(id); err != nil {
		return model.BuildError(err, model.BadParamErrCode)
	}
	return model.BuildSuccess(nil)
}

// ExportAnki exports the given notebook's saved words (all notebooks if id is
// empty) to a native Anki .apkg. Each word becomes a card carrying its stored
// HTML snapshot (images extracted as Anki media); each notebook becomes a deck.
// Prompts the user for a save path via a native dialog.
func (b *App) ExportAnki(notebookId string) *model.Resp {
	rows, err := b.bookmarks.ExportRows(notebookId)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	if len(rows) == 0 {
		return model.BuildError(errors.New("没有可导出的生词"), model.BadParamErrCode)
	}

	path, err := runtime.SaveFileDialog(b.ctx, runtime.SaveDialogOptions{
		Title:           "导出 Anki 包",
		DefaultFilename: "medict.apkg",
		Filters: []runtime.FileFilter{
			{DisplayName: "Anki Package (*.apkg)", Pattern: "*.apkg"},
		},
	})
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	if path == "" {
		return model.BuildSuccess(nil) // user cancelled the dialog
	}

	src := make([]ankiexport.ExportRow, len(rows))
	for i, r := range rows {
		src[i] = ankiexport.ExportRow{
			Word:         r.Word,
			DictName:     r.DictName,
			NotebookName: r.NotebookName,
			HTML:         r.HTML,
		}
	}
	apkg, err := ankiexport.Build(src)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	if err := os.WriteFile(path, apkg, 0o644); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(path)
}

// ExportCurrentEntry renders the current word's definition as a self-contained
// HTML (all resources inlined as data: URLs) and writes it to a user-chosen
// file — for debugging complex entries (#784). Reuses RenderSnapshot (the same
// Search→Locate→@@LINK→WrapContent→InlineResources path used by bookmark
// snapshots). "" html means no match.
func (b *App) ExportCurrentEntry(dictId, word string) *model.Resp {
	if dictId == "" || word == "" {
		return model.BuildError(errors.New("dictId 和 word 不能为空"), model.BadParamErrCode)
	}
	html, err := b.bs.Controller.RenderSnapshot(dictId, word)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	if html == "" {
		return model.BuildError(errors.New("未找到该词条"), model.BadParamErrCode)
	}
	path, err := runtime.SaveFileDialog(b.ctx, runtime.SaveDialogOptions{
		Title:           "导出词条 HTML",
		DefaultFilename: word + ".html",
		Filters: []runtime.FileFilter{
			{DisplayName: "HTML (*.html)", Pattern: "*.html"},
		},
	})
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	if path == "" {
		return model.BuildSuccess(nil) // user cancelled the dialog
	}
	if err := os.WriteFile(path, []byte(html), 0o644); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(path)
}

// userCSSDir returns the directory for per-dictionary CSS overrides (under the
// app config dir, which is always writable — unlike the dict directory itself).
func userCSSDir() (string, error) {
	dir, err := utils.AppConfigDir()
	if err != nil {
		return "", err
	}
	cssDir := filepath.Join(dir, "user_css")
	os.MkdirAll(cssDir, 0755)
	return cssDir, nil
}

// loadUserCSSOverrides loads all per-dictionary CSS overrides from the app
// config dir into the handler's in-memory map at startup (#783).
func loadUserCSSOverrides() {
	cssDir, err := userCSSDir()
	if err != nil {
		return
	}
	entries, err := os.ReadDir(cssDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".css") {
			continue
		}
		dictId := strings.TrimSuffix(name, ".css")
		data, err := os.ReadFile(filepath.Join(cssDir, name))
		if err == nil {
			handler.SetUserCSS(dictId, string(data))
		}
	}
}

// GetDictUserCSS reads the per-dictionary user CSS override (#783).
func (b *App) GetDictUserCSS(dictId string) *model.Resp {
	cssDir, err := userCSSDir()
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	path := filepath.Join(cssDir, dictId+".css")
	if !utils.FileExists(path) {
		return model.BuildSuccess("")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(string(data))
}

// SaveDictUserCSS persists the per-dictionary user CSS override + updates the
// in-memory map so WrapContent injects it immediately (#783).
func (b *App) SaveDictUserCSS(dictId, css string) *model.Resp {
	cssDir, err := userCSSDir()
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	path := filepath.Join(cssDir, dictId+".css")
	if css == "" {
		os.Remove(path)
		handler.SetUserCSS(dictId, "")
		return model.BuildSuccess(nil)
	}
	if err := os.WriteFile(path, []byte(css), 0644); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	handler.SetUserCSS(dictId, css)
	return model.BuildSuccess(nil)
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
