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

package service

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/op/go-logging"
	// CGO-free pure-Go SQLite driver (database/sql).
	_ "modernc.org/sqlite"
)

var bLog = logging.MustGetLogger("bookmark")

// Notebook groups saved words into user-managed collections. Exactly one
// notebook has IsDefault=true at any time; the default cannot be deleted.
type Notebook struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	CreatedAt int64  `json:"created_at"`
}

// Bookmark is one saved word entry (#643). NotebookId ties it to a Notebook;
// the same word may live in several notebooks. The rendered HTML snapshot is
// stored separately (bookmarks.content_html) and fetched via GetSnapshot, so
// list responses don't ship large payloads.
type Bookmark struct {
	Word       string `json:"word"`
	DictId     string `json:"dict_id"`
	DictName   string `json:"dict_name"`
	NotebookId string `json:"notebook_id"`
	SavedAt    int64  `json:"saved_at"`
}

// ExportRow is one saved word bundled with its notebook name and stored HTML
// snapshot — everything the Anki exporter needs to build a card. Returned by
// ExportRows; unlike Bookmark it carries the (potentially large) snapshot HTML.
type ExportRow struct {
	Word         string `json:"word"`
	DictName     string `json:"dict_name"`
	NotebookName string `json:"notebook_name"`
	HTML         string `json:"html"`
	SavedAt      int64  `json:"saved_at"`
}

// BookmarkStore persists notebooks + saved words (with self-contained HTML
// snapshots) to a SQLite file (bookmarks.db) in the app config dir. Uses the
// pure-Go modernc.org/sqlite driver — no CGO. Single serialized connection
// (SetMaxOpenConns(1)) avoids SQLITE_BUSY; bookmark ops are infrequent.
type BookmarkStore struct {
	mu sync.Mutex
	db *sql.DB
}

// NewBookmarkStore opens (creating if absent) the SQLite store and ensures the
// schema + a default notebook exist. No JSON migration.
func NewBookmarkStore(configDir string) (*BookmarkStore, error) {
	dbPath := filepath.Join(configDir, "bookmarks.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open bookmark db: %w", err)
	}
	db.SetMaxOpenConns(1)
	s := &BookmarkStore{db: db}
	if err := s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *BookmarkStore) init() error {
	if _, err := s.db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		return fmt.Errorf("pragma wal: %w", err)
	}
	if _, err := s.db.Exec(`PRAGMA busy_timeout=5000;`); err != nil {
		return fmt.Errorf("pragma busy_timeout: %w", err)
	}
	schema := `
	CREATE TABLE IF NOT EXISTS notebooks (
		id         TEXT    PRIMARY KEY,
		name       TEXT    NOT NULL,
		is_default INTEGER NOT NULL DEFAULT 0,
		created_at INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS bookmarks (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		word         TEXT    NOT NULL,
		dict_id      TEXT    NOT NULL,
		dict_name    TEXT    NOT NULL,
		notebook_id  TEXT    NOT NULL,
		saved_at     INTEGER NOT NULL,
		content_html TEXT    NOT NULL DEFAULT ''
	);
	CREATE UNIQUE INDEX IF NOT EXISTS uq_bookmarks        ON bookmarks(word, dict_id, notebook_id);
	CREATE        INDEX IF NOT EXISTS idx_bookmarks_nb     ON bookmarks(notebook_id);
	`
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return s.ensureDefaultNotebook()
}

// newDefaultNotebook builds the auto-created default notebook.
func newDefaultNotebook() Notebook {
	return Notebook{
		Id:        uuid.NewString(),
		Name:      "默认生词本",
		IsDefault: true,
		CreatedAt: time.Now().Unix(),
	}
}

// ensureDefaultNotebook seeds a default notebook when none exist and collapses
// any duplicate-default state down to exactly one.
func (s *BookmarkStore) ensureDefaultNotebook() error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM notebooks`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		n := newDefaultNotebook()
		if _, err := s.db.Exec(
			`INSERT INTO notebooks(id, name, is_default, created_at) VALUES(?, ?, 1, ?)`,
			n.Id, n.Name, n.CreatedAt,
		); err != nil {
			return fmt.Errorf("seed default notebook: %w", err)
		}
		return nil
	}
	// pick the first existing default (else the first notebook) as the sole default
	if _, err := s.db.Exec(`
		UPDATE notebooks SET is_default = (id = COALESCE(
			(SELECT id FROM notebooks WHERE is_default = 1 ORDER BY created_at LIMIT 1),
			(SELECT id FROM notebooks ORDER BY created_at LIMIT 1)
		))`); err != nil {
		return fmt.Errorf("normalize default notebook: %w", err)
	}
	return nil
}

// defaultId returns the default notebook's id (empty if none).
func (s *BookmarkStore) defaultId() (string, error) {
	var id string
	err := s.db.QueryRow(`SELECT id FROM notebooks WHERE is_default = 1 LIMIT 1`).Scan(&id)
	if err == sql.ErrNoRows {
		err = s.db.QueryRow(`SELECT id FROM notebooks ORDER BY created_at LIMIT 1`).Scan(&id)
	}
	return id, err
}

// Add saves a word (+ its rendered HTML snapshot) into the given notebook.
// Deduped by word+dictId+notebookId (INSERT OR IGNORE → no-op if present).
// An empty notebookId lands in the default notebook.
func (s *BookmarkStore) Add(word, dictId, dictName, notebookId, html string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if notebookId == "" {
		id, err := s.defaultId()
		if err != nil {
			return err
		}
		notebookId = id
	}
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO bookmarks(word, dict_id, dict_name, notebook_id, saved_at, content_html)
		 VALUES(?, ?, ?, ?, ?, ?)`,
		word, dictId, dictName, notebookId, time.Now().Unix(), html,
	)
	return err
}

