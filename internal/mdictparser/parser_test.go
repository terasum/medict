package mdictparser

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestMdictParser_Lookup(t *testing.T) {
	parser := &MdictParser{}
	err := parser.Load("./testdict/testdict.mdx")
	assert.Nil(t, err)
	word := parser.Lookup("zoom")
	t.Log(word)
	parser.Destroy()
}

func TestMdictParser_AllWords(t *testing.T) {
	parser := &MdictParser{}
	err := parser.Load("./testdict/wlghyzd2000.mdx")
	assert.Nil(t, err)
	words, length, err := parser.AllWords()
	assert.Nil(t, err)
	t.Logf("%#v", words[0])
	assert.Equal(t, words[0].KeyWord, "a")
	assert.Equal(t, words[1].KeyWord, "ai")
	assert.Equal(t, words[2].KeyWord, "an")
	t.Log(length)
	assert.Equal(t, length, uint64(31402))
	parser.Destroy()
}

func TestMdictParser_AllWordsr2(t *testing.T) {
	parser := &MdictParser{}
	err := parser.Load("./testdict/wlghyzd2000.mdd")
	assert.Nil(t, err)
	words, length, err := parser.AllWords()
	t.Log(length)

	assert.Nil(t, err)
	for _, w := range words {
		t.Logf("%#v", w)
		def, err := parser.FindDef(w.KeyWord, w.RecordStart)
		assert.Nil(t, err)
		t.Logf(w.KeyWord)
		fname := strings.ReplaceAll(w.KeyWord, "\\", "/")
		bt, err := hex.DecodeString(def)
		assert.Nil(t, err)
		err = ioutil.WriteFile("./testdict/gendata/"+fname, bt, 0644)
		assert.Nil(t, err)
	}

	parser.Destroy()
}
