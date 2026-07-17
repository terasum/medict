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

// Package ecdict is an offline English→Chinese dictionary backed by a compact
// ECDICT (https://github.com/skywind3000/ECDICT) SQLite subset, shipped as a
// preset so Medict has a default EN-CN dictionary out of the box. It implements
// model.GeneralDictionary, so it plugs into the existing dict list / search /
// content pipeline (multi-dict, hover-popup, …) like mdict/stardict.
//
// Content path: Locate/Lookup return an HTML *fragment* (the entry card); the
// existing dicts_controller.HandleWordQueryReq → handler.WrapContent wraps it
// into a full page (with the inner script that powers hover/dblclick lookup),
// so no controller change is needed.
package ecdict

import (
	"database/sql"
	"fmt"
	"html"
	"path/filepath"
	"strings"

	"github.com/terasum/medict/pkg/model"

	// CGO-free pure-Go SQLite driver (database/sql); same driver the bookmark
	// store uses.
	_ "modernc.org/sqlite"
)

// ECDict is an offline EN-CN dictionary over a SQLite file (table `ecdict`:
// word, phonetic, definition, translation, pos, frq). Safe for read-only use.
type ECDict struct {
	db *sql.DB
}

// NewECDict opens ecdict.db inside the dict directory.
func NewECDict(dirItem *model.DirItem) (*ECDict, error) {
	dbPath := filepath.Join(dirItem.CurrentDir, "ecdict.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ecdict: open %s: %w", dbPath, err)
	}
	// Read-only dictionary lookups; single connection avoids SQLITE_BUSY.
	db.SetMaxOpenConns(1)
	return &ECDict{db: db}, nil
}

// BuildIndex ensures a prefix-friendly index on word. word is the PRIMARY KEY
// (already indexed), so this is essentially a no-op safety net.
func (e *ECDict) BuildIndex() error {
	_, err := e.db.Exec(`CREATE INDEX IF NOT EXISTS idx_ecdict_word ON ecdict(word)`)
	return err
}

func (e *ECDict) DictType() model.DictType { return model.DictTypeECDICT }

func (e *ECDict) Name() string { return "简明英汉词典" }

func (e *ECDict) Description() *model.PlainDictionaryInfo {
	return &model.PlainDictionaryInfo{
		Title:       "简明英汉词典",
		Description: "ECDICT 离线英汉(高频词条)",
	}
}

// Lookup returns the entry for an exact word as an HTML fragment ("" if absent).
func (e *ECDict) Lookup(keyword string) ([]byte, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("ecdict: empty keyword")
	}
	var phonetic, definition, translation string
	err := e.db.QueryRow(
		`SELECT phonetic, definition, translation FROM ecdict WHERE word = ?`,
		keyword,
	).Scan(&phonetic, &definition, &translation)
	if err == sql.ErrNoRows {
		return []byte(""), nil
	}
	if err != nil {
		return nil, err
	}
	return []byte(renderCard(keyword, phonetic, definition, translation)), nil
}

// Locate serves the word carried by the entry (offsets are irrelevant for
// ECDICT — only the keyword matters).
func (e *ECDict) Locate(entry *model.KeyQueryIndex) ([]byte, error) {
	kw := ""
	if entry != nil && entry.MdictKeyWordIndex != nil {
		kw = entry.KeyWord
	}
	return e.Lookup(kw)
}

// Search returns prefix matches (frq ascending = most frequent first), capped.
func (e *ECDict) Search(keyword string) ([]*model.KeyQueryIndex, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*model.KeyQueryIndex{}, nil
	}
	rows, err := e.db.Query(
		`SELECT word FROM ecdict WHERE word LIKE ? ORDER BY frq ASC LIMIT 50`,
		keyword+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]*model.KeyQueryIndex, 0, 50)
	for rows.Next() {
		var w string
		if err := rows.Scan(&w); err != nil {
			return nil, err
		}
		out = append(out, &model.KeyQueryIndex{
			IndexType:         string(model.DictTypeECDICT),
			MdictKeyWordIndex: &model.MdictKeyWordIndex{KeyWord: w},
		})
	}
	return out, rows.Err()
}

// LookupResource — ECDICT has no external resources (css/images/fonts).
func (e *ECDict) LookupResource(keyword string) ([]byte, error) {
	return nil, fmt.Errorf("ecdict: no resources (%q)", keyword)
}

// Close releases the database handle.
func (e *ECDict) Close() error {
	if e.db == nil {
		return nil
	}
	return e.db.Close()
}

// renderCard builds the entry HTML fragment (word + phonetic + CN translation +
// EN definition). WrapContent wraps it into the full page. Text is HTML-escaped;
// newlines become <br>.
func renderCard(word, phonetic, definition, translation string) string {
	var b strings.Builder
	b.WriteString(`<style>
.ec-card{font-family:-apple-system,"Segoe UI","Microsoft YaHei",sans-serif;color:#222;max-width:680px;padding:4px 2px;}
.ec-head{display:flex;align-items:baseline;gap:10px;margin-bottom:4px;flex-wrap:wrap;}
.ec-word{font-size:26px;font-weight:700;}
.ec-phon{color:#2a7de1;font-size:15px;}
.ec-sec-h{font-size:12px;color:#999;font-weight:600;margin:12px 0 4px;border-bottom:1px solid #eee;padding-bottom:2px;}
.ec-trans{font-size:16px;line-height:1.7;}
.ec-def{font-size:14px;line-height:1.6;color:#555;}
</style>`)
	b.WriteString(`<div class="ec-card"><div class="ec-head">`)
	fmt.Fprintf(&b, `<span class="ec-word">%s</span>`, html.EscapeString(word))
	if strings.TrimSpace(phonetic) != "" {
		fmt.Fprintf(&b, `<span class="ec-phon">/%s/</span>`, html.EscapeString(phonetic))
	}
	b.WriteString(`</div>`)
	if strings.TrimSpace(translation) != "" {
		b.WriteString(`<div class="ec-sec-h">释义</div><div class="ec-trans">`)
		b.WriteString(nl2br(html.EscapeString(translation)))
		b.WriteString(`</div>`)
	}
	if strings.TrimSpace(definition) != "" {
		b.WriteString(`<div class="ec-sec-h">English</div><div class="ec-def">`)
		b.WriteString(nl2br(html.EscapeString(definition)))
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func nl2br(s string) string {
	return strings.ReplaceAll(s, "\n", "<br>")
}
