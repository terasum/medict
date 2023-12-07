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

package go_mdict

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("default")

type Mdict struct {
	*MdictBase
	rangeTreeRoot *RecordBlockRangeTreeNode
}

func New(filename string) (*Mdict, error) {
	dictType := MdictTypeMdx
	if strings.ToLower(filepath.Ext(filename)) == ".mdd" {
		dictType = MdictTypeMdd
	}

	mdict := &Mdict{
		MdictBase: &MdictBase{
			filePath:      filename,
			fileType:      dictType,
			rangeTreeRoot: new(RecordBlockRangeTreeNode),
		},
	}
	return mdict, mdict.init()
}

func (mdict *Mdict) init() error {
	// 读取词典头
	err := mdict.readDictHeader()
	if err != nil {
		return err
	}

	// 读取 key block 元信息
	err = mdict.readKeyBlockMeta()
	if err != nil {
		return err
	}

	return nil
}

// BuildIndex 构建索引
func (mdict *Mdict) BuildIndex() error {
	err := mdict.readKeyBlockInfo()
	if err != nil {
		return err
	}

	err = mdict.readKeyEntries()
	if err != nil {
		return err
	}

	err = mdict.readRecordBlockMeta()
	if err != nil {
		return err
	}

	err = mdict.readRecordBlockInfo()
	if err != nil {
		return err
	}

	mdict.buildRecordRangeTree()

	return nil
}

func (mdict *Mdict) Name() string {
	_, rawpath := filepath.Split(mdict.filePath)
	rawpath = strings.TrimRight(rawpath, ".mdx")
	if len(rawpath) > 0 {
		return rawpath
	}
	return rawpath
}

func (mdict *Mdict) Title() string {
	return mdict.meta.title

}

func (mdict *Mdict) Description() string {
	return mdict.meta.description
}
func (mdict *Mdict) GeneratedByEngineVersion() string {
	return mdict.meta.generatedByEngineVersion
}
func (mdict *Mdict) CreationDate() string {
	return mdict.meta.creationDate
}
func (mdict *Mdict) Version() string {
	return fmt.Sprintf("%f", mdict.meta.version)
}

func (mdict *Mdict) IsMDD() bool {
	return mdict.fileType == MdictTypeMdd
}

func (mdict *Mdict) IsRecordEncrypted() bool {
	return mdict.meta.encryptType == EncryptRecordEnc
}

func (mdict *Mdict) IsUTF16() bool {
	return mdict.meta.encoding == EncodingUtf16
}

func (mdict *Mdict) Lookup(word string) ([]byte, error) {
	word = strings.TrimSpace(word)
	for id, keyBlockEntry := range mdict.keyBlockData.keyEntries {
		if keyBlockEntry.KeyWord == word {
			log.Infof("mdict.Lookup hit entries[%d/%d] key:(%s), entry-key:(%s), equals(%v)", id, len(mdict.keyBlockData.keyEntries), word, keyBlockEntry.KeyWord, keyBlockEntry.KeyWord == word)
			return mdict.LocateByKeywordEntry(keyBlockEntry)
		}
	}
	return nil, fmt.Errorf("word:(%s) not found", word)
}

func (mdict *Mdict) LocateByKeywordEntry(entry *MDictKeywordEntry) ([]byte, error) {
	if entry == nil {
		return nil, errors.New("invalid mdict keyword entry")
	}
	return mdict.MdictBase.locateByKeywordEntry(entry)
}

func (mdict *Mdict) LocateByKeywordIndex(index *MDictKeywordIndex) ([]byte, error) {
	if index == nil {
		return nil, errors.New("invalid mdict keyword index")
	}
	return mdict.MdictBase.locateByKeywordIndex(index)

}

func (mdict *Mdict) GetKeyWordEntries() ([]*MDictKeywordEntry, error) {
	return mdict.getKeyWordEntries()
}

func (mdict *Mdict) GetKeyWordEntriesSize() int64 {
	return mdict.keyBlockData.keyEntriesSize
}

func (mdict *Mdict) KeywordEntryToIndex(item *MDictKeywordEntry) (*MDictKeywordIndex, error) {
	return mdict.MdictBase.keywordEntryToIndex(item)
}
