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

package utils

import (
	"os"
	"path/filepath"
)

func FetchBaseDirName(fpath string) string {
	return filepath.Base(fpath)
}

func FileAbs(fpath string) string {
	fp, err := filepath.Abs(fpath)
	if err != nil {
		return fpath
	}
	return fp
}

func FileName(fpath string) string {
	paths := filepath.SplitList(fpath)
	if len(paths) <= 1 {
		return fpath
	}
	return paths[len(paths)-1]
}

func FileExists(fpath string) bool {
	if _, err := os.Stat(fpath); err == nil {
		return true
	}
	return false

}
