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
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	goruntime "runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/op/go-logging"
	"github.com/skratchdot/open-golang/open"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/entry"
	"github.com/terasum/medict/internal/static/handler"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/backserver"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
	"github.com/terasum/medict/pkg/service/ankiexport"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

	// cssWindowMode keeps the helper process deliberately lightweight: the
	// editor window can use the CSS persistence methods without loading any
	// dictionaries or opening their indexes a second time.
	cssWindowMode     bool
	cssWindowDictID   string
	cssWindowDictName string
	cssWindowSeedPath string
	cssWindowsMu      sync.Mutex
	cssWindows        map[string]*exec.Cmd
	cssSaveMu         sync.Mutex
	cssCloseMu        sync.Mutex
	cssWindowCanClose bool
}

const (
	cssEditorChangedEvent        = "medict:css-editor-changed"
	cssEditorCloseRequestedEvent = "medict:css-editor-close-requested"
)

var safeDictID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		errorChannel: make(chan error),
		stopChannel:  make(chan int),
		bs:           &backserver.BackServer{Ready: false},
		cssWindows:   make(map[string]*exec.Cmd),
	}

	app.initErr = app.appInit()
	return app
}

func newCSSWindowApp(dictID, dictName, seedPath string) *App {
	return &App{cssWindowMode: true, cssWindowDictID: dictID, cssWindowDictName: dictName, cssWindowSeedPath: seedPath}
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
	if !safeDictID.MatchString(dictId) {
		return model.BuildError(errors.New("invalid dictionary id"), model.BadParamErrCode)
	}
	css, err := readDictUserCSS(dictId)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(css)
}

// GetDictEditorCSS returns the saved user override when present. On first edit
// it falls back to a private snapshot of the dictionary's own stylesheet,
// resolved by the main process before the lightweight editor is launched.
func (b *App) GetDictEditorCSS(dictID string) *model.Resp {
	if !safeDictID.MatchString(dictID) || !b.cssWindowMode || dictID != b.cssWindowDictID {
		return model.BuildError(errors.New("invalid CSS editor context"), model.BadParamErrCode)
	}
	cssDir, err := userCSSDir()
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	css, err := readEditorCSS(filepath.Join(cssDir, dictID+".css"), b.cssWindowSeedPath)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(css)
}

func readEditorCSS(overridePath, seedPath string) (string, error) {
	for _, path := range []string{overridePath, seedPath} {
		if path == "" {
			continue
		}
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data), nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	return "", nil
}

// SaveDictUserCSS persists the per-dictionary user CSS override + updates the
// in-memory map so WrapContent injects it immediately (#783).
func (b *App) SaveDictUserCSS(dictId, css string) *model.Resp {
	if !safeDictID.MatchString(dictId) {
		return model.BuildError(errors.New("invalid dictionary id"), model.BadParamErrCode)
	}
	cssDir, err := userCSSDir()
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	path := filepath.Join(cssDir, dictId+".css")
	b.cssSaveMu.Lock()
	defer b.cssSaveMu.Unlock()
	if css == "" {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return model.BuildError(err, model.InnerSysErrCode)
		}
		handler.SetUserCSS(dictId, "")
		return model.BuildSuccess(nil)
	}
	if err := writeFileAtomic(path, []byte(css), 0644); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	handler.SetUserCSS(dictId, css)
	return model.BuildSuccess(nil)
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".medict-css-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// WindowMode lets the shared frontend select the compact editor surface in the
// helper process. The main application keeps the normal router-driven shell.
func (b *App) WindowMode() string {
	if b.cssWindowMode {
		return "css-editor"
	}
	return "main"
}

func (b *App) CSSWindowDictionary() map[string]string {
	return map[string]string{"id": b.cssWindowDictID, "name": b.cssWindowDictName}
}

func (b *App) requestCSSWindowClose(ctx context.Context) bool {
	b.cssCloseMu.Lock()
	defer b.cssCloseMu.Unlock()
	if b.cssWindowCanClose {
		return false
	}
	runtime.EventsEmit(ctx, cssEditorCloseRequestedEvent)
	return true
}

// CloseCSSWindow is called only after the editor has flushed pending content.
func (b *App) CloseCSSWindow() {
	b.cssCloseMu.Lock()
	b.cssWindowCanClose = true
	b.cssCloseMu.Unlock()
	runtime.Quit(b.ctx)
}

