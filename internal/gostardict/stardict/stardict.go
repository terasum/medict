package stardict

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Definition contains translation items
type Definition struct {
	Parts []*RecordItem
}

// RecordItem contain single translation item
type RecordItem struct {
	Type rune
	Data []byte
}

type Option struct {
	DictDzFilePath string `json:"dictDzFilePath"`
	DictFilePath   string `json:"dictFilePath"`
	IdxFilePath    string `json:"idxFilePath"`
	IfoFilePath    string `json:"infoFilePath"`
}

// StarDict stardict dictionary
type StarDict struct {
	dictFilePath string
	idxFilePath  string
	infoFilePath string

	ready bool

	dict *RawDict
	idx  *RawIdx
	info *RawInfo
}

func (sdict *StarDict) KeyList() []string {
	if !sdict.ready {
		return []string{}
	}
	keys := make([]string, 0)
	for k, _ := range sdict.idx.items {
		keys = append(keys, k)
	}
	return keys
}

func (sdict *StarDict) Lookup(word string) string {
	if !sdict.ready {
		return ""
	}
	senses := sdict.locate(word) // get translations

	def := `<div class="start-dict-definition">`
	if len(senses) == 0 {
		def += `<section class="definition-section"> <p class="definition-section-part"> Nothing Found</p></section></div>`
		return def
	}

	for _, seq := range senses { // for each translation analyze returned parts
		section := `<section class="definition-section">`
		for _, p := range seq.Parts { // write each part contents to user
			part := `<p class="definition-section-part">`
			part = part + strings.ReplaceAll(string(p.Data), "\n", "<br/>")
			part += `</p>`
			section += part
		}
		section += `</section>`
		def += section
	}
	def += `</div>`
	return def
}

// locate translates given item
func (sdict *StarDict) locate(item string) (items []*Definition) {
	if !sdict.ready {
		return []*Definition{}
	}
	senses := sdict.idx.Get(item)

	for _, seq := range senses {
		sense := sdict.dict.GetSequence(seq.Offset, seq.Size)

		var transItems []*RecordItem

		if _, ok := sdict.info.Options["sametypesequence"]; ok {
			transItems = sdict.translateWithSametypesequence(sense)
		} else {
			transItems = sdict.translateWithoutSameTypeSequence(sense)
		}

		items = append(items, &Definition{Parts: transItems})
	}

	return
}

func (sdict *StarDict) translateWithSametypesequence(data []byte) (items []*RecordItem) {
	if !sdict.ready {
		return []*RecordItem{}
	}
	seq := sdict.info.Options["sametypesequence"]

	seqLen := len(seq)

	var dataPos int
	dataSize := len(data)

	for i, t := range seq {
		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			// if last seq item
			if i == seqLen-1 {
				items = append(items, &RecordItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				end := bytes.IndexRune(data[dataPos:], '\000')
				items = append(items, &RecordItem{Type: t, Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			if i == seqLen-1 {
				items = append(items, &RecordItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
				items = append(items, &RecordItem{Type: t, Data: data[dataPos+4 : dataPos+int(size)+5]})
				dataPos += int(size) + 5
			}
		}
	}

	return
}

func (sdict *StarDict) translateWithoutSameTypeSequence(data []byte) (items []*RecordItem) {
	if !sdict.ready {
		return []*RecordItem{}
	}
	var dataPos int
	dataSize := len(data)

	for {
		t := data[dataPos]

		dataPos++

		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			end := bytes.IndexRune(data[dataPos:], '\000')

			if end < 0 { // last item
				items = append(items, &RecordItem{Type: rune(t), Data: data[dataPos:dataSize]})
				dataPos = dataSize
			} else {
				items = append(items, &RecordItem{Type: rune(t), Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
			items = append(items, &RecordItem{Type: rune(t), Data: data[dataPos+4 : dataPos+int(size)+5]})
			dataPos += int(size) + 5
		}

		if dataPos >= dataSize {
			break
		}
	}

	return
}

// GetBookName returns book name
func (sdict *StarDict) GetBookName() string {
	if !sdict.ready {
		_, fname := filepath.Split(sdict.idxFilePath)
		return fname
	}
	return sdict.info.Options["bookname"]
}

// GetWordCount returns number of words in the dictionary
func (sdict *StarDict) GetWordCount() uint64 {
	if !sdict.ready {
		return 0
	}
	num, _ := strconv.ParseUint(sdict.info.Options["wordcount"], 10, 64)
	return num
}

func (sdict *StarDict) BuildIndex() error {
	if sdict.ready {
		return nil
	}
	info, err := readInfo(sdict.infoFilePath)
	if err != nil {
		return err
	}
	sdict.info = info

	idx, err := readIndex(sdict.idxFilePath, info)
	if err != nil {
		return err
	}

	sdict.idx = idx

	dict, err := readDict(sdict.dictFilePath, info)
	if err != nil {
		return err
	}

	sdict.dict = dict
	sdict.ready = true

	return nil
}

// NewStarDict returns a new StarDict
// path - path to dictionary files
// name - name of dictionary to parse
func NewStarDict(opt *Option) (*StarDict, error) {

	dictDzPath, _ := filepath.Abs(opt.DictDzFilePath)
	dictPath, _ := filepath.Abs(opt.DictFilePath)
	idxPath, _ := filepath.Abs(opt.IdxFilePath)
	infoPath, _ := filepath.Abs(opt.IfoFilePath)

	if !fileExists(infoPath) {
		return nil, fmt.Errorf("info file not exists, %s\n", infoPath)
	}

	if !fileExists(idxPath) {
		return nil, fmt.Errorf("index file not exists, %s\n", idxPath)
	}

	// we should have either .dict.dz or .dict file
	if !fileExists(dictDzPath) {
		if !fileExists(dictPath) {
			return nil, fmt.Errorf("dict & dict.dz file not exists, %s\n", idxPath)
		}
	} else {
		dictPath = dictDzPath
	}

	return &StarDict{
		dictFilePath: dictPath,
		idxFilePath:  idxPath,
		infoFilePath: infoPath,
	}, nil

}

func fileExists(fpath string) bool {
	if _, err := os.Stat(fpath); err == nil {
		return true
	}
	return false
}
