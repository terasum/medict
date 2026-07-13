package stardict

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/terasum/medict/pkg/model"
)

// TestStarDict_NotReadyAndResourceNotFound: a not-built StarDict surfaces
// ErrNotFound from Lookup, and LookupResource returns ErrNotFound (stardict
// resource loading isn't implemented — #739). No testdata required.
func TestStarDict_NotReadyAndResourceNotFound(t *testing.T) {
	s := &StarDict{ready: false}

	if _, err := s.Lookup("anything"); !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("Lookup not-ready: want ErrNotFound, got %v", err)
	}
	if _, err := s.LookupResource("anything"); !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("LookupResource: want ErrNotFound, got %v", err)
	}
}

func TestStarDict_Lookup(t *testing.T) {
	// stardict testdata is not checked into the repo; skip when absent so CI
	// stays green. The import cycle (pkg/service -> stardict) is avoided by
	// constructing the DictionaryItem locally instead of via service.NewByDirItem.
	if _, err := os.Stat("testdata/stardict/eedic.pdb.ifo"); err != nil {
		t.Skip("stardict testdata missing: testdata/stardict/eedic.pdb.ifo")
	}

	dirItem := &model.DirItem{
		DictType:           model.DictTypeStarDict,
		StarDictDzAbsPath:  "testdata/stardict/eedic.pdb.dict.dz",
		StarDictAbsPath:    "testdata/stardict/eedic.pdb.dict",
		StarDictIdxAbsPath: "testdata/stardict/eedic.pdb.idx",
		StarDictIfoAbsPath: "testdata/stardict/eedic.pdb.ifo",
	}

	dict, err := NewStardict(dirItem)
	if err != nil {
		t.Fatal(err)
	}
	err = dict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	// Construct DictionaryItem locally to avoid importing pkg/service, which
	// would create an import cycle (pkg/service imports stardict).
	item := &model.DictionaryItem{
		PlainDictionaryItem: &model.PlainDictionaryItem{
			Name:        dict.Name(),
			DictType:    string(dirItem.DictType),
			Description: dict.Description(),
		},
		PathInfo: dirItem,
		Dict: dict,
	}
	t.Logf("%+v", item.ToPlain())
	t.Logf("%+v", item.Name)
	words, err := item.Dict.Search("impair")
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(words, "", " ")
	t.Logf("words: %s", data)
}
