package model

const (
	SuccessErrCode   = 200
	InnerSysErrCode  = 500
	BadParamErrCode  = 401
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

func BuildError(err error, code int) *Resp {
	return &Resp{
		Data: nil,
		Err:  err.Error(),
		Code: code,
	}
}
