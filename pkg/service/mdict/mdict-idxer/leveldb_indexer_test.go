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

func TestLookup(t *testing.T) {

	idxer, err := NewIndexer("./testdata/testleveldb")
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

// Regression for issue #722 P0: strip() used strings.TrimLeft(cutset) where it
// meant strings.TrimPrefix, which trimmed any leading P/F/K/W/#/_ off the
// keyword. That collapsed distinct words onto one key ("King" and "ing" both =>
// "PFKW#_ing", one overwriting the other) and made prefix searches for P/F/K/W
// -initial words return the wrong cluster ("Kin" scanned "PFKW#_in").
func TestStrip_NoCollisionForPrefixInitialWords(t *testing.T) {
	idxer, err := NewIndexer(filepath.Join(t.TempDir(), "testleveldb"))
	if err != nil {
		t.Fatal(err)
	}

	add := func(kw string) {
		if err := idxer.AddRecord(&model.MdictKeyWordIndex{KeyWord: kw}); err != nil {
			t.Fatalf("AddRecord(%q): %v", kw, err)
		}
	}
	// "King" and "ing": the buggy strip() mapped both to "PFKW#_ing".
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
	// cluster the buggy strip() produced (it scanned "PFKW#_in").
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
