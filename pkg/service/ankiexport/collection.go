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

// Package ankiexport turns bookmark snapshots into a native Anki .apkg.
//
// An .apkg is a zip containing:
//   - collection.anki2  — a SQLite DB with Anki's schema (col/notes/cards/...)
//   - media             — JSON mapping numeric keys ("0","1",...) → media filenames
//   - the media bytes    — stored in the zip under the numeric keys
//
// This file builds collection.anki2. The SQLite schema, the col-row JSON
// (models/decks/dconf/conf), the new-card scheduling fields and the field
// checksum mirror Anki's own format (see Anki's db.py / utils.field_checksum,
// and genanki's defaults).
package ankiexport

import (
	"crypto/sha1"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	// CGO-free pure-Go SQLite driver (database/sql). Same driver the bookmark
	// store uses; we write a stand-alone .anki2 (journal OFF, clean close).
	_ "modernc.org/sqlite"
)

// ankiSchemaVer is the collection schema version we write. We deliberately use
// the legacy value 11 (matching genanki): on open, Anki upgrades a v11
// collection, migrating col.models JSON into the modern `notetypes` table.
// Writing a newer ver without the `notetypes` table makes Anki fail with
// "no such table: notetypes".
const ankiSchemaVer = 11

// modelID is the stable note-type id for the "Medict Word" model.
const modelID int64 = 1700000000119

// preparedRow is an ExportRow after media extraction: its Definition HTML has
// every <img>/audio data: URL rewritten to a media filename, and the bytes are
// collected by the media builder for the zip.
type preparedRow struct {
	Word       string
	DictName   string
	Notebook   string
	Definition string
}

// fieldChecksum mirrors Anki's utils.field_checksum: strip HTML/media, drop all
// whitespace, lowercase-not-applied (Anki lowercases only for display), take the
// first 4 bytes of SHA-1 as a little-endian uint32. Used for duplicate detection.
func fieldChecksum(field string) uint32 {
	s := stripHTMLMedia(field)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")
	s = anyWhitespace.ReplaceAllString(s, "")
	sum := sha1.Sum([]byte(s))
	return binary.LittleEndian.Uint32(sum[:4])
}

var (
	anyWhitespace  = regexp.MustCompile(`\s`)
	htmlTagOrMedia = regexp.MustCompile(`(?s)\[\[type:[^]]+\]\]|<[^>]+>|\[sound:[^]]*\]`)
)

// stripHTMLMedia removes HTML tags, [sound:..] refs and [[type:..]] cloze markers
// so the checksum is computed on visible text only (matches Anki.stripHTMLMedia).
func stripHTMLMedia(s string) string {
	return htmlTagOrMedia.ReplaceAllString(s, "")
}

// guidFor returns a stable, unique-enough note guid (Anki stores guids as TEXT;
// 10 hex chars from a content hash is plenty and keeps re-exports deterministic).
func guidFor(word, notebook, dict string) string {
	h := sha1.Sum([]byte(word + "\x1f" + notebook + "\x1f" + dict))
	return fmt.Sprintf("%x", h[:5]) // 10 hex chars
}

// deckID derives a stable, positive deck id from a notebook name.
func deckID(name string) int64 {
	h := sha1.Sum([]byte("deck:" + name))
	id := int64(binary.BigEndian.Uint64(h[:8]) % 1_000_000_000_000)
	if id < 0 {
		id = -id
	}
	if id == 0 {
		id = 1
	}
	return id
}

