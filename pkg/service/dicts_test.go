package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/terasum/medict/internal/config"
	"testing"
)

func TestDictService_Dicts(t *testing.T) {
	ds, err := NewDictService(&config.Config{
		BaseDictDir:      "./testdata/dicts",
		StaticServerPort: 0,
		APIServerPort:    0,
		CacheFileDir:     "",
	})
	assert.Nil(t, err)
	defer ds.GC()
	dicts := ds.Dicts()
	assert.Equal(t, 1, len(dicts))
	t.Logf("%#v", dicts[0])

}

func TestDictService_Dicts2(t *testing.T) {
	ds, err := NewDictService(&config.Config{
		BaseDictDir:      "../../testdict/",
		StaticServerPort: 0,
		APIServerPort:    0,
		CacheFileDir:     "",
	})
	assert.Nil(t, err)
	defer ds.GC()
	dicts := ds.Dicts()
	assert.Equal(t, 2, len(dicts))
	t.Logf("%#v", dicts[0])
}

func TestDictService_AllWordIndexing(t *testing.T) {
	ds, err := NewDictService(&config.Config{
		BaseDictDir:      "../../testdict/",
		StaticServerPort: 0,
		APIServerPort:    0,
		CacheFileDir:     "",
	})
	assert.Nil(t, err)
	defer ds.GC()
	mdxsize, mddsize, err := ds.AllWordIndexing()
	assert.Nil(t, err)
	t.Logf("mddsize %d, mddsize %d", mdxsize, mddsize)

}

func TestDictService_SimWords(t *testing.T) {
	ds, err := NewDictService(&config.Config{
		BaseDictDir:      "../../testdict/",
		StaticServerPort: 0,
		APIServerPort:    0,
		CacheFileDir:     "",
	})
	assert.Nil(t, err)
	defer ds.GC()
	mdxsize, mddsize, err := ds.AllWordIndexing()
	assert.Nil(t, err)
	t.Logf("mddsize %d, mddsize %d", mdxsize, mddsize)

	words, err := ds.SimWords("一鸣惊人")
	assert.Nil(t, err)
	t.Logf("simwords len %d", len(words))
	for _, w := range words {
		t.Log(w.KeyWord, w.RecordStart)
	}
}
