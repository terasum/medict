package gomdict

import (
	"os"
	"testing"
	"time"
)

func TestMdict_Lookup(t *testing.T) {
	mdict, err := New("testdata/dict/oale8.mdx")
	if err != nil {
		t.Fatal(err)
	}
	def, err := mdict.Lookup("about")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(string(def))
}

func TestMdict_Lookup2(t *testing.T) {
	mdict, err := New("testdata/dict/oale8.mdx")
	if err != nil {
		t.Fatal(err)
	}
	def, err := mdict.Lookup("definition")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(string(def))
}

func TestMdict_LookupMdd3(t *testing.T) {
	mdict, err := New("testdata/dict/oale8.mdd")
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
	mdict, err := New("testdata/dict/oale8.mdd")
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
	mdict, err := New("testdata/dict/oale8.mdd")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("testdata/out/olae8mdd.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range mdict.KeyBlockData.KeyEntries {
		file.Write([]byte(entry.KeyWord + "\n"))
	}

	file.Close()
}

func TestMdict_LookupMdd6(t *testing.T) {
	mdict, err := New("testdata/dict/ode3e.mdd")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("testdata/out/ode3emdd.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range mdict.KeyBlockData.KeyEntries {
		file.Write([]byte(entry.KeyWord + "\n"))
	}

	file.Close()
}

func TestMdict_SimSearch(t *testing.T) {
	mdict, err := New("testdata/dict/oale8.mdx")
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	err = mdict.BuildBKTree()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("building bktree costs %dms", time.Now().Sub(start).Milliseconds())

	result, err := mdict.SimSearch("hell", 2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
