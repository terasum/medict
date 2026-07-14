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
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/op/go-logging"
)

var bLog = logging.MustGetLogger("bookmark")

// Bookmark is one saved word entry (#643).
type Bookmark struct {
	Word     string `json:"word"`
	DictId   string `json:"dict_id"`
	DictName string `json:"dict_name"`
	SavedAt  int64  `json:"saved_at"`
}

// BookmarkStore persists saved words to a JSON file in the app config dir.
type BookmarkStore struct {
	mu    sync.Mutex
	items []Bookmark
	path  string
}

func NewBookmarkStore(configDir string) *BookmarkStore {
	path := filepath.Join(configDir, "bookmarks.json")
	s := &BookmarkStore{path: path}
	s.load()
	return s
}

func (s *BookmarkStore) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		// first run / no file yet — start empty
		return
	}
	_ = json.Unmarshal(data, &s.items) // corrupt file → start empty
}

func (s *BookmarkStore) persist() {
	data, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		bLog.Errorf("persist: marshal failed: %s", err)
		return
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		bLog.Errorf("persist: write failed: %s", err)
	}
}

// Add saves a word (deduplicated by word+dictId). No-op if already saved.
func (s *BookmarkStore) Add(word, dictId, dictName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range s.items {
		if b.Word == word && b.DictId == dictId {
			return // already saved
		}
	}
	s.items = append(s.items, Bookmark{
		Word:     word,
		DictId:   dictId,
		DictName: dictName,
		SavedAt:  time.Now().Unix(),
	})
	s.persist()
}

// Remove deletes a saved word by word+dictId.
func (s *BookmarkStore) Remove(word, dictId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, b := range s.items {
		if b.Word == word && b.DictId == dictId {
			s.items = append(s.items[:i], s.items[i+1:]...)
			s.persist()
			return
		}
	}
}

// All returns all saved words, newest first.
func (s *BookmarkStore) All() []Bookmark {
	s.mu.Lock()
	defer s.mu.Unlock()
	// return a copy in reverse (newest first)
	n := len(s.items)
	out := make([]Bookmark, n)
	for i, b := range s.items {
		out[n-1-i] = b
	}
	return out
}

// Has reports whether a word+dictId is already saved.
func (s *BookmarkStore) Has(word, dictId string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range s.items {
		if b.Word == word && b.DictId == dictId {
			return true
		}
	}
	return false
}
