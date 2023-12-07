package apis

import (
	"errors"
	"github.com/terasum/medict/pkg/model"
)

// GetAllDicts
func (dc *DictsController) GetAllDicts(args map[string]interface{}) *model.Resp {
	dicts := dc.ds.Dicts()
	return model.BuildSuccess(dicts)
}

// InitDicts 初始化词典
func (dc *DictsController) InitDicts(args map[string]interface{}) *model.Resp {
	err := dc.ds.InitDicts()
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(nil)
}

// buildIndex
func (dc *DictsController) BuildIndexByDictId(args map[string]interface{}) *model.Resp {
	if id, ok := args["dictid"]; !ok {
		return model.BuildError(errors.New("build index failed, dictid is empty"), model.InnerSysErrCode)
	} else {
		log.Infof("[wails] building dictionary index, dict id is %s", id)
		dict := dc.ds.GetDictById(id.(string))
		if err := dict.MainDict.BuildIndex(); err != nil {
			log.Infof("[wails] building dictionary index, dictionary path is %s", dict.MainDict.Name())
			log.Infof("[wails] building dictionary index failed, err %s", err.Error())
			return model.BuildError(err, model.InnerSysErrCode)
		} else {
			log.Infof("[wails] building dictionary index success, id: %s", dict.MainDict.Name())
			return model.BuildSuccess(dict.MainDict.Name())
		}
	}
}

// SearchWord
func (dc *DictsController) SearchWord(args map[string]interface{}) *model.Resp {
	dictId, ok := args["dict_id"]
	if !ok {
		return model.BuildError(errors.New("dict_id not found"), model.BadParamErrCode)
	}

	word, ok := args["word"]
	if !ok {
		return model.BuildError(errors.New("dict_id not found"), model.BadParamErrCode)
	}

	entries, err := dc.ds.Search(dictId.(string), word.(string))
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	return model.BuildSuccess(entries)
}
