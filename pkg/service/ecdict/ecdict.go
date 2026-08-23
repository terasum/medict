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
// content pipeline (multi-dict, in-entry click lookup, …) like mdict/stardict.
//
// Content path: Locate/Lookup return an HTML *fragment* (the entry card); the
// existing dicts_controller.HandleWordQueryReq → handler.WrapContent wraps it
// into a full page (with the inner script that powers hover hints/click lookup),
// so no controller change is needed.
package ecdict

import (
	"database/sql"
	"fmt"
	"html"
	"path/filepath"
	"strings"
	"sync"

	"github.com/terasum/medict/pkg/model"

	// CGO-free pure-Go SQLite driver (database/sql); same driver the bookmark
	// store uses.
	_ "modernc.org/sqlite"
)

// ECDict is an offline EN-CN dictionary over a SQLite file (table `ecdict`:
// word, phonetic, definition, translation, pos, frq). Safe for read-only use.
type ECDict struct {
	mu     sync.RWMutex
	db     *sql.DB
	dbPath string
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
	// Enable case-sensitive LIKE so SQLite can use the primary key index for
	// LIKE 'prefix%' queries. Without this, LIKE is case-insensitive and forces
	// a full table scan (12ms vs 0.1ms on 50K rows).
	db.Exec("PRAGMA case_sensitive_like = ON")
	return &ECDict{db: db, dbPath: dbPath}, nil
}

// BuildIndex ensures a prefix-friendly index on word. word is the PRIMARY KEY
// (already indexed), so this is essentially a no-op safety net.
func (e *ECDict) BuildIndex() error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, err := e.db.Exec(`CREATE INDEX IF NOT EXISTS idx_ecdict_word ON ecdict(word)`)
	return err
}

func (e *ECDict) DictType() model.DictType { return model.DictTypeECDICT }

func (e *ECDict) Name() string { return "简明英汉词典" }

func (e *ECDict) Description() *model.PlainDictionaryInfo {
	description := "ECDICT 离线英汉（5 万高频词条）"
	if status, err := e.Status(); err == nil && status.Edition == "full" {
		description = fmt.Sprintf("ECDICT 完整离线英汉（%d 词条）", status.EntryCount)
	}
	return &model.PlainDictionaryInfo{
		Title:       "简明英汉词典",
		Description: description,
	}
}

// Lookup returns the entry for an exact word as an HTML fragment ("" if absent).
// On an exact miss it retries with lemmaCandidates ("parts"→"part", "ran"→"run"):
// ECDICT stores base forms almost exclusively, so inflected input would
// otherwise always come back empty (roadmap #786).
func (e *ECDict) Lookup(keyword string) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("ecdict: empty keyword")
	}
	for _, cand := range lookupCandidates(keyword) {
		var phonetic, definition, translation string
		err := e.db.QueryRow(
			`SELECT phonetic, definition, translation FROM ecdict WHERE word = ?`,
			cand,
		).Scan(&phonetic, &definition, &translation)
		if err == sql.ErrNoRows {
			continue // try the next candidate (typed form → lowercase → lemmas)
		}
		if err != nil {
			return nil, err
		}
		return []byte(renderCard(cand, phonetic, definition, translation)), nil
	}
	return []byte(""), nil
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
// With case_sensitive_like = ON, SQLite uses the primary key index for LIKE
// 'prefix%' — no full table scan. On a prefix miss it retries with lemma
// candidates ("parts" → prefix "part"), because ECDICT headwords are base
// forms; inflected input would otherwise return nothing (roadmap #786).
func (e *ECDict) Search(keyword string) ([]*model.KeyQueryIndex, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*model.KeyQueryIndex{}, nil
	}
	// typed form first, then lemma prefixes ("parts" → "part")
	prefixes := []string{keyword}
	seen := map[string]bool{keyword: true}
	for _, c := range lemmaCandidates(keyword) {
		if !seen[c] {
			seen[c] = true
			prefixes = append(prefixes, c)
		}
	}
	out := make([]*model.KeyQueryIndex, 0, 50)
	for _, p := range prefixes {
		rows, err := e.db.Query(
			`SELECT word FROM ecdict WHERE word LIKE ? ORDER BY frq ASC LIMIT 50`,
			p+"%",
		)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var w string
			if err := rows.Scan(&w); err != nil {
				rows.Close()
				return nil, err
			}
			out = append(out, &model.KeyQueryIndex{
				IndexType:         string(model.DictTypeECDICT),
				MdictKeyWordIndex: &model.MdictKeyWordIndex{KeyWord: w},
			})
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if len(out) > 0 {
			break // prefix hit (typed or lemma); stop — no further fallback
		}
	}
	return out, nil
}

// LookupResource — ECDICT has no external resources (css/images/fonts).
func (e *ECDict) LookupResource(keyword string) ([]byte, error) {
	return nil, fmt.Errorf("ecdict: no resources (%q)", keyword)
}

// Close releases the database handle.
func (e *ECDict) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.db == nil {
		return nil
	}
	err := e.db.Close()
	e.db = nil
	return err
}

// renderCard builds the entry HTML fragment (word + phonetic + CN translation +
// EN definition). WrapContent wraps it into the full page. Text is HTML-escaped;
// newlines become <br> — including the literal `\n` sequences ECDICT's CSV uses
// to encode in-field line breaks (stored verbatim by the importer).
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

// nl2br converts newlines to <br>. ECDICT's source CSV escapes in-field line
// breaks as a literal backslash-n (and backslashes as \\), and the importer
// stores those bytes verbatim — so real '\n' and the literal "\n" sequence
// must both become <br>. A literal "\\" un-escapes to a single backslash,
// which HTML-escaping happened to render correctly before, so leave it alone.
func nl2br(s string) string {
	s = strings.ReplaceAll(s, `\n`, "<br>")
	return strings.ReplaceAll(s, "\n", "<br>")
}