// OpenDictCSSWindow launches a second, CSS-only instance of this executable.
// Wails v2 has no public multi-window API; using the same signed executable
// gives every platform a genuine independent native window without loading
// dictionary indexes twice.
func (b *App) OpenDictCSSWindow(dictID, dictName, word string) *model.Resp {
	if !safeDictID.MatchString(dictID) {
		return model.BuildError(errors.New("invalid dictionary id"), model.BadParamErrCode)
	}
	b.cssWindowsMu.Lock()
	if cmd := b.cssWindows[dictID]; cmd != nil && cmd.Process != nil {
		b.cssWindowsMu.Unlock()
		return model.BuildSuccess(nil)
	}
	executable, err := os.Executable()
	if err != nil {
		b.cssWindowsMu.Unlock()
		return model.BuildError(err, model.InnerSysErrCode)
	}
	baseline, _ := readDictUserCSS(dictID)
	seedPath := ""
	cssDir, cssDirErr := userCSSDir()
	if cssDirErr != nil {
		b.cssWindowsMu.Unlock()
		return model.BuildError(cssDirErr, model.InnerSysErrCode)
	}
	if _, statErr := os.Stat(filepath.Join(cssDir, dictID+".css")); errors.Is(statErr, os.ErrNotExist) {
		if baseCSS := b.dictionaryBaseCSS(dictID, word); baseCSS != "" {
			seed, seedErr := os.CreateTemp("", "medict-editor-seed-*.css")
			if seedErr != nil {
				b.cssWindowsMu.Unlock()
				return model.BuildError(seedErr, model.InnerSysErrCode)
			}
			seedPath = seed.Name()
			if _, seedErr = seed.WriteString(baseCSS); seedErr == nil {
				seedErr = seed.Close()
			} else {
				_ = seed.Close()
			}
			if seedErr != nil {
				_ = os.Remove(seedPath)
				b.cssWindowsMu.Unlock()
				return model.BuildError(seedErr, model.InnerSysErrCode)
			}
		}
	} else if statErr != nil {
		b.cssWindowsMu.Unlock()
		return model.BuildError(statErr, model.InnerSysErrCode)
	}
	cmd := newCSSWindowCommand(executable, dictID, dictName, seedPath)
	// A wails dev child needs its own IPC/dev-server port, while it may keep
	// using the parent's Vite URL and asset directory inherited in the env.
	cmd.Env = replaceProcessEnv(os.Environ(), "devserver", "localhost:0")
	if err := cmd.Start(); err != nil {
		_ = os.Remove(seedPath)
		b.cssWindowsMu.Unlock()
		return model.BuildError(err, model.InnerSysErrCode)
	}
	b.cssWindows[dictID] = cmd
	b.cssWindowsMu.Unlock()
	go b.watchCSSWindow(dictID, baseline, seedPath, cmd)
	return model.BuildSuccess(nil)
}

func newCSSWindowCommand(executable, dictID, dictName, seedPath string) *exec.Cmd {
	// A packaged macOS GUI executable needs LaunchServices activation to create
	// a visible NSWindow when started from another GUI process. `open -W` keeps
	// the launcher alive so the parent can still observe window closure.
	if goruntime.GOOS == "darwin" && strings.Contains(executable, ".app/Contents/MacOS/") {
		bundle := filepath.Clean(filepath.Join(filepath.Dir(executable), "..", ".."))
		return exec.Command("/usr/bin/open", "-n", "-W", bundle, "--args", cssEditorWindowFlag, dictID, dictName, seedPath)
	}
	return exec.Command(executable, cssEditorWindowFlag, dictID, dictName, seedPath)
}

func replaceProcessEnv(env []string, key, value string) []string {
	prefix := key + "="
	result := make([]string, 0, len(env)+1)
	for _, item := range env {
		if !strings.HasPrefix(item, prefix) {
			result = append(result, item)
		}
	}
	return append(result, prefix+value)
}

func (b *App) watchCSSWindow(dictID, lastCSS, seedPath string, cmd *exec.Cmd) {
	defer os.Remove(seedPath)
	done := make(chan struct{})
	go func() { _ = cmd.Wait(); close(done) }()
	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			css, err := readDictUserCSS(dictID)
			if err == nil && css != lastCSS {
				lastCSS = css
				handler.SetUserCSS(dictID, css)
				runtime.EventsEmit(b.ctx, cssEditorChangedEvent, dictID, css)
			}
		case <-done:
			css, _ := readDictUserCSS(dictID)
			handler.SetUserCSS(dictID, css)
			runtime.EventsEmit(b.ctx, cssEditorChangedEvent, dictID, css)
			b.cssWindowsMu.Lock()
			delete(b.cssWindows, dictID)
			b.cssWindowsMu.Unlock()
			return
		}
	}
}

