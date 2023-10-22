package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/terasum/medict/pkg/model"
	"testing"
)

func TestMdict_Name(t *testing.T) {
	mdict := &Mdict{
		mdxFilePath:       "/User/yourname/test/dict/test.mdx",
		mddFilePaths:      nil,
		hasBuildIndex:     false,
		buildingIndexLock: nil,
	}
	t.Logf("name is %s", mdict.Name())
	assert.Equal(t, "test", mdict.Name())
}

func TestCreateSqliteIndex(t *testing.T) {

	mdict, err := NewMdict(&model.DirItem{
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
	t.Logf(mdict.Name())
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
