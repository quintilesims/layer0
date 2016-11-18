package errors

import (
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

type ServerError struct {
	Code ErrorCode
	Err  error
}

func Newf(code ErrorCode, format string, tokens ...interface{}) *ServerError {
	return New(code, fmt.Errorf(format, tokens...))
}

func New(code ErrorCode, err error) *ServerError {
	return &ServerError{
		Code: code,
		Err:  err,
	}
}

func (this *ServerError) Error() string {
	return fmt.Sprintf("ServerError (code=%d) %s", this.Code, this.Err.Error())
}

func (this *ServerError) Model() models.ServerError {
	return models.ServerError{
		ErrorCode: int64(this.Code),
		Message:   this.Err.Error(),
	}
}
