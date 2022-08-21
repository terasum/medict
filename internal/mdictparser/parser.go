package mdictparser

import "github.com/terasum/medict/internal/mdictparser/def"

type MdictParser struct {
	dictWrapper *RawDict
}

func NewMDictParser() *MdictParser {
	return &MdictParser{}
}

func (m *MdictParser) Load(dictFile string) error {
	// TODO exists?
	dict, err := LoadDict(dictFile)
	m.dictWrapper = dict
	if err != nil {
		return err
	}
	return nil
}

func (m *MdictParser) Type() string {
	return DictType(m.dictWrapper)
}

func (m *MdictParser) Lookup(word string) string {
	if word == "" {
		return ""
	}
	wordDef := LookUpDict(m.dictWrapper, word)
	return wordDef
}

func (m *MdictParser) FindDef(word string, recordStart uint64) (string, error) {
	definition := ParseDefinition(m.dictWrapper, word, recordStart)
	return definition, nil
}

func (m *MdictParser) AllWords() ([]*def.KeyItem, uint64, error) {
	return AllKeyList(m.dictWrapper)
}

func (m *MdictParser) Destroy() {
	_ = DestroyDict(m.dictWrapper)
}
