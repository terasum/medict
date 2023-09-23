package stardict

import (
	"testing"
)

func TestReadDict(t *testing.T) {
	// init dictionary with path to dictionary files and name of dictionary
	dict, err := NewStarDict(&Option{
		DictDzFilePath: "testdata/stardict/eedic.pdb/eedic.pdb.dict.dz",
		DictFilePath:   "",
		IdxFilePath:    "testdata/stardict/eedic.pdb/eedic.pdb.idx",
		IfoFilePath:    "testdata/stardict/eedic.pdb/eedic.pdb.ifo",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = dict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	for _, k := range dict.KeyList() {
		t.Logf("key: %s", k)
	}

	def := dict.Lookup("oxygen") // get translations
	t.Logf("def:\n %s", def)

}