// Remove deletes a saved word from a specific notebook.
func (s *BookmarkStore) Remove(word, dictId, notebookId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(
		`DELETE FROM bookmarks WHERE word = ? AND dict_id = ? AND notebook_id = ?`,
		word, dictId, notebookId,
	)
	return err
}

// All returns all saved words, newest first (snapshot excluded).
func (s *BookmarkStore) All() ([]Bookmark, error) {
	rows, err := s.db.Query(
		`SELECT word, dict_id, dict_name, notebook_id, saved_at
		 FROM bookmarks ORDER BY saved_at DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Bookmark, 0)
	for rows.Next() {
		var b Bookmark
		if err := rows.Scan(&b.Word, &b.DictId, &b.DictName, &b.NotebookId, &b.SavedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// ExportRows returns saved words with their notebook name + stored HTML snapshot,
// newest first. notebookId == "" returns words from every notebook. Used by the
// Anki exporter; each row is one card, each notebook is one deck.
func (s *BookmarkStore) ExportRows(notebookId string) ([]ExportRow, error) {
	q := `SELECT b.word, b.dict_name, n.name, b.content_html, b.saved_at
	      FROM bookmarks b JOIN notebooks n ON n.id = b.notebook_id`
	var args []any
	if notebookId != "" {
		q += ` WHERE b.notebook_id = ?`
		args = append(args, notebookId)
	}
	q += ` ORDER BY b.saved_at DESC, b.id DESC`
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]ExportRow, 0)
	for rows.Next() {
		var r ExportRow
		if err := rows.Scan(&r.Word, &r.DictName, &r.NotebookName, &r.HTML, &r.SavedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// Has reports whether a word+dictId is already saved (in any notebook).
func (s *BookmarkStore) Has(word, dictId string) bool {
	var n int
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM bookmarks WHERE word = ? AND dict_id = ?`, word, dictId,
	).Scan(&n); err != nil {
		return false
	}
	return n > 0
}

// GetSnapshot returns the stored HTML snapshot for a bookmark ("" if absent).
func (s *BookmarkStore) GetSnapshot(word, dictId, notebookId string) (string, error) {
	var html string
	err := s.db.QueryRow(
		`SELECT content_html FROM bookmarks WHERE word = ? AND dict_id = ? AND notebook_id = ?`,
		word, dictId, notebookId,
	).Scan(&html)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return html, err
}

// Notebooks returns all notebooks, default first then by creation time.
func (s *BookmarkStore) Notebooks() ([]Notebook, error) {
	rows, err := s.db.Query(
		`SELECT id, name, is_default, created_at
		 FROM notebooks ORDER BY is_default DESC, created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Notebook, 0)
	for rows.Next() {
		var n Notebook
		var isDef int
		if err := rows.Scan(&n.Id, &n.Name, &isDef, &n.CreatedAt); err != nil {
			return nil, err
		}
		n.IsDefault = isDef == 1
		out = append(out, n)
	}
	return out, rows.Err()
}

// AddNotebook creates a new (non-default) notebook and returns it.
func (s *BookmarkStore) AddNotebook(name string) (Notebook, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "生词本"
	}
	n := Notebook{
		Id:        uuid.NewString(),
		Name:      name,
		IsDefault: false,
		CreatedAt: time.Now().Unix(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := s.db.Exec(
		`INSERT INTO notebooks(id, name, is_default, created_at) VALUES(?, ?, 0, ?)`,
		n.Id, n.Name, n.CreatedAt,
	); err != nil {
		return Notebook{}, err
	}
	return n, nil
}

// RenameNotebook renames a notebook by id.
func (s *BookmarkStore) RenameNotebook(id, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("生词本名称不能为空")
	}
	res, err := s.db.Exec(`UPDATE notebooks SET name = ? WHERE id = ?`, name, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("生词本不存在")
	}
	return nil
}

// SetDefaultNotebook marks the given notebook as the default (unsetting all others).
func (s *BookmarkStore) SetDefaultNotebook(id string) error {
	var exists int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM notebooks WHERE id = ?`, id).Scan(&exists); err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("生词本不存在")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE notebooks SET is_default = 0`); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`UPDATE notebooks SET is_default = 1 WHERE id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// RemoveNotebook deletes a non-default notebook and moves its words into the
// default notebook (no data loss). The default notebook cannot be removed.
func (s *BookmarkStore) RemoveNotebook(id string) error {
	var isDef int
	err := s.db.QueryRow(`SELECT is_default FROM notebooks WHERE id = ?`, id).Scan(&isDef)
	if err == sql.ErrNoRows {
		return errors.New("生词本不存在")
	}
	if err != nil {
		return err
	}
	if isDef == 1 {
		return errors.New("默认生词本不可删除")
	}
	var defId string
	if err := s.db.QueryRow(`SELECT id FROM notebooks WHERE is_default = 1 LIMIT 1`).Scan(&defId); err != nil {
		return errors.New("未找到默认生词本")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE bookmarks SET notebook_id = ? WHERE notebook_id = ?`, defId, id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`DELETE FROM notebooks WHERE id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Close releases the database handle.
func (s *BookmarkStore) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
