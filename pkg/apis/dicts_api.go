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
	"github.com/terasum/medict/internal/config"
	"strings"

	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
)

type DictsAPI struct {
	DictSvc *service.DictService
}

func NewDictsAPI(conf *config.Config) (*DictsAPI, error) {
	dictSvc, err := service.NewDictService(conf)
	if err != nil {
		return nil, err
	}
	return &DictsAPI{
		DictSvc: dictSvc,
	}, nil
}

func (dicts *DictsAPI) Dicts() []*model.PlainDictionaryItem {
	return dicts.DictSvc.Dicts()
}

func (dicts *DictsAPI) Lookup(dictId string, rawKeyWord string) *model.Resp {
	defData, err := dicts.DictSvc.Lookup(dictId, rawKeyWord)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	def := string(defData)
	def = strings.Replace(def, `\"`, `"`, -1)
	def = `<html><body>` + def + `</body></html>`
	return model.BuildSuccess(base64.StdEncoding.EncodeToString([]byte(def)))
}

func (dicts *DictsAPI) Search(dictId string, rawKeyWord string) *model.Resp {
	entries, err := dicts.DictSvc.Search(dictId, rawKeyWord)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}

	return model.BuildSuccess(entries)
}

func (dicts *DictsAPI) Locate(dictId string, entry *model.KeyBlockEntry) *model.Resp {
	if entry == nil {
		return model.BuildError(errors.New("empty entry"), model.BadParamErrCode)
	}
	def, err := dicts.DictSvc.Locate(dictId, entry)
	if err != nil {
		return model.BuildError(err, model.InnerSysErrCode)
	}
	def = strings.Replace(def, `\"`, `"`, -1)
	def = `<html><body>` + def + `</body></html>`
	return model.BuildSuccess(base64.StdEncoding.EncodeToString([]byte(def)))
}
