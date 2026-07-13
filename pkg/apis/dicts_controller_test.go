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

package apis

import (
	"testing"

	"github.com/terasum/medict/pkg/model"
)

// TestConvertKeyIndex locks the 8-param → struct conversion after the loop
// refactor (#738): every param lands in the right field, empty params default
// to 0, a bad integer errors, and dictType maps to IndexType.
func TestConvertKeyIndex(t *testing.T) {
	got, err := convertKeyIndex("medict", "5", "100", "200", "hello", "300", "400", "500", "600", "700")
	if err != nil {
		t.Fatalf("convertKeyIndex: %v", err)
	}
	if got.IndexType != model.IndexTypeMdict {
		t.Fatalf("IndexType = %q, want %q", got.IndexType, model.IndexTypeMdict)
	}
	if got.KeyWord != "hello" {
		t.Fatalf("KeyWord = %q", got.KeyWord)
	}
	checks := []struct {
		name string
		got  int64
		want int64
	}{
		{"ID", int64(got.ID), 5},
		{"RecordLocateStartOffset", got.RecordLocateStartOffset, 100},
		{"RecordLocateEndOffset", got.RecordLocateEndOffset, 200},
		{"RecordBlockDataStartOffset", got.RecordBlockDataStartOffset, 300},
		{"RecordBlockDataCompressSize", got.RecordBlockDataCompressSize, 400},
		{"RecordBlockDataDeCompressSize", got.RecordBlockDataDeCompressSize, 500},
		{"KeyWordDataStartOffset", got.KeyWordDataStartOffset, 600},
		{"KeyWordDataEndOffset", got.KeyWordDataEndOffset, 700},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Fatalf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}

	// Empty entryId/recordStart/recordEnd default to "0" → 0.
	got0, err := convertKeyIndex("medict", "", "", "", "w", "1", "2", "3", "4", "5")
	if err != nil {
		t.Fatalf("empty-default convertKeyIndex: %v", err)
	}
	if got0.ID != 0 || got0.RecordLocateStartOffset != 0 || got0.RecordLocateEndOffset != 0 {
		t.Fatalf("empty-default handling broke: %+v", got0)
	}

	// A non-integer must error (not panic).
	if _, err := convertKeyIndex("medict", "x", "1", "2", "w", "3", "4", "5", "6", "7"); err == nil {
		t.Fatal("expected error for non-integer entryId")
	}

	// dictType "stardict" → IndexTypeStardict.
	gotS, _ := convertKeyIndex("stardict", "1", "1", "2", "w", "3", "4", "5", "6", "7")
	if gotS.IndexType != model.IndexTypeStardict {
		t.Fatalf("IndexType for stardict = %q", gotS.IndexType)
	}
}
