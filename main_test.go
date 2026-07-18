package main

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/terasum/medict/pkg/model"
)

func TestParseCSSWindowArgs(t *testing.T) {
	dictID, dictName, seedPath, ok := parseCSSWindowArgs([]string{cssEditorWindowFlag, "online-bing", "Bing Dictionary", "/tmp/base.css"})
	if !ok || dictID != "online-bing" || dictName != "Bing Dictionary" || seedPath != "/tmp/base.css" {
		t.Fatalf("unexpected parse result: id=%q name=%q seed=%q ok=%v", dictID, dictName, seedPath, ok)
	}
	if _, _, _, ok := parseCSSWindowArgs([]string{cssEditorWindowFlag, "../escape", "bad", ""}); ok {
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
	cmd := newCSSWindowCommand("/Applications/Medict.app/Contents/MacOS/Medict", "dict-1", "Test Dict", "/tmp/base.css")
	if len(cmd.Args) < 9 || cmd.Args[0] != "/usr/bin/open" || cmd.Args[3] != "/Applications/Medict.app" || cmd.Args[8] != "/tmp/base.css" {
		t.Fatalf("packaged macOS editor must use LaunchServices: %#v", cmd.Args)
	}
}

func TestCSSWindowAppDoesNotInitialiseProductServices(t *testing.T) {
	app := newCSSWindowApp("abc123", "Test Dictionary", "/tmp/base.css")
	if app.WindowMode() != "css-editor" || app.dictSvc != nil || app.bs != nil {
		t.Fatalf("CSS editor must stay lightweight: %#v", app)
	}
	context := app.CSSWindowDictionary()
	if context["id"] != "abc123" || context["name"] != "Test Dictionary" {
		t.Fatalf("unexpected editor context: %#v", context)
	}
}

func TestCollectDictionaryCSSUsesLooseFilesInStableOrder(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "b.css"), []byte(".b{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.css"), []byte(".a{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	css, err := collectDictionaryCSS(dir, []string{"Dictionary.css"}, func(string) ([]byte, error) {
		return nil, model.ErrNotFound
	})
	if err != nil {
		t.Fatal(err)
	}
	if css != "/* a.css */\n.a{}\n\n/* b.css */\n.b{}" {
		t.Fatalf("unexpected CSS composition: %q", css)
	}
}

func TestCollectDictionaryCSSFallsBackToEmbeddedResource(t *testing.T) {
	css, err := collectDictionaryCSS(t.TempDir(), []string{"Fancy.css"}, func(name string) ([]byte, error) {
		if name == `\Fancy.css` {
			return []byte(".embedded{}"), nil
		}
		return nil, model.ErrNotFound
	})
	if err != nil {
		t.Fatal(err)
	}
	if css != "/* Fancy.css */\n.embedded{}" {
		t.Fatalf("unexpected embedded CSS: %q", css)
	}
}

func TestCollectDictionaryCSSPropagatesNotReadyError(t *testing.T) {
	notReady := errors.New("dictionary not ready")
	_, err := collectDictionaryCSS(t.TempDir(), []string{"book.css"}, func(string) ([]byte, error) {
		return nil, notReady
	})
	if !errors.Is(err, notReady) {
		t.Fatalf("expected readiness error, got %v", err)
	}
}

func TestCollectDictionaryCSSIncludesNestedLooseAndEmbeddedSources(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "css"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "css", "style.css"), []byte(".nested{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	css, err := collectDictionaryCSS(dir, []string{"book.css"}, func(name string) ([]byte, error) {
		if name == "book.css" {
			return []byte(".embedded{}"), nil
		}
		return nil, model.ErrNotFound
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "/* css/style.css */\n.nested{}\n\n/* book.css */\n.embedded{}"
	if css != want {
		t.Fatalf("mixed CSS: got %q want %q", css, want)
	}
}

func TestExtractInlinedCSSSupportsArbitraryMultipleStylesheets(t *testing.T) {
	first := base64.StdEncoding.EncodeToString([]byte(".style{}"))
	second := base64.StdEncoding.EncodeToString([]byte(".oald{}"))
	snapshot := `<link href="data:text/css;base64,` + first + `"><link href="data:text/css;base64,` + second + `">`
	if got := extractInlinedCSS(snapshot); got != ".style{}\n\n.oald{}" {
		t.Fatalf("unexpected extracted CSS: %q", got)
	}
}

func TestReadEditorCSSPrefersOverrideThenSeed(t *testing.T) {
	dir := t.TempDir()
	override := filepath.Join(dir, "override.css")
	seed := filepath.Join(dir, "seed.css")
	if err := os.WriteFile(seed, []byte("seed"), 0o600); err != nil {
		t.Fatal(err)
	}
	css, err := readEditorCSS(override, seed)
	if err != nil || css != "seed" {
		t.Fatalf("seed fallback: css=%q err=%v", css, err)
	}
	if err := os.WriteFile(override, []byte("override"), 0o600); err != nil {
		t.Fatal(err)
	}
	css, err = readEditorCSS(override, seed)
	if err != nil || css != "override" {
		t.Fatalf("override priority: css=%q err=%v", css, err)
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
