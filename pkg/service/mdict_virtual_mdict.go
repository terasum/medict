package service

import (
	"errors"
	"fmt"
	"github.com/terasum/medict/internal/gomdict"
	"github.com/terasum/medict/internal/medengine"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
	"path/filepath"
)

type virtualMdict struct {
	instance  *gomdict.Mdict
	idxEngine *medengine.IndexEngine
	filePath  string
	dictType  string
}

const (
	virtualTypeMdx = "dictTypeMdx"
	virtualTypeMdd = "dictTypeMdd"
)

func newVirtual(filePath string) (*virtualMdict, error) {

	mdx, err := gomdict.New(filePath)
	if err != nil {
		return nil, fmt.Errorf("new mdx file failed, %s", err.Error())
	}

	ext := filepath.Ext(filePath)
	var dt string
	if ext == ".mdx" {
		dt = virtualTypeMdx
	} else if ext == ".mdd" {
		dt = virtualTypeMdd
	}

	return &virtualMdict{
		instance:  mdx,
		idxEngine: nil,
		filePath:  filePath,
		dictType:  dt,
	}, nil
}

func (vm *virtualMdict) description() *model.PlainDictionaryInfo {
	if vm.instance == nil {
		return &model.PlainDictionaryInfo{}
	}
	return &model.PlainDictionaryInfo{
		Title:                 vm.instance.Meta.Title,
		Description:           vm.instance.Meta.Description,
		CreateDate:            vm.instance.Meta.CreationDate,
		GenerateEngineVersion: vm.instance.Meta.GeneratedByEngineVersion,
	}
}

func (vm *virtualMdict) keyList() []string {
	result := make([]string, 0)
	for _, e := range vm.instance.KeyBlockData.KeyEntries {
		result = append(result, e.KeyWord)
	}

	return result
}

func (vm *virtualMdict) locate(entry *model.KeyIndex) ([]byte, error) {
	mdictEntry := &gomdict.MDictKeyBlockEntry{
		RecordStartOffset: entry.RecordStartOffset,
		RecordEndOffset:   entry.RecordEndOffset,
		KeyWord:           entry.KeyWord,
		KeyBlockIdx:       entry.KeyBlockIdx,
	}
	return vm.instance.Locate(mdictEntry)
}

func (vm *virtualMdict) lookup(keyword string) ([]byte, error) {
	return vm.instance.Lookup(keyword)
}

func (vm *virtualMdict) searchFromIndex(keyword string) ([]*model.KeyIndex, error) {
	if vm.idxEngine == nil {
		return nil, errors.New("virtual mdict hasn't built the index")
	}

	records, err := vm.idxEngine.Search(keyword)
	if err != nil {
		return nil, err
	}
	results := make([]*model.KeyIndex, 0)
	for idx, record := range records {
		kblockEntry := record.ToKeyBlockEntry()
		kblockEntry.ID = idx
		results = append(results, &model.KeyIndex{
			IndexType:     model.IndexTypeMdict,
			KeyBlockEntry: kblockEntry,
		})
	}

	return results, nil

}

func (vm *virtualMdict) search(keyword string) ([]*model.KeyIndex, error) {

	entries, err := vm.instance.Search(keyword)
	if err != nil {
		return nil, err
	}

	results := make([]*model.KeyIndex, 0)

	for id, e := range entries {
		temp := &model.KeyBlockEntry{
			ID:                id,
			RecordStartOffset: e.RecordStartOffset,
			RecordEndOffset:   e.RecordEndOffset,
			KeyWord:           e.KeyWord,
			KeyBlockIdx:       e.KeyBlockIdx,
		}
		tempIdx := &model.KeyIndex{
			IndexType:     model.IndexTypeMdict,
			KeyBlockEntry: temp,
		}
		results = append(results, tempIdx)

	}
	return results, nil
}

func (vm *virtualMdict) index(idxFilePath string) error {
	var eng *medengine.IndexEngine
	var err error
	//err = vm.instance.BuildIndex()
	//if err != nil {
	//	return err
	//}

	// has already built index
	if utils.FileExists(idxFilePath) {
		eng, err = medengine.NewEngine(idxFilePath)
		if err != nil {
			return err
		}
		vm.idxEngine = eng
		return nil
	}

	// TODO ignore bktree index
	// err = md.mdxins.BuildBKTree()
	// if err != nil {
	// 	return err
	// }

	eng, err = medengine.NewEngine(idxFilePath)
	if err != nil {
		return err
	}

	for _, e := range vm.instance.KeyBlockData.KeyEntries {
		err = eng.AddRecord(medengine.NewIndexRecord(
			e.KeyWord,
			e.RecordStartOffset,
			e.RecordEndOffset,
			e.KeyBlockIdx,
		))
		if err != nil {
			continue
		}
	}

	vm.idxEngine = eng
	return nil
}
