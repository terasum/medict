//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package apis

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
)

type DictsAPI struct {
	dictSvc *service.DictService
}

func NewDictsAPI(config *config.Config) (*DictsAPI, error) {
	svc, err := service.NewDictService(config)
	if err != nil {
		return nil, err
	}
	return &DictsAPI{dictSvc: svc}, nil
}

func (dict *DictsAPI) Dicts() []*model.PlainDictionaryItem {
	return dict.dictSvc.Dicts()
}

func (dict *DictsAPI) Lookup(dictId string, rawKeyWord string) *model.Resp {
	defData, err := dict.dictSvc.Lookup(dictId, rawKeyWord)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	def := string(defData)
	def = strings.Replace(def, `\"`, `"`, -1)
	def = `<html><body>` + def + `</body></html>`
	return model.BuildSuccess(base64.StdEncoding.EncodeToString([]byte(def)))
}

func (dict *DictsAPI) Search(dictId string, rawKeyWord string) *model.Resp {
	entries, err := dict.dictSvc.Search(dictId, rawKeyWord)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}

	return model.BuildSuccess(entries)
}

func (dict *DictsAPI) Locate(dictId string, entry *model.KeyBlockEntry) *model.Resp {
	if entry == nil {
		return model.BuildError(errors.New("empty entry"), model.BadParamErrCode)
	}
	def, err := dict.dictSvc.Locate(dictId, entry)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	def = strings.Replace(def, `\"`, `"`, -1)
	def = `<html><body>` + def + `</body></html>`
	return model.BuildSuccess(base64.StdEncoding.EncodeToString([]byte(def)))
}
