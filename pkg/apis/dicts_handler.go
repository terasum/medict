package apis

import (
	"errors"

	"github.com/terasum/medict/pkg/model"
)

// GetAllDicts returns the list of loaded dictionaries.
func (dc *DictsController) GetAllDicts() *model.Resp {
	return model.BuildSuccess(dc.svc.Dicts())
}

// InitDicts scans the dict directory and loads dictionaries.
func (dc *DictsController) InitDicts() *model.Resp {
	if err := dc.svc.InitDicts(); err != nil {
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
	dict := dc.svc.GetDictById(dictid)
	if err := dict.Dict.BuildIndex(); err != nil {
		log.Infof("[wails] building dictionary index, dictionary path is %s", dict.Dict.Name())
		log.Infof("[wails] building dictionary index failed, err %s", err.Error())
		return model.BuildError(err, model.InnerSysErrCode)
	}
	log.Infof("[wails] building dictionary index success, id: %s", dict.Dict.Name())
	return model.BuildSuccess(dict.Dict.Name())
}

// SearchWord returns the matching keyword entries for a dictionary.
func (dc *DictsController) SearchWord(dictId, word string) *model.Resp {
	entries, err := dc.svc.Search(dictId, word)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(entries)
}
