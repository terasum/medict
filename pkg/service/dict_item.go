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
	"encoding/base64"
	"errors"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
	"os"
)

func NewByDirItem(dirItem *model.DirItem) (*model.DictionaryItem, error) {
	// basic information reading
	dictItem := &model.DictionaryItem{
		PlainDictionaryItem: &model.PlainDictionaryItem{
			ID:       utils.MD5Hash(dirItem.CurrentDir),
			DictDir:  utils.FileAbs(dirItem.CurrentDir),
			Name:     utils.FileName(dirItem.CurrentDir),
			DictType: (string)(dirItem.DictType),
		},
		PathInfo: dirItem,
	}

	if dirItem.CoverImgPath != "" && utils.FileExists(dirItem.CoverImgPath) {
		imgBuffer, err := os.ReadFile(dirItem.CoverImgPath)
		if err != nil {
			log.Errorf("read cover image file failed %s", err.Error())
		} else {
			if dirItem.CoverImgType == model.ImgTypeJPG {
				dictItem.Background = "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(imgBuffer)
			} else {
				dictItem.Background = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBuffer)
			}
		}
	}

	if dirItem.ConfigPath != "" && utils.FileExists(dirItem.ConfigPath) {
		log.Infof("read config file %s\n", dirItem.ConfigPath)
	}

	if dirItem.LicensePath != "" && utils.FileExists(dirItem.LicensePath) {
		log.Infof("read license file %s\n", dirItem.LicensePath)
	}

	if dirItem.DictType == model.DictTypeMdict {
		dict, err := NewMdict(dirItem)
		if err != nil {
			return nil, err
		}
		dictItem.MainDict = dict
		dictItem.Name = dict.Name()
	} else if dirItem.DictType == model.DictTypeStarDict {
		dict, err := NewStardict(dirItem)
		if err != nil {
			return nil, err
		}
		dictItem.Name = dict.Name()
		dictItem.MainDict = dict
	} else {
		return nil, errors.New("not recognized dictionary type")
	}

	dictItem.Description = dictItem.MainDict.Description()

	return dictItem, nil

}
