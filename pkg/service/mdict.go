package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/terasum/medict/pkg/model"
)

var _ model.GeneralDictionary = &Mdict{}

type Mdict struct {
	mdxFilePath       string
	mddFilePaths      []string
	mdx               *virtualMdict
	mdxIdxFilePath    string
	mdds              []*virtualMdict
	mddsIdxFilePaths  []string
	hasBuildIndex     bool
	buildingIndexLock *sync.Mutex
}

func NewMdict(dirItem *model.DirItem) (model.GeneralDictionary, error) {
	mdict := &Mdict{
		mdxFilePath:       dirItem.MdictMdxAbsPath,
		mddFilePaths:      dirItem.MdictMddAbsPath,
		hasBuildIndex:     false,
		mddsIdxFilePaths:  make([]string, len(dirItem.MdictMddAbsPath)),
		buildingIndexLock: new(sync.Mutex),
	}

	mdx, err := newVirtual(dirItem.MdictMdxAbsPath)
	if err != nil {
		return nil, fmt.Errorf("new mdx file failed, %s", err.Error())
	}

	mdds := make([]*virtualMdict, 0)
	for _, mddpath := range dirItem.MdictMddAbsPath {
		mdd, err1 := newVirtual(mddpath)
		if err1 != nil {
			return nil, fmt.Errorf("new mdd file failed, %s", err1.Error())
		}
		mdds = append(mdds, mdd)
	}

	mdict.mdx = mdx
	mdict.mdds = mdds

	return mdict, nil

}

func (md *Mdict) Name() string {
	_, rawpath := filepath.Split(md.mdxFilePath)
	rawpath = strings.TrimRight(rawpath, ".mdx")
	if len(rawpath) > 0 {
		return rawpath
	}
	return md.mdxFilePath
}

func (md *Mdict) Description() *model.PlainDictionaryInfo {
	if md.mdx == nil {
		return &model.PlainDictionaryInfo{}
	}
	return md.mdx.description()
}

func (md *Mdict) KeyList() []string {
	return md.mdx.keyList()
}

func (md *Mdict) BuildIndex() error {
	md.buildingIndexLock.Lock()
	defer md.buildingIndexLock.Unlock()
	if md.hasBuildIndex {
		return nil
	}

	mdxIndexFilePath := md.mdx.filePath + ".meidx"
	err := md.mdx.index(mdxIndexFilePath)
	if err != nil {
		return err
	}
	md.mdxIdxFilePath = mdxIndexFilePath

	for i, mdd := range md.mdds {
		mddIndexFilePath := mdd.filePath + ".meidx"
		err1 := mdd.index(mddIndexFilePath)
		if err1 != nil {
			return err1
		}
		md.mddsIdxFilePaths[i] = mddIndexFilePath
	}

	md.hasBuildIndex = true
	return nil
}

func (md *Mdict) Locate(entry *model.KeyIndex) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	return md.mdx.locate(entry)
}

func (md *Mdict) DictType() model.DictType {
	return model.DictTypeMdict
}

func (md *Mdict) Lookup(keyword string) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	return md.mdx.lookup(keyword)
}

func (md *Mdict) LookupResource(keyword string) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	var err error
	var def []byte

	for _, mdd := range md.mdds {
		def, err = mdd.lookup(keyword)
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

	// search from index file
	return md.mdx.searchFromIndex(keyword)
}
