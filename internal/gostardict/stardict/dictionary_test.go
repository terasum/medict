package stardict

import (
	"testing"
)

func TestReadDict(t *testing.T) {

	// init dictionary with path to dictionary files and name of dictionary
	dict, err := NewDictionary("./testdata/stardict/elements-2.4.2", "elements")

	if err != nil {
		panic(err)
	}

	for _, k := range dict.KeyList() {
		t.Logf("key: %s", k)
	}

	def := dict.Lookup("oxygen") // get translations
	t.Logf("def:\n %s", def)

}
