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

package go_mdict

import (
	"os"
	"testing"
)

func TestMdict_Lookup(t *testing.T) {
	skipIfMissingTestData(t, "testdata/dict/oale8.mdx")
	mdict, err := New("testdata/dict/oale8.mdx")
	if err != nil {
		t.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	def, err := mdict.Lookup("about")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", string(def))
}

func TestMdict_Lookup2(t *testing.T) {
	skipIfMissingTestData(t, "testdata/dict/oale8.mdx")
	mdict, err := New("testdata/dict/oale8.mdx")
	if err != nil {
		t.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}
	def, err := mdict.Lookup("definition")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", string(def))
}

func TestMdict_LookupMdd3(t *testing.T) {
	skipIfMissingTestData(t, "testdata/dict/oale8.mdd")
	mdict, err := New("testdata/dict/oale8.mdd")
	if err != nil {
		t.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	def, err := mdict.Lookup("\\uk\\double_quick_1_gb_1.mp3")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("testdata/out/double_quick_1_gb_1.mp3")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(def)
	if err != nil {
		t.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMdict_LookupMdd4(t *testing.T) {
	skipIfMissingTestData(t, "testdata/dict/oale8.mdd")
	mdict, err := New("testdata/dict/oale8.mdd")
	if err != nil {
		t.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	def, err := mdict.Lookup("\\uk\\abalone__gb_1.mp3")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("testdata/out/abalone__gb_1.mp3")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(def)
	if err != nil {
		t.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMdict_LookupMdd5(t *testing.T) {
	skipIfMissingTestData(t, "testdata/dict/oale8.mdd")
	mdict, err := New("testdata/dict/oale8.mdd")
	if err != nil {
		t.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("testdata/out/olae8mdd.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range mdict.keyBlockData.keyEntries {
		file.Write([]byte(entry.KeyWord + "\n"))
	}

	file.Close()

	pngbytes, err := mdict.Lookup("\\us_pron.png")
	if err != nil {
		t.Fatal(err)
	}
	file, err = os.Create("testdata/out/us_pron.png")
	if err != nil {
		t.Fatal(err)
	}

	file.Write(pngbytes)
}

func TestMdict_LookupMdd6(t *testing.T) {
	skipIfMissingTestData(t, "testdata/dict/ode3e.mdd")
	mdict, err := New("testdata/dict/ode3e.mdd")
	if err != nil {
		t.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("testdata/out/ode3emdd.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range mdict.keyBlockData.keyEntries {
		file.Write([]byte(entry.KeyWord + "\n"))
	}

	file.Close()
}

// TestMdict_Name: regression for #259 — Name() used TrimRight(cutset) which
// ate trailing m/d/x chars ("mad.mdx" -> "ma"). TrimSuffix of the extension is
// correct for both .mdx and .mdd.
func TestMdict_Name(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/dicts/mad.mdx", "mad"},
		{"/dicts/OALD10.mdx", "OALD10"},
		{"/dicts/foox.mdd", "foox"},
		{"/dicts/标题词典.mdx", "标题词典"},
	}
	for _, tc := range cases {
		m := &Mdict{MdictBase: &MdictBase{filePath: tc.path}}
		if got := m.Name(); got != tc.want {
			t.Fatalf("Name(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}
