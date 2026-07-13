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

package model

import (
	"reflect"
	"sort"
	"testing"
)

// TestDictList_SortByName: the list is ordered by Name, case-insensitive
// ascending — not by ID (the MD5 of the dict dir), which was meaningless and
// jumped on every new dict (#740).
func TestDictList_SortByName(t *testing.T) {
	list := DictList{
		&PlainDictionaryItem{ID: "bb", Name: "OALD"},
		&PlainDictionaryItem{ID: "aa", Name: "apple dict"},
		&PlainDictionaryItem{ID: "cc", Name: "Cambridge"},
	}
	sort.Sort(list)

	got := []string{list[0].Name, list[1].Name, list[2].Name}
	want := []string{"apple dict", "Cambridge", "OALD"} // case-insensitive ascending
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}
