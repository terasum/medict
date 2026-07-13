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

package mdict_idxer

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/terasum/medict/pkg/model"
)

// TestIndexer_CloseReleasesHandle verifies Close releases the leveldb handle;
// operations after Close return an error (not a panic). Guards the Close
// lifecycle added for issue #722.
func TestIndexer_CloseReleasesHandle(t *testing.T) {
	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "idxclose"))
	if err != nil {
		t.Fatal(err)
	}
	if err := idxer.AddRecord(&model.MdictKeyWordIndex{KeyWord: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := idxer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	// After Close, operations must error (leveldb ErrClosed), not panic.
	if _, err := idxer.Lookup("x"); err == nil {
		t.Fatal("Lookup after Close: want error, got nil")
	}
}

func TestLookup(t *testing.T) {

	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "testleveldb"))
	if err != nil {
		t.Fatal(err)
	}
	err = idxer.AddRecord(&model.MdictKeyWordIndex{
		ID:      0,
		KeyWord: "hell"})
	if err != nil {
		t.Fatal(err)
	}
	err = idxer.AddRecord(&model.MdictKeyWordIndex{
		ID:      0,
		KeyWord: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	err = idxer.AddRecord(&model.MdictKeyWordIndex{
		ID:      0,
		KeyWord: "helium"})
	if err != nil {
		t.Fatal(err)
	}

	list, err := idxer.Search("hell")
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range list {
		t.Logf("search word: %s", w.KeyWord)
	}
}

// TestDuplicateKeywords_NotCollapsed: same-headword entries (homographs) are
// kept distinct via the offset-suffix key. Lookup returns the first sense;
// Search surfaces all senses plus any longer matches (issue #722 #1).
// On the legacy headword-only key this would collapse to one "bank".
func TestDuplicateKeywords_NotCollapsed(t *testing.T) {
	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "dup"))
	if err != nil {
		t.Fatal(err)
	}
	if err := idxer.AddRecords([]*model.MdictKeyWordIndex{
		{KeyWord: "bank", RecordLocateStartOffset: 100}, // sense 1
		{KeyWord: "bank", RecordLocateStartOffset: 200}, // sense 2 (same headword)
		{KeyWord: "bankrupt", RecordLocateStartOffset: 300},
	}); err != nil {
		t.Fatal(err)
	}

	// Lookup returns the first sense of "bank" (offset 100).
	got, err := idxer.Lookup("bank")
	if err != nil {
		t.Fatalf("Lookup(bank): %v", err)
	}
	if got == nil || got.KeyWord != "bank" || got.RecordLocateStartOffset != 100 {
		t.Fatalf("Lookup(bank) = %+v, want first sense (offset 100)", got)
	}

	// Search("bank") surfaces BOTH senses (100, 200) plus the longer "bankrupt".
	res, err := idxer.Search("bank")
	if err != nil {
		t.Fatalf("Search(bank): %v", err)
	}
	seen := map[int64]bool{}
	for _, r := range res {
		seen[r.RecordLocateStartOffset] = true
	}
	for _, want := range []int64{100, 200, 300} {
		if !seen[want] {
			t.Fatalf("Search(bank) missing offset %d (homograph collapsed?); got %v", want, seen)
		}
	}
}

