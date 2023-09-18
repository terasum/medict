package gincon

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
	"github.com/terasum/medict/internal/utils"
	"github.com/terasum/medict/pkg/model"
)

func (dc *DictsController) GetDictCover(args map[string]interface{}) *model.Resp {
	dictId, ok := args["dict_id"]
	if !ok {
		return model.BuildError(errors.New("dict_id not found"), model.BadParamErrCode)
	}
	coverName, ok := args["cover_name"]
	if !ok {
		return model.BuildError(errors.New("cover_name not found"), model.BadParamErrCode)
	}
	if dict, ok := dc.ds.GetDictPlain(dictId.(string)); ok {
		path := filepath.Join(dict.BaseDictDir, coverName.(string))
		path, _ = filepath.Abs(path)
		log.Infof("cover file path %s", path)
		if utils.FileExists(path) {
			buff, err := os.ReadFile(path)
			if err != nil {
				return model.BuildError(err, model.InnerSysErrCode)
			}
			base64Str := base64.StdEncoding.EncodeToString(buff)
			return model.BuildSuccess(base64Str)
		}
	} else {
		log.Infof("dictionary not found %s", dictId)
	}
	return model.BuildRawError("not exists", model.BadReqCode)
}

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
