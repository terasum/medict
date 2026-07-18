package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseCSSWindowArgs(t *testing.T) {
	dictID, dictName, ok := parseCSSWindowArgs([]string{cssEditorWindowFlag, "online-bing", "Bing Dictionary"})
	if !ok || dictID != "online-bing" || dictName != "Bing Dictionary" {
		t.Fatalf("unexpected parse result: id=%q name=%q ok=%v", dictID, dictName, ok)
	}
	if _, _, ok := parseCSSWindowArgs([]string{cssEditorWindowFlag, "../escape", "bad"}); ok {
		t.Fatal("unsafe dictionary id must be rejected")
	}
}

func TestReplaceProcessEnv(t *testing.T) {
	got := replaceProcessEnv([]string{"PATH=/bin", "devserver=localhost:34115"}, "devserver", "localhost:0")
	if len(got) != 2 || got[0] != "PATH=/bin" || got[1] != "devserver=localhost:0" {
		t.Fatalf("unexpected environment: %#v", got)
	}
}

func TestNewCSSWindowCommand(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS LaunchServices behavior")
	}
	cmd := newCSSWindowCommand("/Applications/Medict.app/Contents/MacOS/Medict", "dict-1", "Test Dict")
	if len(cmd.Args) < 8 || cmd.Args[0] != "/usr/bin/open" || cmd.Args[3] != "/Applications/Medict.app" {
		t.Fatalf("packaged macOS editor must use LaunchServices: %#v", cmd.Args)
	}
}

func TestCSSWindowAppDoesNotInitialiseProductServices(t *testing.T) {
	app := newCSSWindowApp("abc123", "Test Dictionary")
	if app.WindowMode() != "css-editor" || app.dictSvc != nil || app.bs != nil {
		t.Fatalf("CSS editor must stay lightweight: %#v", app)
	}
	context := app.CSSWindowDictionary()
	if context["id"] != "abc123" || context["name"] != "Test Dictionary" {
		t.Fatalf("unexpected editor context: %#v", context)
	}
}

func TestWriteFileAtomicReplacesContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.css")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(path, []byte("new content"), 0o644); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "new content" {
		t.Fatalf("atomic replacement failed: data=%q err=%v", data, err)
	}
}
