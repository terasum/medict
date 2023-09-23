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

package service

import (
	"encoding/json"
	"github.com/terasum/medict/pkg/model"
	"testing"
)

func TestNewByDirItem(t *testing.T) {

	const testDirItemData = `[
	{
	"BaseDir": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts",
	"CurrentDir": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/ccedit",
	"IsValid": true,
	"DictType": "StarDict",
	"CoverImgPath": "",
	"CoverImgType": "",
	"ConfigPath": "",
	"LicensePath": "",
	"MdictMdxFileName": "",
	"MdictMdxAbsPath": "",
	"MdictMddAbsPath": [],
	"StarDictDzAbsPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/ccedit/ccedit.dz",
	"StarDictIdxAbsPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/ccedit/ccedit.idx",
	"StarDictIfoAbsPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/ccedit/ccedit.ifo"
	},
	{
	"BaseDir": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts",
	"CurrentDir": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/mdict-type",
	"IsValid": true,
	"DictType": "Mdict",
	"CoverImgPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/mdict-type/test.jpg",
	"CoverImgType": "jpg",
	"ConfigPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/mdict-type/mdict.toml",
	"LicensePath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/mdict-type/dict.license",
	"MdictMdxFileName": "test",
	"MdictMdxAbsPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/mdict-type/test.mdx",
	"MdictMddAbsPath": [
	"/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/mdict-type/test.mdd"
	],
	"StarDictDzAbsPath": "",
	"StarDictIdxAbsPath": "",
	"StarDictIfoAbsPath": ""
	},
	{
	"BaseDir": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts",
	"CurrentDir": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/oale3",
	"IsValid": true,
	"DictType": "Mdict",
	"CoverImgPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/oale3/test.png",
	"CoverImgType": "png",
	"ConfigPath": "",
	"LicensePath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/oale3/dict.license",
	"MdictMdxFileName": "test",
	"MdictMdxAbsPath": "/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/oale3/test.mdx",
	"MdictMddAbsPath": [
	"/Users/chenquan/Workspace/go/src/github.com/terasum/medict/pkg/service/support/testdata/dicts/oale3/test.mdd"
	],
	"StarDictDzAbsPath": "",
	"StarDictIdxAbsPath": "",
	"StarDictIfoAbsPath": ""
	}
]`

	dirItems := make([]*model.DirItem, 0)
	err := json.Unmarshal([]byte(testDirItemData), &dirItems)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dirItems)

	for _, diritem := range dirItems {
		dictItem, err1 := NewByDirItem(diritem)
		if err1 != nil {
			t.Fatal(err1)
		}
		t.Logf("dictitem: %+v", dictItem)
	}
}
