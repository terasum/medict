package mdict_idxer

import "github.com/terasum/medict/pkg/model"

type Indexer interface {
	Lookup(keyword string) (*model.MdictKeyWordIndex, error)
	SetMeta(key, value string) error
	GetMeta(key string) (value string, err error)
	AddRecord(record *model.MdictKeyWordIndex) error
	Search(keyword string) ([]*model.MdictKeyWordIndex, error)
}
