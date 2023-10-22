package medengine

import (
	"database/sql"
	"github.com/silenceper/pool"
	"github.com/terasum/medict/pkg/model"
	"time"
)

type IndexEngine struct {
	indexFilePath string
	connPool      pool.Pool
}

type IndexRecord struct {
	keyWord           string
	keyBlockIndex     int64
	recordStartOffset int64
	recordEndOffset   int64
}

func NewIndexRecord(
	KeyWord string,
	RecordStartOffset int64,
	RecordEndOffset int64,
	KeyBlockIdx int64,
) *IndexRecord {
	return &IndexRecord{
		keyWord:           KeyWord,
		keyBlockIndex:     KeyBlockIdx,
		recordStartOffset: RecordStartOffset,
		recordEndOffset:   RecordEndOffset,
	}
}

func (idr *IndexRecord) ToKeyBlockEntry() *model.KeyBlockEntry {
	return &model.KeyBlockEntry{
		ID:                0,
		RecordStartOffset: idr.recordStartOffset,
		RecordEndOffset:   idr.recordEndOffset,
		KeyWord:           idr.keyWord,
		KeyBlockIdx:       idr.keyBlockIndex,
	}
}

func NewEngine(indexFilePath string) (*IndexEngine, error) {
	err := CreateMeIndex(indexFilePath)
	if err != nil {
		return nil, err
	}

	connPool, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 1,
		MaxCap:     3,
		MaxIdle:    3,
		Factory: func() (interface{}, error) {
			db, err1 := sql.Open("sqlite3", indexFilePath)
			if err1 != nil {
				return nil, err1
			}
			return db, nil
		},
		Close: func(db interface{}) error {
			return db.(*sql.DB).Close()
		},
		Ping: func(db interface{}) error {
			return nil
		},
		IdleTimeout: 60 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &IndexEngine{
		indexFilePath: indexFilePath,
		connPool:      connPool,
	}, nil
}

func (engine *IndexEngine) AddRecord(record *IndexRecord) error {
	db, err := engine.Acquire()
	if err != nil {
		return err
	}
	defer engine.Release(db)

	sqlfmt := `INSERT INTO meidx_keyword_index (
	   key_word,
	   key_block_index,
	   record_start_offset,
	   record_end_offset)
	   VALUES(?, ?, ?, ?)`
	statement, err := db.Prepare(sqlfmt)
	if err != nil {
		return err
	}
	_, err = statement.Exec(record.keyWord, record.keyBlockIndex, record.recordStartOffset, record.recordEndOffset)
	return err
}

func (engine *IndexEngine) Search(keyword string) ([]*IndexRecord, error) {
	db, err := engine.Acquire()
	if err != nil {
		return nil, err
	}
	defer engine.Release(db)

	sqlfmt := `SELECT key_word, key_block_index, record_start_offset, record_end_offset FROM meidx_keyword_index WHERE key_word LIKE ?`
	statement, err := db.Prepare(sqlfmt)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Query(keyword + "%")
	if err != nil {
		return nil, err
	}
	results := make([]*IndexRecord, 0)
	for result.Next() {
		temp := new(IndexRecord)
		err1 := result.Scan(&(temp.keyWord), &(temp.keyBlockIndex), &(temp.recordStartOffset), &(temp.recordEndOffset))
		if err1 != nil {
			log.Errorf("sql query error %s", err1.Error())
		}
		results = append(results, temp)
	}
	return results, nil
}

func (engine *IndexEngine) Acquire() (*sql.DB, error) {
	db, err := engine.connPool.Get()
	return db.(*sql.DB), err
}

func (engine *IndexEngine) Release(db *sql.DB) {
	err := engine.connPool.Put(db)
	if err != nil {
		log.Errorf(err.Error())
	}
}

func (engine *IndexEngine) Close() {
	engine.connPool.Release()
}
