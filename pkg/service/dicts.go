package service

import (
	"errors"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service/support"
)

type DictService struct {
	config *config.Config
	dicts  map[string]*model.DictItem
}

func NewDictService(config *config.Config) (*DictService, error) {
	ds := &DictService{
		config: config,
		dicts:  make(map[string]*model.DictItem),
	}
	err := ds.walkDicts()
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// GC IMPORTANT: you have to invoke this to reduce the memory
func (ds *DictService) GC() {
	for _, dict := range ds.dicts {
		dict.MDXHandler.Destroy()
		for _, mdd := range dict.MDDHandler {
			mdd.Destroy()
		}
	}
}

func (ds *DictService) Dicts() []*model.PlainDictItem {
	result := make([]*model.PlainDictItem, len(ds.dicts))
	i := 0
	for _, dict := range ds.dicts {
		result[i] = dict.Info
		i++
	}
	return result
}

func (ds *DictService) LookupDefinition(dictId string, rawKeyWord string, recordStart uint64) (string, error) {
	if dict, ok := ds.dicts[dictId]; !ok {
		return "", errors.New("dict not found")
	} else {
		def, err := dict.MDXHandler.FindDef(rawKeyWord, recordStart)
		if err != nil {
			return "", err
		}
		return def, nil
	}
}

func (ds *DictService) SimWords(word string) ([]*model.WrappedWordItem, error) {
	result := make([]*model.WrappedWordItem, 0)
	for _, dict := range ds.dicts {
		simWords, err := dict.SimWord(word)
		if err != nil {
			return nil, err
		}
		result = append(result, simWords...)
	}
	return result, nil
}

func (ds *DictService) AllWordIndexing() (uint64, uint64, error) {
	mdxsize, mddsize := uint64(0), uint64(0)
	var err error
	for _, dict := range ds.dicts {
		err, mdxsize = dict.MDXIndexer.Indexing()
		if err != nil {
			return 0, 0, err
		}
		for _, mdxIdx := range dict.MDDIndexer {
			err, mdds := mdxIdx.Indexing()
			if err != nil {
				return 0, 0, err
			}
			mddsize += mdds
		}
	}

	return mdxsize, mddsize, nil
}

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
		ds.dicts[dictItem.Info.ID] = dictItem
	}
	return nil
}
