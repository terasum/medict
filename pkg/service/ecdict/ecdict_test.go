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
		// ECDICT's CSV encodes in-field line breaks as a literal backslash-n;
		// the importer stores those bytes verbatim (see preset ecdict.db: 10k+
		// rows affected). "go" here mirrors that real-world form.
		{"go", "gəu", "v. move away\\nfrom a place", "vi. 去, 走\\nvt. 忍受", 300},
		// base forms for the inflection-fallback tests (ECDICT stores base
		// forms; parts/studies/boxes/loved are absent on purpose)
		{"part", "pɑ:t", "n. a piece", "n. 部分", 1000},
		{"study", "'stʌdi", "v. learn", "v./n. 学习", 900},
		{"box", "bɔks", "n. container", "n. 盒子", 500},
		{"love", "lʌv", "v. feel affection", "v./n. 爱", 800},
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

	// ECDICT 的 CSV 把字段内换行存为字面量 \n,渲染时必须转为 <br>
	goFrag, err := d.Lookup("go")
	if err != nil {
		t.Fatalf("Lookup go: %v", err)
	}
	s := string(goFrag)
	if strings.Contains(s, `\n`) {
		t.Fatalf("literal \\n leaked into fragment: %q", s)
	}
	if !strings.Contains(s, "去, 走<br>vt. 忍受") {
		t.Fatalf("literal \\n not rendered as <br>: %q", s)
	}

	// 屈折回退:库里只有基础词形,parts/studies/ran 等变形必须回退到 part/study/run
	for _, tc := range []struct{ in, want string }{
		{"parts", "部分"},
		{"studies", "学习"},
		{"boxes", "盒子"},
		{"loved", "爱"},
		{"ran", "跑,运行"},
	} {
		frag, lerr := d.Lookup(tc.in)
		if lerr != nil {
			t.Fatalf("Lookup %s: %v", tc.in, lerr)
		}
		if !strings.Contains(string(frag), tc.want) {
			t.Fatalf("inflection fallback %s: want fragment containing %q, got %q", tc.in, tc.want, string(frag))
		}
	}

	// Search 同样回退:parts 前缀无结果 → 以 part 前缀重查
	sres, serr := d.Search("parts")
	if serr != nil {
		t.Fatalf("Search parts: %v", serr)
	}
	if len(sres) == 0 || sres[0].KeyWord != "part" {
		t.Fatalf("Search parts should fall back to prefix 'part', got %v", sres)
	}

	// 中文等非 ASCII 输入不做屈折回退(规则仅适用英文)
	if frag, _ := d.Lookup("苹果"); strings.Contains(string(frag), "苹果") {
		t.Fatalf("non-English input should not match, got %q", string(frag))
	}
	if lemmaCandidates("苹果") != nil {
		t.Fatalf("lemmaCandidates must yield nothing for non-ASCII input")
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

// TestCaseSensitiveLikeEnabled guards the PRAGMA case_sensitive_like = ON
// invariant set in NewECDict. Without it, SQLite's LIKE is case-insensitive and
// CANNOT use the primary key index → full table scan (12ms vs 0.1ms on 50K
// rows). If this test fails, someone removed the PRAGMA in NewECDict; restore it.
//
// We can't query `PRAGMA case_sensitive_like` (SQLite doesn't support reading
// it). Instead we verify the BEHAVIOR: LIKE 'app%' must be case-sensitive (match
// only "apple", not "Apple").
func TestCaseSensitiveLikeEnabled(t *testing.T) {
	dir := t.TempDir()
	// Create a DB with mixed-case words to test case sensitivity.
	db, err := sql.Open("sqlite", filepath.Join(dir, "ecdict.db"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE ecdict (word TEXT PRIMARY KEY, phonetic TEXT, definition TEXT, translation TEXT, pos TEXT, frq INTEGER)`); err != nil {
		t.Fatal(err)
	}
	for _, w := range []string{"apple", "Apple", "application"} {
		if _, err := db.Exec(`INSERT INTO ecdict(word,frq) VALUES(?,1)`, w); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()

	d, err := NewECDict(&model.DirItem{CurrentDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	// case_sensitive_like=ON → LIKE 'app%' matches "apple","application" but NOT "Apple"
	// case_sensitive_like=OFF → LIKE 'app%' matches all three
	rows, err := d.db.Query(`SELECT word FROM ecdict WHERE word LIKE 'app%'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var matches []string
	for rows.Next() {
		var w string
		rows.Scan(&w)
		matches = append(matches, w)
	}

	hasApple := false
	for _, w := range matches {
		if w == "Apple" {
			hasApple = true
		}
	}
	if hasApple {
		t.Fatalf("case_sensitive_like must be ON: LIKE 'app%%' matched 'Apple' "+
			"(case-insensitive). This means the PRAGMA in NewECDict was removed → "+
			"LIKE will do a full table scan instead of using the PK index "+
			"(12ms vs 0.1ms on 50K rows). Restore: db.Exec(\"PRAGMA case_sensitive_like = ON\"). "+
			"Matches were: %v", matches)
	}
	if len(matches) < 2 {
		t.Fatalf("LIKE 'app%%' should match at least apple+application, got %v", matches)
	}
}

func firstWord(res []*model.KeyQueryIndex) string {
	if len(res) == 0 {
		return ""
	}
	return res[0].KeyWord
}

// ─── Benchmarks (go test -bench=. -benchmem) ───
//
// These cover the core Search/Lookup/Locate operations on the seed DB (4 rows).
// For full-scale benchmarking with the 50K preset + pprof profiling, use:
//
//	go run ./cmd/benchmark
//
// The tests below verify the performance CHARACTERISTICS that matter:
//   - Search (LIKE with case_sensitive_like=ON) should use the PK index
//   - Lookup (exact PK) should be constant-time
//   - Locate should be ≈ Lookup

// benchECDict opens a seeded ECDict for benchmarking (shared setup).
func benchECDict(b *testing.B) *ECDict {
	b.Helper()
	dir := b.TempDir()
	seedDBTB(b, dir)
	d, err := NewECDict(&model.DirItem{CurrentDir: dir})
	if err != nil {
		b.Fatal(err)
	}
	if err := d.BuildIndex(); err != nil {
		b.Fatal(err)
	}
	return d
}

// seedDBTB is the testing.B variant of seedDB.
func seedDBTB(b *testing.B, dir string) {
	b.Helper()
	db, err := sql.Open("sqlite", filepath.Join(dir, "ecdict.db"))
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE ecdict (word TEXT PRIMARY KEY, phonetic TEXT, definition TEXT, translation TEXT, pos TEXT, frq INTEGER)`); err != nil {
		b.Fatal(err)
	}
	for _, r := range []struct{ word, phon, def, trans string; frq int }{
		{"apple", "", "fruit", "n. 苹果", 2695},
		{"application", "", "request", "n. 应用", 800},
		{"book", "", "written work", "n. 书", 241},
		{"run", "", "move fast", "v. 跑", 202},
	} {
		if _, err := db.Exec(`INSERT INTO ecdict(word,phonetic,definition,translation,frq) VALUES(?,?,?,?,?)`, r.word, r.phon, r.def, r.trans, r.frq); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_PrefixHit(b *testing.B) {
	d := benchECDict(b)
	defer d.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := d.Search("ap"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_Miss(b *testing.B) {
	d := benchECDict(b)
	defer d.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := d.Search("zzzzz"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLookup_Hit(b *testing.B) {
	d := benchECDict(b)
	defer d.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := d.Lookup("apple"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLookup_Miss(b *testing.B) {
	d := benchECDict(b)
	defer d.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := d.Lookup("xxnotaword"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLocate_Pipeline(b *testing.B) {
	d := benchECDict(b)
	defer d.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := d.Search("ap")
		if err != nil || len(res) == 0 {
			b.Fatal(err)
		}
		if _, err := d.Locate(res[0]); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuildIndex measures the CREATE INDEX IF NOT EXISTS call
// (should be near-instant on a 4-row DB; for real scale see cmd/benchmark).
func BenchmarkBuildIndex(b *testing.B) {
	d := benchECDict(b)
	defer d.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := d.BuildIndex(); err != nil {
			b.Fatal(err)
		}
	}
}
