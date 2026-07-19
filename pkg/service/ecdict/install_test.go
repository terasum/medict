package ecdict

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/terasum/medict/pkg/model"
)

const fixtureCSV = `word,phonetic,definition,translation,pos,collins,oxford,tag,bnc,frq,exchange,detail,audio
apple,'æpl,a round fruit,n. 苹果,n:100,5,1,,100,2695,,,
book,buk,a written work,"n. 书,书籍",n:80,5,1,,50,241,,,
empty,,,"",,,,,,,,,
`

func TestBuildDatabaseFromCSV(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "ecdict.db")
	count, err := buildDatabaseFromCSV(context.Background(), strings.NewReader(fixtureCSV), dbPath)
	if err != nil {
		t.Fatalf("buildDatabaseFromCSV: %v", err)
	}
	if count != 2 {
		t.Fatalf("count=%d want 2", count)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var translation string
	if err := db.QueryRow(`SELECT translation FROM ecdict WHERE word = 'book'`).Scan(&translation); err != nil {
		t.Fatal(err)
	}
	if translation != "n. 书,书籍" {
		t.Fatalf("translation=%q", translation)
	}
}

func TestInstallFromCSVReplacesOpenDatabase(t *testing.T) {
	dir := t.TempDir()
	seedDB(t, dir)
	dict, err := NewECDict(&model.DirItem{CurrentDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	defer dict.Close()

	count, err := dict.installFromCSV(context.Background(), io.NopCloser(strings.NewReader(fixtureCSV)), 2)
	if err != nil {
		t.Fatalf("installFromCSV: %v", err)
	}
	if count != 2 {
		t.Fatalf("count=%d want 2", count)
	}
	definition, err := dict.Lookup("apple")
	if err != nil || !strings.Contains(string(definition), "苹果") {
		t.Fatalf("replacement lookup=%q err=%v", definition, err)
	}
	if definition, _ := dict.Lookup("run"); string(definition) != "" {
		t.Fatalf("old database still active: %q", definition)
	}
}

func TestInstallFromCSVPreservesDatabaseOnInvalidSource(t *testing.T) {
	dir := t.TempDir()
	seedDB(t, dir)
	dict, err := NewECDict(&model.DirItem{CurrentDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	defer dict.Close()

	if _, err := dict.installFromCSV(context.Background(), io.NopCloser(strings.NewReader("not,a,valid,ecdict\n")), 1); err == nil {
		t.Fatal("expected invalid source error")
	}
	definition, err := dict.Lookup("run")
	if err != nil || !strings.Contains(string(definition), "运行") {
		t.Fatalf("original database was not preserved: %q err=%v", definition, err)
	}
}

func TestInstallFullLive(t *testing.T) {
	if os.Getenv("MEDICT_TEST_LIVE_ECDICT") == "" {
		t.Skip("set MEDICT_TEST_LIVE_ECDICT=1 to test the official full dataset")
	}
	dir := t.TempDir()
	seedDB(t, dir)
	dict, err := NewECDict(&model.DirItem{CurrentDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	defer dict.Close()

	count, err := dict.InstallFull(context.Background(), nil)
	if err != nil {
		t.Fatalf("InstallFull: %v", err)
	}
	if count < minimumFullRows {
		t.Fatalf("count=%d want >=%d", count, minimumFullRows)
	}
	t.Logf("installed %d official ECDICT entries", count)
	status, err := dict.Status()
	if err != nil || status.Edition != "full" || status.EntryCount != count {
		t.Fatalf("status=%+v err=%v count=%d", status, err, count)
	}
}
