//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package medengine

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
	"github.com/terasum/medict/internal/utils"
)

var log = logging.MustGetLogger("default")

// CreateMeIndex creates Medict index file
// Medict words index file format:
// filename: meidx
// internal format: sqlite unzipped
// sqlite table name: meidx_keyword_index
// table columns:
// --------------------------
// | idx_no | key_word | key_block_index | record_start_offset | record_end_offset | compressed_size | decompressed_size | dict_type |
func CreateMeIndex(idxFilePath string) error {
	if utils.FileExists(idxFilePath) {
		return nil
	}

	db, err := sql.Open("sqlite3", idxFilePath)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
	DROP TABLE if EXISTS meidx_keyword_index;
	CREATE TABLE meidx_keyword_index (
	   idx_no integer primary key autoincrement not null, 
	   key_word varchar(512) unique,
	   key_block_index long ,
	   record_start_offset long ,
	   record_end_offset long);
-- 	CREATE INDEX index_meidx_keyword_index_keyword ON meidx_keyword_index(key_word);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}
