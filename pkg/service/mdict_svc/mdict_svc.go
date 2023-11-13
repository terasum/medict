package mdict_svc

import (
	"errors"
	"fmt"
	"github.com/terasum/medict/pkg/model"
)

var _ model.GeneralDictionary = &MdictSvcImpl{}

type MdictSvcImpl struct {
	hasBuildIndex bool
	mdx           *mdictHolder
	mdds          []*mdictHolder
}

func NewMdictSvc(dirItem *model.DirItem) (model.GeneralDictionary, error) {
	mdx, err := newMdictHolder(dirItem.MdictMdxAbsPath)
	if err != nil {
		return nil, fmt.Errorf("new mdx file failed, %s", err.Error())
	}

	mdds := make([]*mdictHolder, 0)
	for _, mddpath := range dirItem.MdictMddAbsPath {
		mdd, err1 := newMdictHolder(mddpath)
		if err1 != nil {
			return nil, fmt.Errorf("new mdd file failed, %s", err1.Error())
		}
		mdds = append(mdds, mdd)
	}
	mdict := &MdictSvcImpl{
		mdx:  mdx,
		mdds: mdds,
	}
	return mdict, nil
}

func (md *MdictSvcImpl) DictType() model.DictType {
	return model.DictTypeMdict
}

func (md *MdictSvcImpl) Name() string {
	return md.mdx.rawdict.Name()
}

func (md *MdictSvcImpl) Description() *model.PlainDictionaryInfo {
	return &model.PlainDictionaryInfo{
		Title:                 "",
		Description:           "",
		CreateDate:            "",
		GenerateEngineVersion: "",
	}
}

func (md *MdictSvcImpl) KeyList() []string {
	return nil
}

func (md *MdictSvcImpl) BuildIndex() error {
	err := md.mdx.BuildIndex()
	if err != nil {
		return err
	}

	for _, mdd := range md.mdds {
		err = mdd.BuildIndex()
		if err != nil {
			return err
		}
	}
	return nil
}

func (md *MdictSvcImpl) Locate(entry *model.KeyQueryIndex) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	return nil, nil
	//return md.mdx.rawdict.LocateByKeywordIndex(entry)
}

func (md *MdictSvcImpl) Lookup(keyword string) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	return md.mdx.rawdict.Lookup(keyword)
}

func (md *MdictSvcImpl) LookupResource(keyword string) ([]byte, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}
	var err error
	var def []byte
	for _, mdd := range md.mdds {
		def, err = mdd.rawdict.Lookup(keyword)
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

func (md *MdictSvcImpl) Search(keyword string) ([]*model.KeyQueryIndex, error) {
	if !md.hasBuildIndex {
		return nil, errors.New("dictionary not ready, building index first")
	}

	_, err := md.mdx.idxer.Search(keyword)
	if err != nil {
		return nil, err
	}
	//return indeses, nil
	return nil, nil

}
