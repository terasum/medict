package entry

import (
	_ "embed"
	"testing"
)

func TestWritePresetDictionary(t *testing.T) {
	err := WritePresetDictionary("testdata/out")
	if err != nil {
		t.Fatal(err)
	}
}
