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
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestWalkDir(t *testing.T) {
	result, err := WalkDir("./testdata/dicts")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	expectPath, _ := filepath.Abs("testdata/dicts/oale3/test.mdx")
	assert.Equal(t, expectPath, result[0].MdxAbsPath)
}
