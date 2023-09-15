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

package model

import (
	"github.com/terasum/medict/internal/gomdict"
	"github.com/terasum/medict/internal/utils"
)

type DictType string

const MDD DictType = "MDD"
const MDX DictType = "MDX"

type DirItem struct {
	BaseDir    string
	CurrentDir string
	MdxAbsPath string
	MddAbsPath []string
}

type PlainDictionaryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type DictionaryItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PathInfo *DirItem
	MDX      *gomdict.Mdict
	MDDS     []*gomdict.Mdict
}

type KeyBlockEntry struct {
	ID                int    `json:"id"`
	RecordStartOffset int64  `json:"record_start_offset"`
	RecordEndOffset   int64  `json:"record_end_offset"`
	KeyWord           string `json:"key_word"`
	KeyBlockIdx       int64  `json:"key_block_idx"`
}

func NewByDirItem(dirItem *DirItem) (*DictionaryItem, error) {
	dictItem := &DictionaryItem{
		ID:       utils.MD5Hash(dirItem.CurrentDir),
		Name:     utils.FetchBaseDirName(dirItem.CurrentDir),
		PathInfo: dirItem,
	}

	mdx, err := gomdict.New(dirItem.MdxAbsPath)
	if err != nil {
		return nil, err
	}

	mdds := make([]*gomdict.Mdict, 0)

	for _, mddpath := range dirItem.MddAbsPath {
		mdd, err1 := gomdict.New(mddpath)
		if err1 != nil {
			return nil, err1
		}
		mdds = append(mdds, mdd)
	}

	dictItem.MDX = mdx
	dictItem.MDDS = mdds

	return dictItem, nil

}
