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

package support

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func touchFile(t *testing.T, root, rel string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
}

// TestWalkDir_NestedNoDup: a category folder holding two dicts must find each
// dict exactly once and must NOT register the category folder itself as a dict.
// Regression for #257 (recursive auto-scan previously double-counted nested
// dicts and mis-attributed them to the category dir).
func TestWalkDir_NestedNoDup(t *testing.T) {
	root := t.TempDir()
	touchFile(t, root, "English/OALD9/oald9.mdx")
	touchFile(t, root, "English/OALD9/oald9.mdd")
	touchFile(t, root, "English/COBUILD/cobuild.mdx")
	touchFile(t, root, "cc-cedict/cc-cedict.mdx")
	touchFile(t, root, "stardict1/stardict.ifo")
	touchFile(t, root, "stardict1/stardict.idx")
	touchFile(t, root, "stardict1/stardict.dict.dz")

	items, err := WalkDir(root)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, 0, len(items))
	for _, it := range items {
		got = append(got, filepath.Base(it.CurrentDir))
	}
	sort.Strings(got)

	want := []string{"COBUILD", "OALD9", "cc-cedict", "stardict1"}
	if len(got) != len(want) {
		t.Fatalf("found %d dicts %v, want %d %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("dict[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
	count := func(name string) int {
		c := 0
		for _, n := range got {
			if n == name {
				c++
			}
		}
		return c
	}
	if count("English") != 0 {
		t.Errorf("category dir English must not be a dict: %v", got)
	}
	if count("OALD9") != 1 {
		t.Errorf("OALD9 must appear once, got %d: %v", count("OALD9"), got)
	}
}

// TestWalkDir_FlatInTemp: classic flat layout via temp dir (one dict dir → one dict).
func TestWalkDir_FlatInTemp(t *testing.T) {
	root := t.TempDir()
	touchFile(t, root, "oale3/test.mdx")
	touchFile(t, root, "oale3/test.mdd")
	touchFile(t, root, "oale3/test.png") // cover candidate same basename as mdx
	items, err := WalkDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || filepath.Base(items[0].CurrentDir) != "oale3" {
		t.Fatalf("expected one dict oale3, got %+v", items)
	}
	if items[0].MdictMdxAbsPath == "" || len(items[0].MdictMddAbsPath) != 1 {
		t.Errorf("mdx/mdd not captured: %+v", items[0])
	}
	if items[0].CoverImgType != "png" || items[0].CoverImgPath == "" {
		t.Errorf("cover png (same basename as mdx) not captured: %+v", items[0])
	}
}
