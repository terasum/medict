package service

import (
	"github.com/terasum/medict/pkg/model"
)

var _ model.GeneralDictionary = &StarDict{}

type StarDict struct {
	DzFilePath  string
	IfoFilePath string
	IdxFilePath string
}

func (s StarDict) Name() string {
	panic("implement me")
}

func (s StarDict) Description() *model.PlainDictionaryInfo {
	panic("implement me")
}

func (s StarDict) BuildIndex() error {
	//TODO implement me
	panic("implement me")
}

func (s StarDict) DictType() model.DictType {
	//TODO implement me
	panic("implement me")
}

func (s StarDict) Lookup(keyword string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s StarDict) LookupResource(keyword string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s StarDict) Locate(entry *model.KeyIndex) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s StarDict) Search(keyword string) ([]*model.KeyIndex, error) {
	//TODO implement me
	panic("implement me")
}

func NewStardict(dirItem *model.DirItem) (model.GeneralDictionary, error) {

	return StarDict{}, nil

}
