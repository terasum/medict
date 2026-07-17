//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// GPL-3.0.

package ecdict

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/terasum/medict/pkg/model"

	_ "modernc.org/sqlite"
)

// seedDB creates a tiny ecdict.db (table `ecdict`) in dir for testing.
func seedDB(t *testing.T, dir string) {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(dir, "ecdict.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	schema := `CREATE TABLE ecdict (word TEXT PRIMARY KEY, phonetic TEXT, definition TEXT, translation TEXT, pos TEXT, frq INTEGER)`
	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}
	rows := []struct {
		word, phon, def, trans string
		frq                    int
	}{
		{"apple", "'æpl", "a round fruit", "n. 苹果", 2695},
		{"application", "ˌæplɪˈkeɪʃn", "a formal request", "n. 应用,申请", 800},
		{"book", "buk", "a written work", "n. 书,书籍\nv. 预订", 241},
		{"run", "rʌn", "to move fast", "v. 跑,运行\nn. 跑", 202},
	}
	for _, r := range rows {
		if _, err := db.Exec(
			`INSERT INTO ecdict(word,phonetic,definition,translation,frq) VALUES(?,?,?,?,?)`,
			r.word, r.phon, r.def, r.trans, r.frq,
		); err != nil {
			t.Fatal(err)
		}
	}
}

func TestECDict_LookupAndSearch(t *testing.T) {
	dir := t.TempDir()
	seedDB(t, dir)
	d, err := NewECDict(&model.DirItem{CurrentDir: dir})
	if err != nil {
		t.Fatalf("NewECDict: %v", err)
	}
	defer d.Close()

	if err := d.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}
	if d.DictType() != model.DictTypeECDICT {
		t.Fatalf("DictType=%s want ECDICT", d.DictType())
	}

	// Lookup 命中:fragment 含中文释义 + 词头
	def, err := d.Lookup("apple")
	if err != nil {
		t.Fatalf("Lookup apple: %v", err)
	}
	if s := string(def); !strings.Contains(s, "苹果") || !strings.Contains(s, "apple") {
		t.Fatalf("apple fragment missing CN/word: %q", s)
	}
	// Lookup 未命中 → 空片段
	if def, _ := d.Lookup("no_such_word"); string(def) != "" {
		t.Fatalf("missing word should give empty fragment, got %q", string(def))
	}

	// Locate 用 entry 的 keyword
	def2, _ := d.Locate(&model.KeyQueryIndex{MdictKeyWordIndex: &model.MdictKeyWordIndex{KeyWord: "book"}})
	if !strings.Contains(string(def2), "书") {
		t.Fatalf("book fragment missing CN: %q", string(def2))
	}

	// Search 前缀:ap → apple + application
	res, err := d.Search("ap")
	if err != nil {
		t.Fatalf("Search ap: %v", err)
	}
	words := make(map[string]bool)
	for _, r := range res {
		words[r.KeyWord] = true
	}
	if !words["apple"] || !words["application"] {
		t.Fatalf("Search ap results = %v, want apple+application", words)
	}
	// frq ASC:application(800) 应排在 apple(2695) 前
	if len(res) < 2 || res[0].KeyWord != "application" {
		t.Fatalf("Search ap order wrong: first=%q", firstWord(res))
	}

	// Search 无匹配 → 空
	if res, _ := d.Search("zzz"); len(res) != 0 {
		t.Fatalf("Search zzz should be empty, got %d", len(res))
	}

	// LookupResource:无资源
	if _, err := d.LookupResource("anything"); err == nil {
		t.Fatalf("LookupResource should error")
	}
}

func firstWord(res []*model.KeyQueryIndex) string {
	if len(res) == 0 {
		return ""
	}
	return res[0].KeyWord
}
