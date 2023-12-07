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
	"github.com/terasum/medict/internal/static/handler"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service/support"
)

var singltonInstanceDictService *DictService

type DictService struct {
	config   *config.Config
	dicts    map[string]*model.DictionaryItem
	dictLock *sync.Mutex
}

func NewDictService(config *config.Config) (*DictService, error) {
	if singltonInstanceDictService != nil {
		return singltonInstanceDictService, nil
	}

	ds := &DictService{
		config:   config,
		dicts:    make(map[string]*model.DictionaryItem),
		dictLock: new(sync.Mutex),
	}

	singltonInstanceDictService = ds

	return ds, nil
}

// InitDicts initialize the dictionaries
func (ds *DictService) InitDicts() error {
	return ds.walkDicts()
}

func (ds *DictService) FindFromDir(dictId string, key string) ([]byte, error) {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	if dict, ok := ds.dicts[dictId]; ok {
		key = strings.ReplaceAll(key, "\\", string(os.PathSeparator))
		key = strings.TrimLeft(key, string("."))
		key = strings.TrimLeft(key, string(os.PathSeparator))
		fullPath := path.Join(dict.PathInfo.CurrentDir, key)
		if utils.FileExists(fullPath) {
			log.Infof("FindFromDir hit %s", fullPath)
			return os.ReadFile(fullPath)
		}
		log.Infof("FindFromDir missed %s", fullPath)

	}
	return nil, errors.New("not found from dir")
}

func (ds *DictService) Dicts() []*model.PlainDictionaryItem {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	result := make([]*model.PlainDictionaryItem, len(ds.dicts))
	i := 0
	for _, dict := range ds.dicts {
		result[i] = dict.ToPlain()
		result[i].Description.Description = handler.WrapDesc(dict.ID, result[i].Name, result[i].Description.Description)
		i++
	}
	list := (model.DictList)(result)
	sort.Sort(list)
	return list
}

func (ds *DictService) GetDictById(id string) *model.DictionaryItem {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	if dict, ok := ds.dicts[id]; ok {
		return dict
	}
	return nil
}

func (ds *DictService) BuildIndexById(dictId string) error {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	if dict, ok := ds.dicts[dictId]; ok {
		err := dict.MainDict.BuildIndex()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *DictService) GetDictPlain(id string) (*model.PlainDictionaryItem, bool) {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	dict, ok := ds.dicts[id]
	return dict.ToPlain(), ok
}

func (ds *DictService) Lookup(dictId string, keyword string) ([]byte, error) {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

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
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

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

func (ds *DictService) Locate(dictid string, idx *model.KeyQueryIndex) (string, error) {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	if dict, ok := ds.dicts[dictid]; !ok {
		return "", errors.New("dict not found")
	} else {
		idxType := model.IndexTypeMdict
		if dict.DictType == (string)(model.DictTypeStarDict) {
			idxType = model.IndexTypeStardict
		}
		defData, err := dict.MainDict.Locate(&model.KeyQueryIndex{
			IndexType:         idxType,
			MdictKeyWordIndex: idx.MdictKeyWordIndex,
		})
		if err != nil {
			return "", err
		}
		return string(defData), nil
	}
}

func (ds *DictService) Search(dictId string, keyword string) ([]*model.KeyQueryIndex, error) {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	log.Infof("search %s %s", dictId, keyword)
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
	baseDir := ds.config.EnsureDictsDir()
	items, err := support.WalkDir(baseDir)
	if err != nil {
		return fmt.Errorf("walk dir failed, basedir %s,  %s", baseDir, err.Error())
	}
	for _, dirItem := range items {
		dictItem, err1 := NewByDirItem(dirItem)
		if err1 != nil {
			return fmt.Errorf("new dir item failed, %s", err1.Error())
		}
		ds.dictLock.Lock()
		ds.dicts[dictItem.ID] = dictItem
		ds.dictLock.Unlock()
	}
	return nil
}
