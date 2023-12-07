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
	"github.com/creasty/go-levenshtein"
	"github.com/terasum/medict/internal/libs/bktree"
	"github.com/terasum/medict/internal/utils"
)

type DirDcitType string

// DirItem
// 词典文件单元，以文件夹为基本单元，一个文件夹代表一个文件
// 词典类型优先级:
// 0: mdict: mdx/mdd
// 1: startdict: dz/ifo/idx
// 特殊文件
// - .dtype 文件: 词典类型文件，可以是 _mdict.dtype 或是 _startdict.dtype 固定名称类型
// - _cover.jpg 文件: 词典封面文件，仅支持 jpg 类型
// - _dict.toml 文件: 当前词典配置文件，控制词典行为，暂未实现
// - _dict.license 文件: 当前词典的授权文件, 通常为 GPL 授权
type DirItem struct {
	BaseDir    string
	CurrentDir string

	IsValid bool

	DictType DictType

	// 特殊文件类型
	CoverImgPath string
	CoverImgType ImgType
	ConfigPath   string
	LicensePath  string

	// medict 文件路径
	MdictMdxFileName string
	MdictMdxAbsPath  string
	MdictMddAbsPath  []string

	// StartDict文件路径
	StarDictDzAbsPath  string
	StarDictAbsPath    string
	StarDictIdxAbsPath string
	StarDictIfoAbsPath string
}

type PlainDictionaryItem struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	DictDir string `json:"base_dir"`

	Background string `json:"background"`
	DictType   string `json:"dict_type"`

	Description *PlainDictionaryInfo `json:"description"`
}

type DictList []*PlainDictionaryItem

func (dlist DictList) Len() int {
	return len(dlist)
}

func (dlist DictList) Swap(i, j int) {
	temp := dlist[i]
	dlist[i] = dlist[j]
	dlist[j] = temp
}

func (dlist DictList) Less(i, j int) bool {
	return dlist[i].ID > dlist[j].ID
}

type PlainDictionaryInfo struct {
	Title                 string `json:"title"`
	Description           string `json:"description"`
	CreateDate            string `json:"createDate"`
	GenerateEngineVersion string `json:"generateEngineVersion"`
}

type DictionaryItem struct {
	*PlainDictionaryItem
	MainDict GeneralDictionary
	PathInfo *DirItem
}

func (dict *DictionaryItem) ToPlain() *PlainDictionaryItem {
	return &PlainDictionaryItem{
		ID:         dict.PlainDictionaryItem.ID,
		Name:       dict.PlainDictionaryItem.Name,
		DictDir:    dict.PlainDictionaryItem.DictDir,
		Background: dict.PlainDictionaryItem.Background,
		DictType:   dict.PlainDictionaryItem.DictType,
		Description: &PlainDictionaryInfo{
			Title:                 dict.PlainDictionaryItem.Description.Title,
			Description:           dict.PlainDictionaryItem.Description.Description,
			CreateDate:            dict.PlainDictionaryItem.Description.CreateDate,
			GenerateEngineVersion: dict.PlainDictionaryItem.Description.GenerateEngineVersion,
		},
	}
}

const IndexTypeMdict = "IndexTypeMdict"
const IndexTypeStardict = "IndexTypeMdict"

type KeyQueryIndex struct {
	IndexType string `json:"index_type"`
	*MdictKeyWordIndex
}

type MdictKeyWordIndex struct {
	ID                            int    `json:"id"`
	KeyWord                       string `json:"keyword"`
	RecordLocateStartOffset       int64  `json:"record_start_offset"`
	RecordLocateEndOffset         int64  `json:"record_end_offset"`
	IsUTF16                       int    `json:"is_utf16"`
	IsRecordEncrypt               int    `json:"is_record_encrypt"`
	IsMDD                         int    `json:"is_mdd"`
	RecordBlockDataStartOffset    int64  `json:"record_block_data_start_offset"`
	RecordBlockDataCompressSize   int64  `json:"record_block_data_compress_size"`
	RecordBlockDataDeCompressSize int64  `json:"record_block_data_decompress_size"`
	KeyWordDataStartOffset        int64  `json:"keyword_data_start_offset"`
	KeyWordDataEndOffset          int64  `json:"keyword_data_end_offset"`
}

// Distance calculates levenshtein distance.
func (x *MdictKeyWordIndex) Distance(e bktree.Entry) int {
	a := x.KeyWord
	b := e.(*MdictKeyWordIndex).KeyWord
	a = utils.StrToUnicode(a)
	b = utils.StrToUnicode(b)

	return levenshtein.Distance(a, b)
}

type MdictMeta struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Filepath         string `json:"filepath"`
	Description      string `json:"description"`
	IsRecordEncoding bool   `json:"is_record_encoding"`
	IsUTF16          bool   `json:"is_utf_16"`
	IsMDD            bool   `json:"is_mdd"`
}
