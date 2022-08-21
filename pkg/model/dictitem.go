package model

import (
	"github.com/terasum/medict/internal/mdictparser"
	"github.com/terasum/medict/internal/utils"
)

type PlainDictItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DictItem struct {
	Info     *PlainDictItem
	PathInfo *DirItem

	MDXIndexer *Indexer
	MDXHandler *mdictparser.MdictParser

	MDDIndexer []*Indexer
	MDDHandler []*mdictparser.MdictParser
}

func NewByDirItem(dirItem *DirItem) (*DictItem, error) {
	dictItem := &DictItem{
		Info: &PlainDictItem{
			ID:   utils.MD5Hash(dirItem.CurrentDir),
			Name: utils.FetchBaseDirName(dirItem.CurrentDir),
		},
		PathInfo:   dirItem,
		MDDHandler: make([]*mdictparser.MdictParser, 0),
		MDDIndexer: make([]*Indexer, 0),
	}

	mdx := mdictparser.NewMDictParser()
	err := mdx.Load(dirItem.MdxAbsPath)
	if err != nil {
		mdx.Destroy()
		return nil, err
	}
	dictItem.MDXHandler = mdx
	dictItem.MDXIndexer = NewIndexer(MDX, mdx)

	for _, mddpath := range dirItem.MddAbsPath {
		mdd := mdictparser.NewMDictParser()
		err := mdd.Load(mddpath)
		if err != nil {
			for _, tmdd := range dictItem.MDDHandler {
				tmdd.Destroy()
			}
			return nil, err
		}
		dictItem.MDDHandler = append(dictItem.MDDHandler, mdd)
		dictItem.MDDIndexer = append(dictItem.MDDIndexer, NewIndexer(MDD, mdd))
	}
	return dictItem, nil

}

func (di *DictItem) SimWord(word string) ([]*WrappedWordItem, error) {
	simwords, err := di.MDXIndexer.SimWord(word, 2)
	if err != nil {
		return nil, err
	}
	res := make([]*WrappedWordItem, len(simwords))
	for i, simw := range simwords {
		res[i] = &WrappedWordItem{
			DictId:   di.Info.ID,
			WordItem: simw,
		}
	}
	return res, err
}

func (di *DictItem) LocateResource(resourceKey string) {

}
