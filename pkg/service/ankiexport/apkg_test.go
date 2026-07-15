//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ankiexport

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

// png1x1B64 is a valid 1x1 transparent PNG.
const png1x1B64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+M8AAAMBAQDJ/IvNAAAAAElFTkSuQmCC"

func TestFieldChecksum_MatchesDefinition(t *testing.T) {
	// Anki's field_checksum = LE uint32 of first 4 bytes of SHA-1 over the
	// HTML-stripped, whitespace-stripped text.
	sum := sha1.Sum([]byte("hello"))
	want := binary.LittleEndian.Uint32(sum[:4])
	got := fieldChecksum("hello")
	if got != want {
		t.Fatalf("fieldChecksum(hello)=%d want %d", got, want)
	}
	// HTML + whitespace must be stripped before hashing.
	if fieldChecksum("<b>hello</b>") != fieldChecksum("hello") {
		t.Fatalf("fieldChecksum should ignore HTML tags")
	}
	if fieldChecksum("h e llo") != fieldChecksum("hello") {
		t.Fatalf("fieldChecksum should ignore whitespace")
	}
}

func TestStripHTMLMedia(t *testing.T) {
	got := stripHTMLMedia(`<b>hi</b> [sound:a.mp3] <img src="x.png"> [[type:cloze]]`)
	if got != "hi  " && strings.TrimSpace(got) != "hi" {
		t.Fatalf("stripHTMLMedia = %q", got)
	}
}

func TestParseDataURL(t *testing.T) {
	raw, err := base64.StdEncoding.DecodeString(png1x1B64)
	if err != nil {
		t.Fatal(err)
	}
	mime, data, ok := parseDataURL("data:image/png;base64," + png1x1B64)
	if !ok {
		t.Fatal("expected ok")
	}
	if mime != "image/png" {
		t.Fatalf("mime=%q want image/png", mime)
	}
	if !bytes.Equal(data, raw) {
		t.Fatalf("decoded bytes mismatch")
	}
	if _, _, ok := parseDataURL("https://x/y"); ok {
		t.Fatal("non-data URL should not parse")
	}
}

func TestExtractDefinition_BodyAndStyle(t *testing.T) {
	doc := `<html><head><style>.x{color:red}</style></head><body><p>hello</p></body></html>`
	def := extractDefinition(doc)
	if !strings.Contains(def, ".x{color:red}") {
		t.Fatalf("style block should be hoisted: %q", def)
	}
	if !strings.Contains(def, "<p>hello</p>") {
		t.Fatalf("body inner should be present: %q", def)
	}
	if strings.Contains(def, "<body") {
		t.Fatalf("body tag should be stripped: %q", def)
	}
}

func TestMediaCollector_RewriteAndDedup(t *testing.T) {
	m := newMediaCollector()
	img := `<img src="data:image/png;base64,` + png1x1B64 + `">`
	out := m.rewrite(img + img) // same image twice
	if len(m.files) != 1 {
		t.Fatalf("expected 1 deduped media file, got %d", len(m.files))
	}
	if !strings.Contains(m.files[0].name, ".png") {
		t.Fatalf("media filename should end with .png: %q", m.files[0].name)
	}
	// Both occurrences rewritten to the same filename.
	if c := strings.Count(out, m.files[0].name); c != 2 {
		t.Fatalf("expected 2 refs to media file, got %d in %q", c, out)
	}
	if strings.Contains(out, "data:image/png") {
		t.Fatalf("data: URL should be replaced: %q", out)
	}
}

