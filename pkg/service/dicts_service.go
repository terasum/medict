package service

import (
	"errors"
	"fmt"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/gomdict"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service/support"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"os"
	"path"
)

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
		fmt.Printf("dict %s, currentDir is %s\n", dictId, dict.PathInfo.CurrentDir)
		fullPath := path.Join(dict.PathInfo.CurrentDir, key)
		if fileutil.Exist(fullPath) {
			return os.ReadFile(fullPath)
		}
	}
	return nil, errors.New("not found")
}

func (ds *DictService) Dicts() []*model.PlainDictionaryItem {
	result := make([]*model.PlainDictionaryItem, len(ds.dicts))
	i := 0
	for _, dict := range ds.dicts {
		result[i] = &model.PlainDictionaryItem{
			ID:   dict.ID,
			Name: dict.Name,
			Path: dict.PathInfo.MdxAbsPath,
		}
		i++
	}
	return result
}

func (ds *DictService) Lookup(dictId string, keyword string) ([]byte, error) {
	if dict, ok := ds.dicts[dictId]; !ok {
		return nil, errors.New("dict not found")
	} else {
		data, err := dict.MDX.Lookup(keyword)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (ds *DictService) LookupResource(dictId string, keyword string) ([]byte, error) {
	if dict, ok := ds.dicts[dictId]; !ok {
		fmt.Printf("MDD resource search [%s] dict not found\n", keyword)
		return nil, fmt.Errorf("resource (%s) not found", keyword)
	} else {
		for _, mdd := range dict.MDDS {
			data, err := mdd.Lookup(keyword)
			if err != nil {
				fmt.Printf("MDD resource search (%s)[%s] failed %s\n", mdd.FilePath, keyword, err.Error())
				continue
			}
			fmt.Printf("MDD resource search (%s)[%s] success\n", mdd.FilePath, keyword)
			return data, nil
		}
	}
	fmt.Printf("MDD resource search [%s] failed, not found\n", keyword)
	return nil, fmt.Errorf("resource (%s) not found", keyword)
}

func (ds *DictService) Locate(dictid string, entry *model.KeyBlockEntry) (string, error) {
	if dict, ok := ds.dicts[dictid]; !ok {
		return "", errors.New("dict not found")
	} else {
		mdictEntry := &gomdict.MDictKeyBlockEntry{
			RecordStartOffset: entry.RecordStartOffset,
			RecordEndOffset:   entry.RecordEndOffset,
			KeyWord:           entry.KeyWord,
			KeyBlockIdx:       entry.KeyBlockIdx,
		}
		defData, err := dict.MDX.Locate(mdictEntry)
		if err != nil {
			return "", err
		}
		return string(defData), nil
	}
}

func (ds *DictService) Search(dictId string, keyword string) ([]*model.KeyBlockEntry, error) {
	if dict, ok := ds.dicts[dictId]; !ok {
		return nil, errors.New("dict not found")
	} else {
		entries, err := dict.MDX.Search(keyword)
		if err != nil {
			return nil, err
		}
		results := make([]*model.KeyBlockEntry, 0)
		for id, e := range entries {
			temp := &model.KeyBlockEntry{
				ID:                id,
				RecordStartOffset: e.RecordStartOffset,
				RecordEndOffset:   e.RecordEndOffset,
				KeyWord:           e.KeyWord,
				KeyBlockIdx:       e.KeyBlockIdx,
			}
			results = append(results, temp)

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
		return err
	}
	items, err := support.WalkDir(baseDir)
	if err != nil {
		return err
	}
	for _, dirItem := range items {
		dictItem, err := model.NewByDirItem(dirItem)
		if err != nil {
			return err
		}
		ds.dicts[dictItem.ID] = dictItem
	}
	return nil
}
