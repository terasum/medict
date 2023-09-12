package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terasum/medict/internal/config"
)

func TestDictService_Dicts(t *testing.T) {
	ds, err := NewDictService(&config.Config{
		BaseDictDir:      "./testdata/dicts",
		StaticServerPort: 0,
		CacheFileDir:     "",
	})
	assert.Nil(t, err)
	dicts := ds.Dicts()
	assert.Equal(t, 1, len(dicts))
	t.Logf("%#v", dicts[0])

}

func TestDictService_Dicts2(t *testing.T) {
	ds, err := NewDictService(&config.Config{
		BaseDictDir:      "../../testdict/",
		StaticServerPort: 0,
		CacheFileDir:     "",
	})
	assert.Nil(t, err)
	dicts := ds.Dicts()
	assert.Equal(t, 2, len(dicts))
	t.Logf("%#v", dicts[0])
}
