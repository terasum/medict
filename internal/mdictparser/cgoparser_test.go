package mdictparser

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestLoadAndDestory(t *testing.T) {
	dict, _ := LoadDict("./testdict/testdict.mdx")
	err := DestroyDict(dict)
	assert.Nil(t, err)
}

func TestLookUpDict(t *testing.T) {
	dict, _ := LoadDict("./testdict/testdict.mdx")
	result := LookUpDict(dict, "zoom")
	t.Logf(result)
	err := DestroyDict(dict)
	assert.Nil(t, err)
}

func TestAllKeyList(t *testing.T) {
	dict, _ := LoadDict("./testdict/testdict.mdx")
	result, length, err := AllKeyList(dict)
	assert.Nil(t, err)
	t.Logf(result[3].KeyWord)
	t.Logf(strconv.FormatUint(result[3].RecordStart, 10))
	t.Logf(strconv.FormatUint(length, 10))
	err = DestroyDict(dict)
	assert.Nil(t, err)
}

func TestParseDefinition(t *testing.T) {
	dict, _ := LoadDict("./testdict/testdict.mdx")
	result, length, err := AllKeyList(dict)
	assert.Nil(t, err)
	t.Logf(result[3].KeyWord)
	t.Logf(strconv.FormatUint(result[3].RecordStart, 10))
	t.Logf(strconv.FormatUint(length, 10))
	defResult := ParseDefinition(dict, result[3].KeyWord, result[3].RecordStart)
	t.Logf(defResult)
	err = DestroyDict(dict)
	assert.Nil(t, err)
}
