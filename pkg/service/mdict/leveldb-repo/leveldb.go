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

// Prefix returns the keys with the given prefix, in leveldb's sorted order.
func (lvdb *LvDB) Prefix(prefix string) ([]string, error) {
	ite := lvdb.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer ite.Release()

	result := make([]string, 0)
	for ite.Next() {
		result = append(result, string(ite.Key()))
	}
	if err := ite.Error(); err != nil {
		return nil, err
	}
	return result, nil
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
