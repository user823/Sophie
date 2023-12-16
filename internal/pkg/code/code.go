package code

import "github.com/user823/Sophie/pkg/errors"

type ErrCode struct {
	C    int
	HTTP int
	Msg  string
}

func (coder ErrCode) Code() int {
	return coder.C
}

func (coder ErrCode) HTTPStatus() int {
	return coder.HTTP
}

func (coder ErrCode) Message() string {
	return coder.Msg
}

var _ errors.Coder = &ErrCode{}

func register(code int, httpStatus int, message string) {
	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Msg:  message,
	}
	errors.MustRegister(coder)
}
