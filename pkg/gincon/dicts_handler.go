package gincon

import (
	"errors"
	"github.com/terasum/medict/pkg/model"
)

// GetAllDicts
func (dc *DictsController) GetAllDicts(args map[string]interface{}) *model.Resp {
	resp := dc.ds.Dicts()
	return model.BuildSuccess(resp)
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
