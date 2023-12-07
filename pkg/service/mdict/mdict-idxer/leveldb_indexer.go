package mdict_idxer

import (
	"encoding/json"
	"errors"
	"github.com/terasum/medict/pkg/model"
	lvdb "github.com/terasum/medict/pkg/service/mdict/leveldb-repo"
	"regexp"
	"strings"
)

var _ Indexer = &MedictDBIndexer{}

const prefixKeyword = "PFKW#_"
const prefixMeta = "PFMT#_"

var regex = regexp.MustCompile("/[., '\\\\@_\\$#\\%\\:\\/]/g")

type MedictDBIndexer struct {
	indexFileDirPath string
	lvdb             *lvdb.LvDB
	//searchTree       *searchTree
}

func strip(key string) string {
	return prefixKeyword + strings.TrimLeft(regex.ReplaceAllString(key, ""), prefixKeyword)
}

func NewIndexer(fpath string) (*MedictDBIndexer, error) {
	db, err := lvdb.NewLvDB(fpath)
	if err != nil {
		return nil, err
	}
	idxer := &MedictDBIndexer{
		indexFileDirPath: fpath,
		lvdb:             db,
	}

	//err = idxer.initSearchTree()
	//if err != nil {
	//	return nil, err
	//}
	return idxer, nil

}

func (m *MedictDBIndexer) Lookup(keyword string) (*model.MdictKeyWordIndex, error) {
	key := strip(keyword)
	data, err := m.lvdb.Get(key)
	if err != nil {
		return nil, err
	}
	indexData := new(model.MdictKeyWordIndex)
	err = json.Unmarshal(data, indexData)
	if err != nil {
		return nil, err
	}
	return indexData, nil
}

func (m *MedictDBIndexer) SetMeta(key, value string) error {
	key = prefixMeta + key
	return m.lvdb.Put(key, []byte(value))
}

func (m *MedictDBIndexer) GetMeta(key string) (string, error) {
	key = prefixMeta + key
	value, err := m.lvdb.Get(key)
	return string(value), err
}

func (m *MedictDBIndexer) AddRecord(record *model.MdictKeyWordIndex) (resErr error) {
	startTime := logstart("MedictDBIndexer.AddRecord", record)
	defer logend("MedictDBIndexer.AddRecord", startTime, resErr)

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	key := strip(record.KeyWord)
	err = m.lvdb.Put(key, data)
	if err != nil {
		return err
	}
	return nil
	//return m.searchTree.addKeyValue(key, data)
}

func (m *MedictDBIndexer) Search(keyword string) (res []*model.MdictKeyWordIndex, resErr error) {
	startTime := logstart("MedictDBIndexer.Search", keyword)
	defer logend("MedictDBIndexer.Search", startTime, resErr)

	keyword = strip(keyword)
	//values, err := m.searchTree.search(keyword)
	//if err != nil {
	//	return nil, err
	//}

	values, err := m.lvdb.Prefix(keyword)
	if err != nil {
		return nil, err
	}

	list := make([]*model.MdictKeyWordIndex, 0)
	for idx, value := range values {
		data, err1 := m.lvdb.Get(value)
		if err1 != nil {
			log.Errorf("query db failed, key: %s", value)
			continue
		}

		vi := new(model.MdictKeyWordIndex)
		err1 = json.Unmarshal(data, vi)
		if err1 != nil {
			log.Errorf("unmarshal value failed, %s, data: %s", err1.Error(), data)
			continue
		}
		vi.ID = idx
		list = append(list, vi)
	}

	if len(list) == 0 {
		return nil, errors.New("result not found")
	}

	return list, nil
}
