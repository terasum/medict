package apis

import (
	"errors"

	"github.com/terasum/medict/pkg/model"
)

// GetAllDicts returns the list of loaded dictionaries.
func (dc *DictsController) GetAllDicts() *model.Resp {
	return model.BuildSuccess(dc.ds.Dicts())
}

// InitDicts scans the dict directory and loads dictionaries.
func (dc *DictsController) InitDicts() *model.Resp {
	if err := dc.ds.InitDicts(); err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(nil)
}

// BuildIndexByDictId builds the keyword index for one dictionary.
func (dc *DictsController) BuildIndexByDictId(dictid string) *model.Resp {
	if dictid == "" {
		return model.BuildError(errors.New("build index failed, dictid is empty"), model.InnerSysErrCode)
	}
	log.Infof("[wails] building dictionary index, dict id is %s", dictid)
	dict := dc.ds.GetDictById(dictid)
	if err := dict.MainDict.BuildIndex(); err != nil {
		log.Infof("[wails] building dictionary index, dictionary path is %s", dict.MainDict.Name())
		log.Infof("[wails] building dictionary index failed, err %s", err.Error())
		return model.BuildError(err, model.InnerSysErrCode)
	}
	log.Infof("[wails] building dictionary index success, id: %s", dict.MainDict.Name())
	return model.BuildSuccess(dict.MainDict.Name())
}

// SearchWord returns the matching keyword entries for a dictionary.
func (dc *DictsController) SearchWord(dictId, word string) *model.Resp {
	entries, err := dc.ds.Search(dictId, word)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(entries)
}
