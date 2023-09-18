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

package model

const (
	SuccessErrCode   = 200
	InnerSysErrCode  = 500
	BadParamErrCode  = 401
	BadReqCode       = 400
	UnknownErrorCode = 999
)

type Resp struct {
	Data interface{} `json:"data"`
	Err  string      `json:"err"`
	Code int         `json:"code"`
}

func BuildSuccess(data interface{}) *Resp {
	return &Resp{
		Data: data,
		Err:  "success",
		Code: SuccessErrCode,
	}
}

func BuildRawError(err string, code int) *Resp {
	return &Resp{
		Data: nil,
		Err:  err,
		Code: code,
	}
}

func BuildError(err error, code int) *Resp {
	return &Resp{
		Data: nil,
		Err:  err.Error(),
		Code: code,
	}
}