// TestAllRecords: AllRecords returns every keyword record (and excludes the
// PFMT#_ meta keys). Used by the holder to build the fuzzy BK-tree without
// re-parsing the source dict (issue #722 #2).
func TestAllRecords(t *testing.T) {
	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "all"))
	if err != nil {
		t.Fatal(err)
	}
	if err := idxer.SetMeta("Title", "demo"); err != nil {
		t.Fatal(err)
	}
	if err := idxer.AddRecords([]*model.MdictKeyWordIndex{
		{KeyWord: "apple", RecordLocateStartOffset: 10},
		{KeyWord: "banana", RecordLocateStartOffset: 20},
		{KeyWord: "cherry", RecordLocateStartOffset: 30},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := idxer.AllRecords()
	if err != nil {
		t.Fatalf("AllRecords: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("AllRecords returned %d records, want 3 (meta must be excluded)", len(got))
	}
	for _, r := range got {
		if r.KeyWord == "" {
			t.Fatalf("AllRecords returned a record with empty keyword: %+v", r)
		}
	}
}

// TestAddRecords_BatchEquivalence: AddRecords (leveldb batch) writes records
// that are individually lookable-up and carry their fields through — i.e. the
// batch path produces the same observable state as per-record AddRecord.
func TestAddRecords_BatchEquivalence(t *testing.T) {
	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "batch"))
	if err != nil {
		t.Fatal(err)
	}
	words := []string{"apple", "apply", "banana", "cherry"}
	records := make([]*model.MdictKeyWordIndex, len(words))
	for i, w := range words {
		records[i] = &model.MdictKeyWordIndex{KeyWord: w, RecordLocateStartOffset: int64(i * 100)}
	}
	if err := idxer.AddRecords(records); err != nil {
		t.Fatalf("AddRecords: %v", err)
	}
	for i, w := range words {
		got, err := idxer.Lookup(w)
		if err != nil {
			t.Fatalf("Lookup(%s): %v", w, err)
		}
		if got.KeyWord != w {
			t.Fatalf("Lookup(%s) keyword = %q", w, got.KeyWord)
		}
		if got.RecordLocateStartOffset != int64(i*100) {
			t.Fatalf("Lookup(%s) offset = %d, want %d", w, got.RecordLocateStartOffset, i*100)
		}
	}
	res, err := idxer.Search("app")
	if err != nil {
		t.Fatalf("Search(app): %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("Search(app) returned %d, want 2 (apple, apply)", len(res))
	}
}

// TestSearch_EmptyResultIsNotError: a prefix that matches nothing returns
// (nil, nil), not an error (issue #722 P3 error-as-control-flow).
func TestSearch_EmptyResultIsNotError(t *testing.T) {
	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "empty"))
	if err != nil {
		t.Fatal(err)
	}
	if err := idxer.AddRecords([]*model.MdictKeyWordIndex{{KeyWord: "hello"}}); err != nil {
		t.Fatal(err)
	}
	res, err := idxer.Search("zzz")
	if err != nil {
		t.Fatalf("Search(zzz) error: %v (issue #722 P3: empty != error)", err)
	}
	if len(res) != 0 {
		t.Fatalf("Search(zzz) = %v, want empty", res)
	}
}

// Regression for issue #722 P0: keywordKey() used strings.TrimLeft(cutset) where it
// meant strings.TrimPrefix, which trimmed any leading P/F/K/W/#/_ off the
// keyword. That collapsed distinct words onto one key ("King" and "ing" both =>
// "PFKW#_ing", one overwriting the other) and made prefix searches for P/F/K/W
// -initial words return the wrong cluster ("Kin" scanned "PFKW#_in").
func TestKeywordKey_NoCollisionForPrefixInitialWords(t *testing.T) {
	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "testleveldb"))
	if err != nil {
		t.Fatal(err)
	}

	add := func(kw string) {
		if err := idxer.AddRecord(&model.MdictKeyWordIndex{KeyWord: kw}); err != nil {
			t.Fatalf("AddRecord(%q): %v", kw, err)
		}
	}
	// "King" and "ing": the buggy keywordKey() mapped both to "PFKW#_ing".
	add("King")
	add("ing")

	for _, kw := range []string{"King", "ing"} {
		got, err := idxer.Lookup(kw)
		if err != nil {
			t.Fatalf("Lookup(%q) error: %v (strip collision?)", kw, err)
		}
		if got.KeyWord != kw {
			t.Fatalf("Lookup(%q) returned %q (strip collision overwrote it)", kw, got.KeyWord)
		}
	}

	// Prefix search "Kin" must return only King-prefixed words, never the "ing"
	// cluster the buggy keywordKey() produced (it scanned "PFKW#_in").
	res, err := idxer.Search("Kin")
	if err != nil {
		t.Fatalf("Search(Kin) error: %v", err)
	}
	for _, r := range res {
		if !strings.HasPrefix(r.KeyWord, "Kin") {
			t.Fatalf("Search(Kin) returned %q — strip TrimLeft bug leaked the wrong cluster", r.KeyWord)
		}
	}
}
