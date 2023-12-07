package leveldb_repo

import (
	"github.com/silenceper/pool"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"time"
)

type LvDB struct {
	dbFileDirPath string
	connPool      pool.Pool
}

func NewLvDB(fpath string) (*LvDB, error) {
	connPool, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 1,
		MaxCap:     3,
		MaxIdle:    3,
		Factory: func() (interface{}, error) {
			db, err1 := leveldb.OpenFile(fpath, nil)
			if err1 != nil {
				return nil, err1
			}
			return db, nil
		},
		Close: func(db interface{}) error {
			return db.(*leveldb.DB).Close()
		},
		Ping: func(db interface{}) error {
			return nil
		},
		IdleTimeout: 60 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	idxer := &LvDB{
		dbFileDirPath: fpath,
		connPool:      connPool,
	}
	return idxer, nil
}

func (lvdb *LvDB) acquire() (*leveldb.DB, error) {
	db, err := lvdb.connPool.Get()
	return db.(*leveldb.DB), err
}

func (lvdb *LvDB) release(db *leveldb.DB) {
	err := lvdb.connPool.Put(db)
	if err != nil {
		log.Errorf(err.Error())
	}
}

// Prefix retrieve keys list from leveldb
func (lvdb *LvDB) Prefix(prefix string) ([]string, error) {
	db, err := lvdb.acquire()
	if err != nil {
		return nil, err
	}
	defer lvdb.release(db)
	ite := db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

	result := make([]string, 0)
	for ite.Next() {
		result = append(result, string(ite.Key()))
	}
	return result, nil
}

// Put a key-value pair into leveldb
func (lvdb *LvDB) Put(key string, value []byte) error {
	db, err := lvdb.acquire()
	if err != nil {
		return err
	}
	defer lvdb.release(db)

	return db.Put([]byte(key), value, nil)
}

// Get a key-value pair from leveldb
func (lvdb *LvDB) Get(key string) ([]byte, error) {
	db, err := lvdb.acquire()
	if err != nil {
		return nil, err
	}
	defer lvdb.release(db)
	value, err := db.Get([]byte(key), nil)
	return value, err
}

//func sortList(keyword string, list []*model.MdictKeyWordIndex) (error) {
//	wrapper := &entryWrapper{
//		list: make([]*entryWrapperItem, 0),
//	}
//
//	distance := levenshtein.Distance(keyword, valueIndex.KeyWord)
//	wrapper.list = append(wrapper.list, &entryWrapperItem{
//		entry:    valueIndex,
//		distance: distance,
//	})
//	sort.Sort(wrapper)
//	return nil
//}
