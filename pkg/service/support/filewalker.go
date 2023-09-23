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
	"github.com/op/go-logging"
	"io/fs"
	"path/filepath"

	"github.com/terasum/medict/internal/utils"

	"github.com/terasum/medict/pkg/model"
)

var log = logging.MustGetLogger("default")

// WalkDir 遍历第一层的所有文件夹，忽略文件
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
		if dirpath == path || path == "." || path == ".." {
			return nil
		}

		// skip non-dir
		if !d.IsDir() {
			return nil
		}
		// 遍历第二层
		item, err := innerWalkLevel2(dirpath, path)
		if err != nil {
			// 二层遍历失败，继续遍历下一个
			log.Errorf("inner walker failed , path:[%s], %s", path, err.Error())
			return nil
		}

		if item != nil && item.IsValid {
			list = append(list, item)
		}

		return nil
	})
	return list, err
}

func innerWalkLevel2(level1Path, level2Path string) (*model.DirItem, error) {
	if level1Path == "" {
		// this should not take event
		return nil, fmt.Errorf("level2dir walk failed, level1path is empty, err: %s", level1Path)
	}

	item := &model.DirItem{
		BaseDir:         utils.FileAbs(level1Path),
		CurrentDir:      utils.FileAbs(level2Path),
		MdictMddAbsPath: make([]string, 0),
		IsValid:         false,
	}

	err1 := filepath.Walk(level2Path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("level2 inner walk dir failed, path %s, %s", path, err)
		}

		// skip directory
		if info.IsDir() {
			return nil
		}

		// verify multiple directory types, such as
		// order0: medict.type file
		// order1: medict type
		// order2: stardict type

		// predefined file type
		if info.Name() == "cover.jpg" {
			item.CoverImgPath = utils.FileAbs(path)
			item.CoverImgType = model.ImgTypeJPG
		} else if info.Name() == "cover.png" {
			item.CoverImgPath = utils.FileAbs(path)
			item.CoverImgType = model.ImgTypePNG
		} else if info.Name() == "mdict.toml" {
			item.DictType = model.DictTypeMdict
			item.ConfigPath = utils.FileAbs(path)
		} else if info.Name() == "stardict.toml" {
			item.DictType = model.DictTypeStarDict
			item.ConfigPath = utils.FileAbs(path)
		} else if info.Name() == "dict.license" {
			item.LicensePath = utils.FileAbs(path)
		}

		// if mdx
		if filepath.Ext(info.Name()) == ".mdx" {
			item.MdictMdxAbsPath, _ = filepath.Abs(path)
			item.MdictMdxFileName = utils.FileNameWithoutExt(path)
			baseDir := utils.FileBaseDir(path)
			fmt.Printf("path is %s, basedir is %s\n", path, baseDir)
			pngPath := filepath.Join(baseDir, item.MdictMdxFileName+"."+"png")
			fmt.Printf("pngpath: %s\n", pngPath)
			if utils.FileExists(pngPath) {
				item.CoverImgPath = utils.FileAbs(pngPath)
				item.CoverImgType = model.ImgTypePNG
			}

			jpgPath := filepath.Join(baseDir, item.MdictMdxFileName+"."+"jpg")
			fmt.Printf("jpgpath: %s\n", pngPath)
			if utils.FileExists(jpgPath) {
				item.CoverImgPath = utils.FileAbs(jpgPath)
				item.CoverImgType = model.ImgTypeJPG
			}
			item.DictType = model.DictTypeMdict
			item.IsValid = true
		} else if filepath.Ext(info.Name()) == ".mdd" {
			// MDD append
			mddAbs, _ := filepath.Abs(path)
			item.MdictMddAbsPath = append(item.MdictMddAbsPath, mddAbs)
		}

		// if stardict
		if filepath.Ext(info.Name()) == ".dz" {
			item.StarDictDzAbsPath, _ = filepath.Abs(path)
			item.DictType = model.DictTypeStarDict
			item.IsValid = true
		} else if filepath.Ext(info.Name()) == ".dict" {
			item.StarDictAbsPath, _ = filepath.Abs(path)
			item.DictType = model.DictTypeStarDict
			item.IsValid = true
		} else if filepath.Ext(info.Name()) == ".ifo" {
			item.StarDictIfoAbsPath, _ = filepath.Abs(path)
		} else if filepath.Ext(info.Name()) == ".idx" {
			item.StarDictIdxAbsPath, _ = filepath.Abs(path)
		}

		return nil
	})
	return item, err1
}
