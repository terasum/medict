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

package gomdict

import (
	"errors"
	"fmt"
	"github.com/agatan/bktree"
	"github.com/creasty/go-levenshtein"
	"path/filepath"
	"strings"
)

type Mdict struct {
	bktree *bktree.BKTree
	*MdictBase
}

func New(filename string) (*Mdict, error) {
	dictType := MdictTypeMdx
	if strings.ToLower(filepath.Ext(filename)) == ".mdd" {
		dictType = MdictTypeMdd
	}
	mdict := &Mdict{
		MdictBase: &MdictBase{FilePath: filename, FileType: dictType},
	}
	return mdict, mdict.init()
}

func (mdict *Mdict) init() error {
	err := mdict.ReadDictHeader()
	if err != nil {
		return err
	}

	err = mdict.ReadKeyBlockMeta()
	if err != nil {
		return err
	}

	err = mdict.ReadKeyBlockInfo()
	if err != nil {
		return err
	}

	err = mdict.ReadKeyEntries()
	if err != nil {
		return err
	}

	err = mdict.ReadRecordBlockMeta()
	if err != nil {
		return err
	}

	err = mdict.ReadRecordBlockInfo()
	if err != nil {
		return err
	}

	err = mdict.BuildBKTree()
	if err != nil {
		return err
	}

	return nil
}

func (mdict *Mdict) Lookup(word string) ([]byte, error) {

	for _, keyBlockEntry := range mdict.KeyBlockData.KeyEntries {
		if strings.TrimSpace(keyBlockEntry.KeyWord) == strings.TrimSpace(word) {
			return mdict.LocateRecordDefinition(keyBlockEntry)
		}
	}
	return nil, fmt.Errorf("word:(%s) not found", word)
}

func (mdict *Mdict) Locate(entry *MDictKeyBlockEntry) ([]byte, error) {
	if entry == nil {
		return nil, errors.New("invalid mdict key block entry")
	}
	return mdict.LocateRecordDefinition(entry)
}

func (mdict *Mdict) Search(word string) ([]*MDictKeyBlockEntry, error) {
	if mdict.bktree == nil {
		return nil, errors.New("bktree hasn't build yet")
	}
	result, err := mdict.SimSearch(word, 1)
	return result, err
}

// Distance calculates hamming distance.
func (x *MDictKeyBlockEntry) Distance(e bktree.Entry) int {
	a := x.KeyWord
	b := e.(*MDictKeyBlockEntry).KeyWord

	return levenshtein.Distance(a, b)
}

func (mdict *Mdict) BuildBKTree() error {
	mdict.bktree = &bktree.BKTree{}
	for _, e := range mdict.KeyBlockData.KeyEntries {
		mdict.bktree.Add(e)
	}
	return nil
}

func (mdict *Mdict) SimSearch(word string, tolerance int) ([]*MDictKeyBlockEntry, error) {
	entry := &MDictKeyBlockEntry{KeyWord: word}
	results := mdict.bktree.Search(entry, tolerance)
	resultWords := make([]*MDictKeyBlockEntry, 0)
	for _, result := range results {
		resultWords = append(resultWords, result.Entry.(*MDictKeyBlockEntry))
	}
	return resultWords, nil
}
