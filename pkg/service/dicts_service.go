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
	"github.com/terasum/medict/pkg/service/onlinedict"
	"github.com/terasum/medict/pkg/service/support"
)

type DictService struct {
	config   *config.Config
	dicts    map[string]*model.DictionaryItem
	dictLock *sync.Mutex
}

// NewDictService constructs a DictService. It is a plain constructor — no
// process-wide singleton — so the caller (App) owns the instance and injects
// it where needed (issue #727).
func NewDictService(config *config.Config) (*DictService, error) {
	return &DictService{
		config:   config,
		dicts:    make(map[string]*model.DictionaryItem),
		dictLock: new(sync.Mutex),
	}, nil
}

// Close releases resources held by all loaded dictionaries (e.g. leveldb
// handles). Intended for app shutdown: errors from individual dicts are logged
// and joined, and never block closing the rest.
func (ds *DictService) Close() error {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	var errs []error
	for _, dict := range ds.dicts {
		if dict.Dict == nil {
			continue
		}
		if err := dict.Dict.Close(); err != nil {
			log.Errorf("close dict %s failed: %s", dict.ID, err.Error())
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// InitDicts initialize the dictionaries
func (ds *DictService) InitDicts() error {
	if err := ds.walkDicts(); err != nil {
		return err
	}
	ds.registerVirtualDicts()
	return nil
}

// registerVirtualDicts adds non-file (online) dictionaries to the dict map, so
// they appear in the dict list like file dicts. Currently: Bing (free, no key).
func (ds *DictService) registerVirtualDicts() {
	bing := onlinedict.NewBing()
	item := &model.DictionaryItem{
		PlainDictionaryItem: &model.PlainDictionaryItem{
			ID:          "online-bing",
			Name:        bing.Name(),
			DictType:    string(model.DictTypeOnline),
			Description: bing.Description(),
		},
		Dict: bing,
		// PathInfo 留 nil:在线词典无目录。
	}
	ds.dictLock.Lock()
	ds.dicts[item.ID] = item
	ds.dictLock.Unlock()
}

func (ds *DictService) FindFromDir(dictId string, key string) ([]byte, error) {
	ds.dictLock.Lock()
	defer ds.dictLock.Unlock()

	if dict, ok := ds.dicts[dictId]; ok {
		if dict.PathInfo == nil {
			return nil, errors.New("virtual dict has no directory")
		}
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
		err := dict.Dict.BuildIndex()
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
	// 锁内仅从 dicts map 取出 *model.DictionaryItem 指针（快），
	// 慢 I/O（dict.Dict.Lookup：磁盘/解压）放到锁外执行，避免持锁序列化。
	ds.dictLock.Lock()
	dict, ok := ds.dicts[dictId]
	ds.dictLock.Unlock()

	if !ok {
		return nil, errors.New("dict not found")
	}
	data, err := dict.Dict.Lookup(keyword)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ds *DictService) LookupResource(dictId string, keyword string) ([]byte, error) {
	// 锁内仅取出 dict 指针，资源查询（磁盘/解压）在锁外进行。
	ds.dictLock.Lock()
	dict, ok := ds.dicts[dictId]
	ds.dictLock.Unlock()

	if !ok {
		log.Infof("LookResource dict not found [%s]", keyword)
		return nil, fmt.Errorf("dictionary (%s) not found", keyword)
	}
	keyword = strings.TrimSpace(keyword)
	data, err := dict.Dict.LookupResource(keyword)
	if err != nil {
		log.Infof("LookupResource search (%s):[%s] failed, err: %s\n", dict.ToPlain().Name, keyword, err.Error())
		return nil, err
	}
	log.Infof("LookupResource search  (%s)[%s] success\n", dict.ToPlain().Name, keyword)
	return data, nil
}

func (ds *DictService) Locate(dictid string, idx *model.KeyQueryIndex) (string, error) {
	// 锁内仅取出 dict 指针；Locate 内部的磁盘/解压读取在锁外执行。
	ds.dictLock.Lock()
	dict, ok := ds.dicts[dictid]
	ds.dictLock.Unlock()

	if !ok {
		return "", errors.New("dict not found")
	}
	idxType := model.IndexTypeMdict
	if dict.DictType == (string)(model.DictTypeStarDict) {
		idxType = model.IndexTypeStardict
	}
	defData, err := dict.Dict.Locate(&model.KeyQueryIndex{
		IndexType:         idxType,
		MdictKeyWordIndex: idx.MdictKeyWordIndex,
	})
	if err != nil {
		return "", err
	}
	return string(defData), nil
}

func (ds *DictService) Search(dictId string, keyword string) ([]*model.KeyQueryIndex, error) {
	log.Infof("search %s %s", dictId, keyword)
	// 锁内仅取出 dict 指针；Search（遍历索引、构建新切片）在锁外执行。
	ds.dictLock.Lock()
	dict, ok := ds.dicts[dictId]
	ds.dictLock.Unlock()

	if !ok {
		return nil, errors.New("dict not found")
	}
	results, err := dict.Dict.Search(keyword)
	if err != nil {
		return nil, err
	}
	return results, nil
}

/**
 * walkDicts walk the dict dir and load all dicts
 * @return error
 */
func (ds *DictService) walkDicts() error {
	baseDir, err := ds.config.EnsureDictsDir()
	if err != nil {
		return fmt.Errorf("ensure dicts dir failed: %s", err.Error())
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
		ds.dictLock.Lock()
		ds.dicts[dictItem.ID] = dictItem
		ds.dictLock.Unlock()
	}
	return nil
}
