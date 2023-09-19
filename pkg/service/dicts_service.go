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
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/op/go-logging"

	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service/support"
)

var log = logging.MustGetLogger("default")

var singltonInstanceDictService *DictService

type DictService struct {
	config *config.Config
	dicts  map[string]*model.DictionaryItem
}

func NewDictService(config *config.Config) (*DictService, error) {
	if singltonInstanceDictService != nil {
		return singltonInstanceDictService, nil
	}
	ds := &DictService{
		config: config,
		dicts:  make(map[string]*model.DictionaryItem),
	}
	err := ds.walkDicts()
	if err != nil {
		return nil, err
	}
	singltonInstanceDictService = ds
	return ds, nil
}

func (ds *DictService) FindFromDir(dictId string, key string) ([]byte, error) {
	if dict, ok := ds.dicts[dictId]; ok {
		key = strings.ReplaceAll(key, "\\", string(os.PathSeparator))
		key = strings.TrimLeft(key, string("."))
		key = strings.TrimLeft(key, string(os.PathSeparator))
		fullPath := path.Join(dict.PathInfo.CurrentDir, key)
		if utils.FileExists(fullPath) {
			log.Infof("FindFromDir hitted %s", fullPath)
			return os.ReadFile(fullPath)
		}
		log.Infof("FindFromDir missed %s", fullPath)

	}
	return nil, errors.New("not found from dir")
}

func (ds *DictService) Dicts() []*model.PlainDictionaryItem {
	result := make([]*model.PlainDictionaryItem, len(ds.dicts))
	i := 0
	for _, dict := range ds.dicts {
		result[i] = dict.ToPlain()
		i++
	}
	return result
}

func (ds *DictService) BuildIndex() error {
	var err error
	errIdx := make([]string, 0)
	for idx, dict := range ds.dicts {
		err = dict.MainDict.BuildIndex()
		if err != nil {
			log.Errorf("building dictionary index failed, (%s): %s", dict.ToPlain().Name, err.Error())
			errIdx = append(errIdx, idx)
		}
	}
	if len(errIdx) > 0 {
		for _, idx := range errIdx {
			delete(ds.dicts, idx)
		}
	}
	if len(ds.dicts) == 0 {
		return errors.New("all dictionary index building failed")
	}
	return nil
}

func (ds *DictService) GetDictPlain(id string) (*model.PlainDictionaryItem, bool) {
	dict, ok := ds.dicts[id]
	return dict.ToPlain(), ok
}

func (ds *DictService) Lookup(dictId string, keyword string) ([]byte, error) {
	if dict, ok := ds.dicts[dictId]; !ok {
		return nil, errors.New("dict not found")
	} else {
		data, err := dict.MainDict.Lookup(keyword)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (ds *DictService) LookupResource(dictId string, keyword string) ([]byte, error) {
	if dict, ok := ds.dicts[dictId]; !ok {
		log.Infof("LookResource dict not found [%s]", keyword)
		return nil, fmt.Errorf("dictionary (%s) not found", keyword)
	} else {
		keyword = strings.TrimSpace(keyword)
		data, err := dict.MainDict.LookupResource(keyword)
		if err != nil {
			log.Infof("LookupResource search (%s):[%s] failed, err: %s\n", dict.ToPlain().Name, keyword, err.Error())
			return nil, err
		}
		log.Infof("LookupResource search  (%s)[%s] success\n", dict.ToPlain().Name, keyword)
		return data, nil
	}
}

func (ds *DictService) Locate(dictid string, entry *model.KeyBlockEntry) (string, error) {
	if dict, ok := ds.dicts[dictid]; !ok {
		return "", errors.New("dict not found")
	} else {
		defData, err := dict.MainDict.Locate(&model.KeyIndex{
			IndexType:     model.IndexTypeMdict,
			KeyBlockEntry: entry,
		})
		if err != nil {
			return "", err
		}
		return string(defData), nil
	}
}

func (ds *DictService) Search(dictId string, keyword string) ([]*model.KeyIndex, error) {
	if dict, ok := ds.dicts[dictId]; !ok {
		return nil, errors.New("dict not found")
	} else {
		results, err := dict.MainDict.Search(keyword)
		if err != nil {
			return nil, err
		}
		return results, nil
	}
}

/**
 * walkDicts walk the dict dir and load all dicts
 * @return error
 */
func (ds *DictService) walkDicts() error {
	baseDir, err := utils.ReplaceHome(ds.config.BaseDictDir)
	if err != nil {
		return fmt.Errorf("replace home failed, %s", err.Error())
	}

	items, err := support.WalkDir(baseDir)
	if err != nil {
		return fmt.Errorf("walk dir failed, basedir %s,  %s", baseDir, err.Error())
	}
	for _, dirItem := range items {
		dictItem, err1 := NewByDirItem(dirItem)
		if err1 != nil {
			return fmt.Errorf("new dir item failed, %s", err1.Error())
		}
		ds.dicts[dictItem.ID] = dictItem
	}
	return nil
}
