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
	"os"
	"path/filepath"

	"github.com/terasum/medict/internal/utils"

	"github.com/terasum/medict/pkg/model"
)

var log = logging.MustGetLogger("service.filewalker")

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

		// filepath.WalkDir does not follow symlinks, and a symlink's DirEntry
		// reports ModeSymlink (not a directory) even when it points to one.
		// Resolve symlinks so a dictionary linked into the dicts dir is scanned.
		scanPath := path
		if !d.IsDir() {
			if d.Type()&fs.ModeSymlink == 0 {
				// a real, non-directory file at this level: skip
				return nil
			}
			resolved, rerr := filepath.EvalSymlinks(path)
			if rerr != nil {
				log.Errorf("resolve symlink failed, path:[%s], %s", path, rerr.Error())
				return nil
			}
			fi, serr := os.Stat(resolved)
			if serr != nil || !fi.IsDir() {
				// broken link or link to a non-directory: skip
				return nil
			}
			scanPath = resolved
		}
		// 遍历第二层
		item, err := innerWalkLevel2(dirpath, scanPath)
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

	// 只看 DIRECT 子文件(不递归进子目录):一个目录直接含词典文件才算一个词典;
	// 分类/分组目录(词典在更深层)由 WalkDir 在更深层各自评估,避免嵌套时重复发现
	// 同一个词典、或把分类目录误当词典。(#257)
	entries, err1 := os.ReadDir(level2Path)
	if err1 != nil {
		return item, err1
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		path := filepath.Join(level2Path, name)

		// predefined file type
		if name == "cover.jpg" {
			item.CoverImgPath = utils.FileAbs(path)
			item.CoverImgType = model.ImgTypeJPG
		} else if name == "cover.png" {
			item.CoverImgPath = utils.FileAbs(path)
			item.CoverImgType = model.ImgTypePNG
		} else if name == "mdict.toml" {
			item.DictType = model.DictTypeMdict
			item.ConfigPath = utils.FileAbs(path)
		} else if name == "stardict.toml" {
			item.DictType = model.DictTypeStarDict
			item.ConfigPath = utils.FileAbs(path)
		} else if name == "dict.license" {
			item.LicensePath = utils.FileAbs(path)
		}

		// if mdx
		if filepath.Ext(name) == ".mdx" {
			item.MdictMdxAbsPath, _ = filepath.Abs(path)
			item.MdictMdxFileName = utils.FileNameWithoutExt(path)
			baseDir := utils.FileBaseDir(path)
			log.Debugf("mdx path: %s, basedir: %s", path, baseDir)
			pngPath := filepath.Join(baseDir, item.MdictMdxFileName+"."+"png")
			log.Debugf("cover png candidate: %s", pngPath)
			if utils.FileExists(pngPath) {
				item.CoverImgPath = utils.FileAbs(pngPath)
				item.CoverImgType = model.ImgTypePNG
			}

			jpgPath := filepath.Join(baseDir, item.MdictMdxFileName+"."+"jpg")
			log.Debugf("cover jpg candidate: %s", jpgPath)
			if utils.FileExists(jpgPath) {
				item.CoverImgPath = utils.FileAbs(jpgPath)
				item.CoverImgType = model.ImgTypeJPG
			}
			item.DictType = model.DictTypeMdict
			item.IsValid = true
		} else if filepath.Ext(name) == ".mdd" {
			// MDD append
			mddAbs, _ := filepath.Abs(path)
			item.MdictMddAbsPath = append(item.MdictMddAbsPath, mddAbs)
		}

		// if stardict
		if filepath.Ext(name) == ".dz" {
			item.StarDictDzAbsPath, _ = filepath.Abs(path)
			item.DictType = model.DictTypeStarDict
			item.IsValid = true
		} else if filepath.Ext(name) == ".dict" {
			item.StarDictAbsPath, _ = filepath.Abs(path)
			item.DictType = model.DictTypeStarDict
			item.IsValid = true
		} else if filepath.Ext(name) == ".ifo" {
			item.StarDictIfoAbsPath, _ = filepath.Abs(path)
		} else if filepath.Ext(name) == ".idx" {
			item.StarDictIdxAbsPath, _ = filepath.Abs(path)
		}

		// ecdict:离线英汉 preset(目录里有 ecdict.db 即识别)
		if name == "ecdict.db" {
			item.DictType = model.DictTypeECDICT
			item.IsValid = true
		}
	}
	return item, nil
}
