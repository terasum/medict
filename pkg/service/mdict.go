package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/terasum/medict/internal/gomdict"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
)

var _ model.GeneralDictionary = &Mdict{}

type Mdict struct {
	mdxFilePath  string
	mddFilePaths []string
	mdxins       *gomdict.Mdict
	mddinss      []*gomdict.Mdict

	hasBuildIndex     bool
	buildingIndexLock *sync.Mutex
}

func NewMdict(dirItem *model.DirItem) (model.GeneralDictionary, error) {
	mdict := &Mdict{
		mdxFilePath:       dirItem.MdictMdxAbsPath,
		mddFilePaths:      dirItem.MdictMddAbsPath,
		hasBuildIndex:     false,
		buildingIndexLock: new(sync.Mutex),
	}

	mdx, err := gomdict.New(dirItem.MdictMdxAbsPath)
	if err != nil {
		return nil, err
	}

	mdds := make([]*gomdict.Mdict, 0)

	for _, mddpath := range dirItem.MdictMddAbsPath {
		mdd, err1 := gomdict.New(mddpath)
		if err1 != nil {
			return nil, err1
		}
		mdds = append(mdds, mdd)
	}

	mdict.mdxins = mdx
	mdict.mddinss = mdds

	return mdict, nil

}

func (md *Mdict) Name() string {
	return strings.TrimRight(utils.FileName(md.mdxFilePath), ".mdx")
}

func (md *Mdict) Description() *model.PlainDictionaryInfo {
	if md.mdxins == nil {
		return &model.PlainDictionaryInfo{}
	}
	return &model.PlainDictionaryInfo{
		Title:                 md.mdxins.Meta.Title,
		Description:           md.mdxins.Meta.Description,
		CreateDate:            md.mdxins.Meta.CreationDate,
		GenerateEngineVersion: md.mdxins.Meta.GeneratedByEngineVersion,
	}
}

func (md *Mdict) BuildIndex() error {
	md.buildingIndexLock.Lock()
	defer md.buildingIndexLock.Unlock()
	if md.hasBuildIndex {
		return nil
	}

	err := md.mdxins.BuildIndex()
	if err != nil {
		return err
	}

	for _, mdd := range md.mddinss {
		err1 := mdd.BuildIndex()
		if err1 != nil {
			return err1
		}
	}

	md.hasBuildIndex = true
	return nil
}

func (md *Mdict) Locate(entry *model.KeyIndex) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	mdictEntry := &gomdict.MDictKeyBlockEntry{
		RecordStartOffset: entry.RecordStartOffset,
		RecordEndOffset:   entry.RecordEndOffset,
		KeyWord:           entry.KeyWord,
		KeyBlockIdx:       entry.KeyBlockIdx,
	}
	def, err := md.mdxins.Locate(mdictEntry)
	return def, err
}

func (md *Mdict) DictType() model.DictType {

	return model.DictTypeMdict
}

func (md *Mdict) Lookup(keyword string) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	def, err := md.mdxins.Lookup(keyword)
	return def, err
}

func (md *Mdict) LookupResource(keyword string) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	var err error
	var def []byte

	for _, mdd := range md.mddinss {
		def, err = mdd.Lookup(keyword)
		if err != nil {
			log.Infof("mdict.LookupResource failed, key [%s] not found", keyword)
			continue
		} else {
			break
		}
	}

	if def == nil || err != nil {
		if err != nil {
			log.Infof("mdict.LookupResource failed, key [%s] not found, error: %s", keyword, err.Error())
		}
		return nil, fmt.Errorf("mdict resource not found: [%s]", keyword)
	}

	return def, nil
}

func (md *Mdict) Search(keyword string) ([]*model.KeyIndex, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}

	entries, err := md.mdxins.Search(keyword)
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
