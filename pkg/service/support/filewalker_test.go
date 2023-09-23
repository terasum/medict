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

package support

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/terasum/medict/pkg/model"
	"path/filepath"
	"testing"
)

func TestWalkDir(t *testing.T) {
	result, err := WalkDir("./testdata/dicts")
	t.Logf("result: %+v", result)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(result))

	expectPath, _ := filepath.Abs("testdata/dicts/ccedit/ccedit.idx")
	assert.Equal(t, expectPath, result[0].StarDictIdxAbsPath)
	assert.Equal(t, model.DictTypeStarDict, result[0].DictType)

	expectPath2, _ := filepath.Abs("testdata/dicts/mdict-type/test.mdx")
	assert.Equal(t, expectPath2, result[1].MdictMdxAbsPath)
	assert.Equal(t, model.ImgTypeJPG, result[1].CoverImgType)
	license, _ := filepath.Abs("testdata/dicts/mdict-type/dict.license")
	assert.Equal(t, license, result[1].LicensePath)

	expectPath3, _ := filepath.Abs("testdata/dicts/oale3/test.mdx")
	assert.Equal(t, expectPath3, result[2].MdictMdxAbsPath)
	assert.Equal(t, model.ImgTypePNG, result[2].CoverImgType)
	license2, _ := filepath.Abs("testdata/dicts/oale3/dict.license")
	assert.Equal(t, license2, result[2].LicensePath)

	data, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(string(data))
}