func (b *App) dictionaryBaseCSS(dictID, word string) string {
	if b.dictSvc == nil {
		return ""
	}
	dict := b.dictSvc.GetDictById(dictID)
	if dict == nil || dict.Dict == nil {
		return ""
	}
	// The current entry is the authoritative source for arbitrary and multiple
	// stylesheet names (including nested loose files and MDD resources). Reuse
	// the snapshot pipeline because it already resolves those references with
	// the same directory-first/resource-candidate behavior as runtime serving.
	if word = strings.TrimSpace(word); word != "" && b.bs != nil && b.bs.Controller != nil {
		if snapshot, err := b.bs.Controller.RenderSnapshot(dictID, word); err == nil {
			if css := extractInlinedCSS(snapshot); css != "" {
				return css
			}
		} else {
			log.Warningf("discover entry CSS for %s/%s failed: %s", dictID, word, err)
		}
	}
	names := make([]string, 0, 2)
	if dict.PathInfo != nil && dict.PathInfo.MdictMdxFileName != "" {
		names = append(names, dict.PathInfo.MdictMdxFileName+".css")
	}
	if name := strings.TrimSpace(dict.Name); name != "" {
		candidate := name + ".css"
		if len(names) == 0 || names[0] != candidate {
			names = append(names, candidate)
		}
	}
	css, err := collectDictionaryCSS(dict.DictDir, names, dict.Dict.LookupResource)
	if err != nil {
		log.Warningf("load dictionary CSS for %s failed: %s", dictID, err)
		return ""
	}
	return handler.InlineCSSResources(css, func(key string) ([]byte, bool) {
		if raw, err := b.dictSvc.FindFromDir(dictID, key); err == nil {
			return raw, true
		}
		for _, candidate := range dictionaryResourceCandidates(key) {
			if raw, err := b.dictSvc.LookupResource(dictID, candidate); err == nil {
				return raw, true
			}
		}
		return nil, false
	})
}

func dictionaryResourceCandidates(key string) []string {
	clean := strings.TrimSpace(key)
	trimmed := strings.TrimLeft(clean, `/\`)
	candidates := []string{clean, trimmed, `\` + trimmed, `/` + trimmed}
	seen := make(map[string]bool)
	result := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate != "" && !seen[candidate] {
			seen[candidate] = true
			result = append(result, candidate)
		}
	}
	return result
}

func collectDictionaryCSS(dir string, embeddedNames []string, lookup func(string) ([]byte, error)) (string, error) {
	parts := make([]string, 0)
	seen := make(map[string]bool)
	if dir != "" {
		names := make([]string, 0)
		err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".css") {
				rel, relErr := filepath.Rel(dir, path)
				if relErr != nil {
					return relErr
				}
				names = append(names, rel)
			}
			return nil
		})
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		sort.Strings(names)
		for _, name := range names {
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return "", err
			}
			content := string(data)
			seen[content] = true
			parts = append(parts, fmt.Sprintf("/* %s */\n%s", strings.ReplaceAll(filepath.ToSlash(name), "*/", "* /"), content))
		}
	}
	if lookup != nil {
		var lookupErr error
		for _, name := range embeddedNames {
			found := false
			for _, candidate := range []string{name, `\` + name, `/` + name} {
				data, err := lookup(candidate)
				if err == nil && len(data) > 0 {
					content := string(data)
					if !seen[content] {
						seen[content] = true
						parts = append(parts, fmt.Sprintf("/* %s */\n%s", strings.ReplaceAll(name, "*/", "* /"), content))
					}
					found = true
					break
				}
				if err != nil && !errors.Is(err, model.ErrNotFound) {
					lookupErr = err
					break
				}
			}
			if !found && lookupErr != nil {
				return "", lookupErr
			}
		}
	}
	return strings.Join(parts, "\n\n"), nil
}

var inlinedCSSPattern = regexp.MustCompile(`href="data:text/css;base64,([A-Za-z0-9+/=]+)"`)

func extractInlinedCSS(snapshot string) string {
	parts := make([]string, 0)
	seen := make(map[string]bool)
	for _, match := range inlinedCSSPattern.FindAllStringSubmatch(snapshot, -1) {
		data, err := base64.StdEncoding.DecodeString(match[1])
		if err != nil || len(data) == 0 || seen[string(data)] {
			continue
		}
		seen[string(data)] = true
		parts = append(parts, string(data))
	}
	return strings.Join(parts, "\n\n")
}

func readDictUserCSS(dictID string) (string, error) {
	cssDir, err := userCSSDir()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(cssDir, dictID+".css"))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	return string(data), err
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
