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
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/terasum/medict/internal/utils"

	"github.com/terasum/medict/pkg/model"
)

func WalkDir(dirpath string) ([]*model.DirItem, error) {
	list := make([]*model.DirItem, 0)
	err := filepath.WalkDir(dirpath, func(path string, d fs.DirEntry, err error) error {
		if dirpath == "" {
			return fmt.Errorf("walkdir failed, path is empty, path: [%s] %s", path, err.Error())
		}
		if err != nil {
			return fmt.Errorf("walkdir failed, path: [%s] %s", path, err.Error())
		}
		// skip self
		if dirpath == path {
			return nil
		}
		// skip non-dir
		if !d.IsDir() {
			return nil
		}
		item, err := innerWalker(dirpath, path, err)
		if err != nil {
			return fmt.Errorf("inner walker failed , path:[%s], %s", path, err.Error())
		}
		if item.MdictMdxAbsPath != "" {
			list = append(list, item)
		}
		return nil
	})
	return list, err
}

func innerWalker(rootpath, subpath string, err error) (*model.DirItem, error) {
	if rootpath == "" {
		return nil, fmt.Errorf("inner walkdir failed, path is empty, path: [%s] %s", rootpath, err.Error())
	}
	if err != nil {
		return nil, fmt.Errorf("inner walk entry failed, path %s, err %s", rootpath, err.Error())
	}
	item := &model.DirItem{
		BaseDir:         rootpath,
		CurrentDir:      subpath,
		MdictMdxAbsPath: "",
		MdictMddAbsPath: make([]string, 0),
	}

	err = filepath.Walk(subpath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("inner walk dir failed, path %s, %s", path, err)
		}
		if info.IsDir() {
			return nil
		}

		// if mdx
		if filepath.Ext(info.Name()) == ".mdx" {
			item.MdictMdxAbsPath, _ = filepath.Abs(path)
		} else if filepath.Ext(info.Name()) == ".mdd" {
			// pop stack
			mddabs, _ := filepath.Abs(path)
			item.MdictMddAbsPath = append(item.MdictMddAbsPath, mddabs)
		} else if info.Name() == "cover.jpg" || info.Name() == "cover.png" {
			item.CoverImgPath = utils.FileAbs(path)
		} else {
			// TODO copy css/js files
		}
		return nil
	})
	return item, err
}
