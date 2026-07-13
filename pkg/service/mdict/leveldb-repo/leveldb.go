package leveldb_repo

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// LvDB wraps a single embedded leveldb handle.
//
// goleveldb's *leveldb.DB is itself goroutine-safe, so one shared handle serves
// all callers — no connection pool needed. The previous silenceper/pool design
// could hand out a nil connection when its idle reaper tried to re-open the
// same file and tripped goleveldb's flock, then panicked on the unchecked
// db.(*leveldb.DB) assertion in acquire (issue #722, live panic during
// BuildIndex). Holding one handle removes that failure class entirely.
type LvDB struct {
	dbFileDirPath string
	db            *leveldb.DB
}

func NewLvDB(fpath string) (*LvDB, error) {
	db, err := leveldb.OpenFile(fpath, nil)
	if err != nil {
		return nil, err
	}
	return &LvDB{
		dbFileDirPath: fpath,
		db:            db,
	}, nil
}

// KV is a single leveldb key-value pair returned by Prefix.
type KV struct {
	Key   []byte
	Value []byte
}

// Prefix returns all key-value pairs whose key starts with prefix, in leveldb's
// sorted order. Keys and values are copied (the iterator's buffers are
// invalidated on Next), so the returned slices are safe to retain.
func (lvdb *LvDB) Prefix(prefix string) ([]KV, error) {
	ite := lvdb.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer ite.Release()

	result := make([]KV, 0)
	for ite.Next() {
		result = append(result, KV{
			Key:   append([]byte(nil), ite.Key()...),
			Value: append([]byte(nil), ite.Value()...),
		})
	}
	if err := ite.Error(); err != nil {
		return nil, err
	}
	return result, nil
}

// Write applies a batch atomically. Used for bulk index builds.
func (lvdb *LvDB) Write(batch *leveldb.Batch) error {
	return lvdb.db.Write(batch, nil)
}

// Put a key-value pair into leveldb.
func (lvdb *LvDB) Put(key string, value []byte) error {
	return lvdb.db.Put([]byte(key), value, nil)
}

// Get a key-value pair from leveldb.
func (lvdb *LvDB) Get(key string) ([]byte, error) {
	return lvdb.db.Get([]byte(key), nil)
}

// Close releases the leveldb handle. Operations after Close return leveldb's
// ErrClosed rather than panicking.
func (lvdb *LvDB) Close() error {
	return lvdb.db.Close()
}
