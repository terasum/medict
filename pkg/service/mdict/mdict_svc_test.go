package mdict

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terasum/medict/pkg/model"
)

func TestMdict_Name(t *testing.T) {
	// An empty mdictSvcImpl{} has a nil mdx holder, so Name() panics on the
	// nil md.mdx.rawdict dereference inside Name(). Constructing a valid
	// holder requires a real .mdx file, which is not checked into the repo.
	if _, err := os.Stat("testdata/mdict/testdict.mdx"); err != nil {
		t.Skip("requires real mdx testdata: testdata/mdict/testdict.mdx")
	}
	mdict, err := NewMdictSvc(&model.DirItem{
		DictType:         model.DictTypeMdict,
		MdictMdxAbsPath:  "testdata/mdict/testdict.mdx",
		MdictMddAbsPath:  []string{"testdata/mdict/testdict.mdd"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("name is %s", mdict.Name())
	assert.Equal(t, "test", mdict.Name())
}

func TestCreateSqliteIndex(t *testing.T) {
	if _, err := os.Stat("testdata/mdict/testdict.mdx"); err != nil {
		t.Skip("requires real mdx testdata: testdata/mdict/testdict.mdx")
	}

	mdict, err := NewMdictSvc(&model.DirItem{
		BaseDir:            "testdata",
		CurrentDir:         "testdata/mdict",
		IsValid:            true,
		DictType:           model.DictTypeMdict,
		CoverImgPath:       "",
		CoverImgType:       "",
		ConfigPath:         "",
		LicensePath:        "",
		MdictMdxFileName:   "testdict",
		MdictMdxAbsPath:    "testdata/mdict/testdict.mdx",
		MdictMddAbsPath:    []string{"testdata/mdict/testdict.mdd"},
		StarDictDzAbsPath:  "",
		StarDictAbsPath:    "",
		StarDictIdxAbsPath: "",
		StarDictIfoAbsPath: "",
	})

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", mdict.Name())
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	records, err := mdict.Search("hello")
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		data, _ := json.Marshal(record)
		t.Logf("%s", data)
		def, err := mdict.Locate(record)
		t.Logf("def: %s, err: %v", def, err)
	}

}
