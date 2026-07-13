package mdict_idxer

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/terasum/medict/pkg/model"
	lvdb "github.com/terasum/medict/pkg/service/mdict/leveldb-repo"
	"strings"
)

var _ Indexer = &MedictDBIndexer{}

const prefixKeyword = "PFKW#_"
const prefixMeta = "PFMT#_"

type MedictDBIndexer struct {
	indexFileDirPath string
	lvdb             *lvdb.LvDB
	//searchTree       *searchTree
}

// strip namespaces a keyword under the PFKW#_ prefix so keyword keys never
// collide with the PFMT#_-prefixed meta keys. TrimPrefix (not TrimLeft, which
// treats its arg as a *cutset* and would silently eat leading P/F/K/W/#/_
// chars off the keyword itself — see issue #722) avoids double-prefixing.
func strip(key string) string {
	return prefixKeyword + strings.TrimPrefix(key, prefixKeyword)
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

// Close releases the underlying leveldb handle.
func (m *MedictDBIndexer) Close() error {
	if m.lvdb == nil {
		return nil
	}
	return m.lvdb.Close()
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

func (m *MedictDBIndexer) AddRecords(records []*model.MdictKeyWordIndex) (resErr error) {
	startTime := logstart("MedictDBIndexer.AddRecords", len(records))
	defer logend("MedictDBIndexer.AddRecords", startTime, resErr)

	if len(records) == 0 {
		return nil
	}
	batch := new(leveldb.Batch)
	for _, r := range records {
		data, err := json.Marshal(r)
		if err != nil {
			return err
		}
		batch.Put([]byte(strip(r.KeyWord)), data)
	}
	return m.lvdb.Write(batch)
}

func (m *MedictDBIndexer) Search(keyword string) (res []*model.MdictKeyWordIndex, resErr error) {
	startTime := logstart("MedictDBIndexer.Search", keyword)
	defer logend("MedictDBIndexer.Search", startTime, resErr)

	// Single iteration: Prefix now returns key+value, so we avoid the previous
	// N+1 (one Prefix scan then one Get per key) — issue #722 P2.
	kvs, err := m.lvdb.Prefix(strip(keyword))
	if err != nil {
		return nil, err
	}

	list := make([]*model.MdictKeyWordIndex, 0, len(kvs))
	for idx, kv := range kvs {
		vi := new(model.MdictKeyWordIndex)
		if err1 := json.Unmarshal(kv.Value, vi); err1 != nil {
			log.Errorf("unmarshal value failed, %s", err1)
			continue
		}
		vi.ID = idx
		list = append(list, vi)
	}

	// Empty result is not an error (issue #722 P3 error-as-control-flow) —
	// callers check len(result), not the error.
	if len(list) == 0 {
		return nil, nil
	}
	return list, nil
}
