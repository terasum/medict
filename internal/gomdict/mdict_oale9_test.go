package gomdict

import (
	"testing"
)

func TestOALE9(t *testing.T) {
	dict, err := New("testdata/mdx/testdict.mdx")
	if err != nil {
		t.Error(err)
	}

	err = dict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Digest:\n-----------------------------\n%s\n", dict.Digest())

	keywordEntries, err := dict.GetKeyWordEntries()
	if err != nil {
		t.Fatal(err)
	}
	for idx, entry := range keywordEntries {
		if idx > 10 {
			break
		}
		t.Logf("\n\n-----------\n\n")
		t.Logf("keyword: %s", entry.KeyWord)
		index, err := dict.KeywordEntryToIndex(entry)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("index: %+v", index)

		def, err := dict.LocateByKeywordIndex(index)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("def: %s", def)

	}

}
