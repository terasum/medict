package mdict_idxer

import "github.com/terasum/medict/pkg/model"

type Indexer interface {
	Lookup(keyword string) (*model.MdictKeyWordIndex, error)
	SetMeta(key, value string) error
	GetMeta(key string) (value string, err error)
	AddRecord(record *model.MdictKeyWordIndex) error
	// AddRecords bulk-writes records (used during index builds).
	AddRecords(records []*model.MdictKeyWordIndex) error
	Search(keyword string) ([]*model.MdictKeyWordIndex, error)
	// AllRecords returns every indexed keyword record, used to build in-memory
	// structures (e.g. the fuzzy BK-tree) without re-parsing the source dict.
	AllRecords() ([]*model.MdictKeyWordIndex, error)
	// Close releases the underlying store handle. Use after Close returns an
	// error (not a panic).
	Close() error
}