// buildCollection writes a valid collection.anki2 for the prepared rows into a
// temp file and returns its bytes. prepared must be non-empty.
func buildCollection(prepared []preparedRow, modMs, crtSec int64) ([]byte, error) {
	tmp, err := os.CreateTemp("", "medict-anki-*.anki2")
	if err != nil {
		return nil, fmt.Errorf("anki: temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	db, err := sql.Open("sqlite", tmpPath)
	if err != nil {
		return nil, fmt.Errorf("anki: open db: %w", err)
	}
	// Keep the file a single self-contained .anki2 (no -wal/-shm sidecars).
	if _, err := db.Exec(`PRAGMA journal_mode = OFF;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("anki: pragma journal: %w", err)
	}

	if err := createAnkiSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	// One deck per notebook (name "Medict::<notebook>"), plus the Default deck 1.
	decks := map[string]any{
		"1": deckDict(1, "Default", modMs),
	}
	for _, name := range uniqueNotebooks(prepared) {
		did := deckID(name)
		decks[strconv.FormatInt(did, 10)] = deckDict(did, "Medict::"+name, modMs)
	}

	models := map[string]any{strconv.FormatInt(modelID, 10): medictModel(modMs)}
	colRow, err := buildColRow(models, decks, defaultDconf(modMs), defaultConf(int64(len(prepared))), modMs, crtSec)
	if err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(colInsertSQL, colRow...); err != nil {
		db.Close()
		return nil, fmt.Errorf("anki: insert col: %w", err)
	}

	noteStmt, err := db.Prepare(noteInsertSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("anki: prepare note: %w", err)
	}
	defer noteStmt.Close()
	cardStmt, err := db.Prepare(cardInsertSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("anki: prepare card: %w", err)
	}
	defer cardStmt.Close()

	baseID := modMs
	var idCounter int64
	for _, r := range prepared {
		noteID := baseID + idCounter
		idCounter++
		cardID := baseID + idCounter
		idCounter++

		tags := cleanTag(r.Notebook)
		if d := cleanTag(r.DictName); d != "" {
			tags = strings.TrimSpace(tags + " " + d)
		}
		flds := r.Word + "\x1f" + r.Definition + "\x1f" + r.DictName
		if _, err := noteStmt.Exec(noteID, guidFor(r.Word, r.Notebook, r.DictName), modelID, modMs, -1,
			tags, flds, r.Word, int64(fieldChecksum(r.Word)), 0, ""); err != nil {
			db.Close()
			return nil, fmt.Errorf("anki: insert note %q: %w", r.Word, err)
		}
		// New card: type=0(new), queue=0(new), due=position, rest zeroed.
		duePos := idCounter / 2 // monotonic new-card position
		if _, err := cardStmt.Exec(cardID, noteID, deckID(r.Notebook), 0, modMs, -1,
			0, 0, duePos, 0, 0, 0, 0, 0, 0, 0, 0, ""); err != nil {
			db.Close()
			return nil, fmt.Errorf("anki: insert card %q: %w", r.Word, err)
		}
	}

	if err := db.Close(); err != nil {
		return nil, fmt.Errorf("anki: close db: %w", err)
	}
	return os.ReadFile(tmpPath)
}

// createAnkiSchema creates Anki's tables + indexes (DDL mirrors Anki db.py).
func createAnkiSchema(db *sql.DB) error {
	const ddl = `
CREATE TABLE col (
    crt     INTEGER NOT NULL,
    mod     INTEGER NOT NULL,
    scm     INTEGER NOT NULL,
    ver     INTEGER NOT NULL,
    dty     INTEGER NOT NULL,
    usn     INTEGER NOT NULL,
    ls      INTEGER NOT NULL,
    conf    TEXT NOT NULL,
    models  TEXT NOT NULL,
    decks   TEXT NOT NULL,
    dconf   TEXT NOT NULL,
    tags    TEXT NOT NULL
);
CREATE TABLE notes (
    id    INTEGER PRIMARY KEY,
    guid  TEXT NOT NULL,
    mid   INTEGER NOT NULL,
    mod   INTEGER NOT NULL,
    usn   INTEGER NOT NULL,
    tags  TEXT NOT NULL,
    flds  TEXT NOT NULL,
    sfld  INTEGER NOT NULL,
    csum  INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    data  TEXT NOT NULL
);
CREATE TABLE cards (
    id      INTEGER PRIMARY KEY,
    nid     INTEGER NOT NULL,
    did     INTEGER NOT NULL,
    ord     INTEGER NOT NULL,
    mod     INTEGER NOT NULL,
    usn     INTEGER NOT NULL,
    type    INTEGER NOT NULL,
    queue   INTEGER NOT NULL,
    due     INTEGER NOT NULL,
    ivl     INTEGER NOT NULL,
    factor  INTEGER NOT NULL,
    reps    INTEGER NOT NULL,
    lapses  INTEGER NOT NULL,
    left    INTEGER NOT NULL,
    odue    INTEGER NOT NULL,
    odid    INTEGER NOT NULL,
    flags   INTEGER NOT NULL,
    data    TEXT NOT NULL
);
CREATE TABLE revlog (
    id      INTEGER PRIMARY KEY,
    cid     INTEGER NOT NULL,
    usn     INTEGER NOT NULL,
    ease    INTEGER NOT NULL,
    ivl     INTEGER NOT NULL,
    lastIvl INTEGER NOT NULL,
    factor  INTEGER NOT NULL,
    time    INTEGER NOT NULL,
    type    INTEGER NOT NULL
);
CREATE TABLE graves (
    usn  INTEGER NOT NULL,
    oid  INTEGER NOT NULL,
    type INTEGER NOT NULL
);
CREATE INDEX ix_notes_usn ON notes (usn);
CREATE INDEX ix_notes_csum ON notes (csum);
CREATE INDEX ix_notes_dup ON notes (mid, csum);
CREATE INDEX ix_cards_usn ON cards (usn);
CREATE INDEX ix_cards_nid ON cards (nid);
CREATE INDEX ix_cards_sched ON cards (did, queue, due);
CREATE INDEX ix_revlog_usn ON revlog (usn);
CREATE INDEX ix_revlog_cid ON revlog (cid);
`
	_, err := db.Exec(ddl)
	return err
}

const (
	colInsertSQL = `INSERT INTO col(crt,mod,scm,ver,dty,usn,ls,conf,models,decks,dconf,tags) VALUES(?,?,?,?,?,?,?,?,?,?,?,?);`
	noteInsertSQL = `INSERT INTO notes(id,guid,mid,mod,usn,tags,flds,sfld,csum,flags,data) VALUES(?,?,?,?,?,?,?,?,?,?,?);`
	cardInsertSQL = `INSERT INTO cards(id,nid,did,ord,mod,usn,type,queue,due,ivl,factor,reps,lapses,left,odue,odid,flags,data) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
)

// buildColRow returns the 12 col-column values, marshalling the four JSON blobs.
func buildColRow(models, decks, dconf, conf map[string]any, modMs, crtSec int64) ([]any, error) {
	mkJSON := func(v any) (any, error) {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	}
	confJ, err := mkJSON(conf)
	if err != nil {
		return nil, fmt.Errorf("anki: marshal conf: %w", err)
	}
	modelsJ, err := mkJSON(models)
	if err != nil {
		return nil, fmt.Errorf("anki: marshal models: %w", err)
	}
	decksJ, err := mkJSON(decks)
	if err != nil {
		return nil, fmt.Errorf("anki: marshal decks: %w", err)
	}
	dconfJ, err := mkJSON(dconf)
	if err != nil {
		return nil, fmt.Errorf("anki: marshal dconf: %w", err)
	}
	// tags column is parsed as JSON by Anki's backend on open (ver=11 upgrade),
	// so it must be a valid JSON object ("{}"), not an empty string (which throws
	// "EOF while parsing a value").
	return []any{crtSec, modMs, crtSec * 1000, ankiSchemaVer, 0, -1, 0, confJ, modelsJ, decksJ, dconfJ, "{}"}, nil
}

// medictModel returns the "Medict Word" note type: fields Word/Definition/Dict,
// one forward template (front = Word, back = Definition).
func medictModel(modMs int64) map[string]any {
	return map[string]any{
		"id":      modelID,
		"name":    "Medict Word",
		"type":    0,
		"mod":     modMs,
		"usn":     -1,
		"sortf":   0,
		"did":     nil,
		"tmpls": []map[string]any{{
			"name":  "Word-Definition",
			"ord":   0,
			"qfmt":  "{{Word}}",
			"afmt":  `{{FrontSide}}<hr id="answer">{{Definition}}<br><div style="color:#999;font-size:14px">{{Dict}}</div>`,
			"bqfmt": "",
			"bafmt": "",
			"did":   nil,
			"mode":  0,
		}},
		"flds": []map[string]any{
			{"name": "Word", "ord": 0, "sticky": false, "rtl": false, "font": "Arial", "size": 28, "media": []any{}},
			{"name": "Definition", "ord": 1, "sticky": false, "rtl": false, "font": "Arial", "size": 20, "media": []any{}},
			{"name": "Dict", "ord": 2, "sticky": false, "rtl": false, "font": "Arial", "size": 14, "media": []any{}},
		},
		"css":           `.card{font-family:Arial;font-size:20px;text-align:left;color:#222;background:#fff;line-height:1.5;}`,
		"req":           [][]any{{"Word-Definition", "any", []any{0}}},
		"latexPostamble": "",
		"latexPreamble":  "",
		"latexsvg":       false,
		"vers":           []any{},
		"tags":           []any{},
		"version":        0,
	}
}

// deckDict returns one entry of the decks JSON.
func deckDict(id int64, name string, modMs int64) map[string]any {
	return map[string]any{
		"id":                id,
		"name":              name,
		"mod":               modMs,
		"usn":               -1,
		"desc":              "",
		"dyn":               0,
		"collapsed":         false,
		"browserCollapsed":  false,
		"extendNew":         10,
		"extendRev":         50,
		"conf":              1,
		"newToday":          []any{0, 0},
		"revToday":          []any{0, 0},
		"lrnToday":          []any{0, 0},
		"timeToday":         []any{0, 0},
	}
}

// defaultDconf returns the default deck config (id 1) Anki ships with.
func defaultDconf(modMs int64) map[string]any {
	return map[string]any{
		"1": map[string]any{
			"id":        1,
			"mod":       modMs,
			"usn":       -1,
			"name":      "Default",
			"autoplay":  true,
			"dyn":       0,
			"maxTaken":  60,
			"timer":     0,
			"replayq":   true,
			"hash":      "37c3a0581407a8d3",
			"new": map[string]any{
				"bury":         false,
				"delays":       []any{1, 10},
				"initialFactor": 2500,
				"ints":          []any{1, 4, 7},
				"order":         1,
				"perDay":        20,
				"separate":      false,
			},
			"rev": map[string]any{
				"bury":       false,
				"ease4":      1.3,
				"fuzz":       0.05,
				"ivlFct":     1.0,
				"maxIvl":     36500,
				"minSpace":   1,
				"perDay":     200,
				"hardFactor": 1.2,
			},
			"lapse": map[string]any{
				"delays":      []any{10},
				"leechAction": 1,
				"leechFails":  8,
				"minInt":      1,
				"mult":        0.5,
			},
		},
	}
}

// defaultConf returns the col.conf scheduler blob. nextPos is the next new-card
// position after this import.
func defaultConf(nextPos int64) map[string]any {
	return map[string]any{
		"activeDecks":    []any{int64(1)},
		"curDeck":        int64(1),
		"curModel":       nil,
		"dueCounts":      true,
		"estTimes":       true,
		"newBury":        true,
		"newSpread":      int64(0),
		"nextPos":        nextPos,
		"sortBackwards":  false,
		"schedVer":       int64(1),
	}
}

// uniqueNotebooks returns the sorted unique notebook names among the rows.
func uniqueNotebooks(rows []preparedRow) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0)
	for _, r := range rows {
		if r.Notebook == "" {
			r.Notebook = "Default"
		}
		if _, ok := seen[r.Notebook]; ok {
			continue
		}
		seen[r.Notebook] = struct{}{}
		out = append(out, r.Notebook)
	}
	sort.Strings(out)
	return out
}

// cleanTag turns an arbitrary label into an Anki-safe tag (lowercased, no
// spaces/special chars; spaces → underscores).
func cleanTag(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "/", "_")
	return s
}
