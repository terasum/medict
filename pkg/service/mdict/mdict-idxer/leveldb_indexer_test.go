package mdict_idxer

import (
	"github.com/terasum/medict/pkg/model"
	"testing"
)

func TestLookup(t *testing.T) {

	idxer, err := NewIndexer("./testdata/testleveldb")
	if err != nil {
		t.Fatal(err)
	}
	err = idxer.AddRecord(&model.MdictKeyWordIndex{
		ID:      0,
		KeyWord: "hell"})
	if err != nil {
		t.Fatal(err)
	}
	err = idxer.AddRecord(&model.MdictKeyWordIndex{
		ID:      0,
		KeyWord: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	err = idxer.AddRecord(&model.MdictKeyWordIndex{
		ID:      0,
		KeyWord: "helium"})
	if err != nil {
		t.Fatal(err)
	}

	list, err := idxer.Search("hell")
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range list {
		t.Logf("search word: %s", w.KeyWord)
	}
}