// TestBuild_endToEnd builds a 2-notebook apkg, reopens the zip + the SQLite
// collection, and asserts the Anki-shaped contents.
func TestBuild_endToEnd(t *testing.T) {
	rows := []ExportRow{
		{
			Word:         "apple",
			DictName:     "OALD",
			NotebookName: "Fruit",
			HTML:         `<html><head><style>p{color:#000}</style></head><body><p>a fruit</p><img src="data:image/png;base64,` + png1x1B64 + `"></body></html>`,
		},
		{
			Word:         "dog",
			DictName:     "OALD",
			NotebookName: "Animals",
			HTML:         `<html><body><p>a pet</p></body></html>`,
		},
	}
	apkg, err := Build(rows)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(apkg), int64(len(apkg)))
	if err != nil {
		t.Fatalf("zip read: %v", err)
	}
	files := map[string]*zip.File{}
	for _, f := range zr.File {
		files[f.Name] = f
	}
	if _, ok := files["collection.anki2"]; !ok {
		t.Fatal("collection.anki2 missing from apkg")
	}
	mediaEntry, ok := files["media"]
	if !ok {
		t.Fatal("media map missing from apkg")
	}

	// The apple card extracted one image → one media file under key "0".
	mediaMap := readJSONMap(t, mediaEntry)
	mediaName, hasImg := mediaMap["0"]
	if !hasImg {
		t.Fatalf("expected media entry 0, got %v", mediaMap)
	}
	if !strings.HasSuffix(mediaName, ".png") {
		t.Fatalf("media 0 should be png: %q", mediaName)
	}
	if _, ok := files["0"]; !ok {
		t.Fatal("media file bytes missing under key 0")
	}

	// Reopen collection.anki2 and assert Anki shape.
	col := openCol(t, files["collection.anki2"])
	defer col.Close()

	var noteCnt, cardCnt int
	if err := col.QueryRow(`SELECT COUNT(*) FROM notes`).Scan(&noteCnt); err != nil {
		t.Fatal(err)
	}
	if err := col.QueryRow(`SELECT COUNT(*) FROM cards`).Scan(&cardCnt); err != nil {
		t.Fatal(err)
	}
	if noteCnt != 2 || cardCnt != 2 {
		t.Fatalf("notes=%d cards=%d, want 2/2", noteCnt, cardCnt)
	}

	// The apple note's flds must reference the extracted image filename and carry
	// the word in field 0 (flds joined by 0x1f).
	var appleFlds string
	if err := col.QueryRow(`SELECT flds FROM notes WHERE tags LIKE ?`, "%fruit%").Scan(&appleFlds); err != nil {
		// tag is lowercased notebook name; fall back to scanning all if tag guess misses
		if err := col.QueryRow(`SELECT flds FROM notes ORDER BY id LIMIT 1 OFFSET 0`).Scan(&appleFlds); err != nil {
			t.Fatal(err)
		}
	}
	if !strings.Contains(appleFlds, mediaName) {
		t.Fatalf("apple flds should reference media %q: %q", mediaName, appleFlds)
	}
	parts := strings.Split(appleFlds, "\x1f")
	if len(parts) < 2 || parts[0] != "apple" {
		t.Fatalf("field0 should be the word 'apple': %q", appleFlds)
	}

	// col row present with models/decks JSON.
	var modelsJSON, decksJSON string
	if err := col.QueryRow(`SELECT models, decks FROM col`).Scan(&modelsJSON, &decksJSON); err != nil {
		t.Fatal(err)
	}
	// Anki-compat guards: ver must be 11 (older schema so Anki upgrades
	// col.models JSON → notetypes table on open) and tags must be "{}" (the
	// backend parses col.tags as JSON; an empty string throws EOF). These two
	// are the exact bugs found by the AnkiPackageImporter validation.
	var colVer int
	var colTags string
	if err := col.QueryRow(`SELECT ver, tags FROM col`).Scan(&colVer, &colTags); err != nil {
		t.Fatal(err)
	}
	if colVer != 11 {
		t.Fatalf("col.ver=%d want 11", colVer)
	}
	if colTags != "{}" {
		t.Fatalf("col.tags=%q want {}", colTags)
	}
	if !strings.Contains(modelsJSON, "Medict Word") {
		t.Fatalf("models JSON missing Medict Word: %s", modelsJSON)
	}
	if !strings.Contains(decksJSON, "Medict::Fruit") || !strings.Contains(decksJSON, "Medict::Animals") {
		t.Fatalf("decks JSON missing notebook decks: %s", decksJSON)
	}
	// Each note's card lives in its notebook's deck.
	var dids []int64
	rowsQ, err := col.Query(`SELECT DISTINCT did FROM cards`)
	if err != nil {
		t.Fatal(err)
	}
	for rowsQ.Next() {
		var d int64
		rowsQ.Scan(&d)
		dids = append(dids, d)
	}
	rowsQ.Close()
	if len(dids) != 2 {
		t.Fatalf("expected 2 distinct deck ids, got %d (%v)", len(dids), dids)
	}
}

func TestBuild_Empty(t *testing.T) {
	if _, err := Build(nil); err != ErrNoBookmarks {
		t.Fatalf("Build(nil) err=%v want ErrNoBookmarks", err)
	}
}

// --- helpers ---

func readJSONMap(t *testing.T, f *zip.File) map[string]string {
	t.Helper()
	rc, err := f.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(rc); err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// openCol writes the zipped collection bytes to a temp file and opens it.
func openCol(t *testing.T, f *zip.File) *sql.DB {
	t.Helper()
	rc, err := f.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(rc); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "collection.anki2")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
