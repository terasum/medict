package mdict_idxer

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/terasum/medict/internal/libs/bktree"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
	"time"

	"github.com/silenceper/pool"
)

type MdictIdxer struct {
	indexFilePath string
	bktree        *bktree.BKTree
	connPool      pool.Pool
}

func NewIdxer(indexFilePath string) (*MdictIdxer, error) {
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

	engine := &MdictIdxer{
		indexFilePath: indexFilePath,
		connPool:      connPool,
	}

	err = engine.migrate()
	return engine, err
}

// migrate creates Medict index file
// Medict words index file format:
// filename: meidx
// internal format: sqlite unzipped
// sqlite table name: meidx_keyword_index
// table columns:
// --------------------------
// | idx_no | keyword | key_block_index | record_start_offset | record_end_offset | compressed_size | decompressed_size | dict_type |
func (idxer *MdictIdxer) migrate() error {
	if utils.FileExists(idxer.indexFilePath) {
		return nil
	}

	db, err := sql.Open("sqlite3", idxer.indexFilePath)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
	DROP TABLE if EXISTS meidx_keyword_index;
	CREATE TABLE meidx_keyword_index (
	   idx_no 							 integer primary key autoincrement not null, 
	   keyword 						     varchar(512) unique,
	   record_locate_start_offset   	 long,
	   record_locate_end_offset     	 long,
	   record_block_data_start_offset 	 long,
	   record_block_data_compress_size   long ,
	   record_block_data_decompress_size long ,
       keyword_data_start_offset    	 long ,
       keyword_data_end_offset      	 long);
-- 	CREATE INDEX index_meidx_keyword_index_keyword ON meidx_keyword_index(keyword);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
	DROP TABLE if EXISTS meidx_dict_meta;
	CREATE TABLE meidx_dict_meta (key varchar(256), value varchar(512));`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func (idxer *MdictIdxer) acquire() (*sql.DB, error) {
	db, err := idxer.connPool.Get()
	return db.(*sql.DB), err
}

func (idxer *MdictIdxer) release(db *sql.DB) {
	err := idxer.connPool.Put(db)
	if err != nil {
		log.Errorf(err.Error())
	}
}

func (idxer *MdictIdxer) close() {
	idxer.connPool.Release()
}

func (idxer *MdictIdxer) Lookup(keyword string) (*model.MdictKeyWordIndex, error) {
	db, err := idxer.acquire()
	if err != nil {
		return nil, err
	}
	defer idxer.release(db)

	sqlfmt := `SELECT 
    idx_no,
       keyword,
	   record_locate_start_offset,
	   record_locate_end_offset, 
	   record_block_data_start_offset,
	   record_block_data_compress_size,
	   record_block_data_decompress_size,
       keyword_data_start_offset,
       keyword_data_end_offset
	 FROM meidx_keyword_index WHERE keyword = ?`
	result := db.QueryRow(sqlfmt, keyword)
	temp := new(model.MdictKeyWordIndex)
	err1 := result.Scan(
		&(temp.ID),
		&(temp.KeyWord),
		&(temp.RecordLocateStartOffset),
		&(temp.RecordLocateEndOffset),
		&(temp.RecordBlockDataStartOffset),
		&(temp.RecordBlockDataCompressSize),
		&(temp.RecordBlockDataDeCompressSize),
		&(temp.KeyWordDataStartOffset),
		&(temp.KeyWordDataEndOffset))
	if err1 != nil {
		log.Errorf("sql query error %s", err1.Error())
		return nil, err1
	}
	return temp, nil

}

func (idxer *MdictIdxer) SetMeta(key, value string) error {
	db, err := idxer.acquire()
	if err != nil {
		return err
	}
	defer idxer.release(db)
	sqlfmt := `
	BEGIN;
	INSERT OR IGNORE INTO meidx_dict_meta (key, value) VALUES(?, ?);
	UPDATE meidx_dict_meta SET value=? WHERE key = ?;
	COMMIT;
	`
	_, err = db.Exec(sqlfmt, key, value, value, key)
	return err
}

func (idxer *MdictIdxer) GetMeta(key string) (value string, err error) {
	db, err := idxer.acquire()
	if err != nil {
		return "", err
	}
	defer idxer.release(db)

	sqlfmt := `SELECT value FROM meidx_dict_meta WHERE key = ?`
	statement, err := db.Prepare(sqlfmt)
	if err != nil {
		return "", err
	}

	err = statement.QueryRow(key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil

}

func (idxer *MdictIdxer) AddRecord(record *model.MdictKeyWordIndex) error {
	db, err := idxer.acquire()
	if err != nil {
		return err
	}
	defer idxer.release(db)

	sqlfmt := `INSERT OR IGNORE INTO meidx_keyword_index (
       keyword,
	   record_locate_start_offset,
	   record_locate_end_offset, 
	   record_block_data_start_offset,
	   record_block_data_compress_size,
	   record_block_data_decompress_size,
       keyword_data_start_offset,
       keyword_data_end_offset)
	   VALUES(?, ?, ?, ?, ?, ?, ?, ?)`

	statement, err := db.Prepare(sqlfmt)
	if err != nil {
		return err
	}
	_, err = statement.Exec(
		record.KeyWord,
		record.RecordLocateStartOffset,
		record.RecordLocateEndOffset,
		record.RecordBlockDataStartOffset,
		record.RecordBlockDataCompressSize,
		record.RecordBlockDataDeCompressSize,
		record.KeyWordDataStartOffset,
		record.KeyWordDataEndOffset)
	return err
}

func (idxer *MdictIdxer) Search(keyword string) ([]*model.MdictKeyWordIndex, error) {
	db, err := idxer.acquire()
	if err != nil {
		return nil, err
	}
	defer idxer.release(db)

	sqlfmt := `SELECT 
    	idx_no,
       keyword,
	   record_locate_start_offset,
	   record_locate_end_offset, 
	   record_block_data_start_offset,
	   record_block_data_compress_size,
	   record_block_data_decompress_size,
       keyword_data_start_offset,
       keyword_data_end_offset
	 FROM meidx_keyword_index WHERE keyword LIKE ?`
	statement, err := db.Prepare(sqlfmt)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Query(keyword + "%")
	if err != nil {
		return nil, err
	}
	results := make([]*model.MdictKeyWordIndex, 0)
	for result.Next() {
		temp := new(model.MdictKeyWordIndex)
		err1 := result.Scan(
			&(temp.ID),
			&(temp.KeyWord),
			&(temp.RecordLocateStartOffset),
			&(temp.RecordLocateEndOffset),
			&(temp.RecordBlockDataStartOffset),
			&(temp.RecordBlockDataCompressSize),
			&(temp.RecordBlockDataDeCompressSize),
			&(temp.KeyWordDataStartOffset),
			&(temp.KeyWordDataEndOffset))
		if err1 != nil {
			log.Errorf("sql query error %s", err1.Error())
		}
		results = append(results, temp)
	}
	return results, nil
}
