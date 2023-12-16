package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var unknownCoder = defaultCoder{code: 0, http: http.StatusInternalServerError, msg: "An internal server error occurred"}

// 将业务码、Http状态码、message绑定
type Coder interface {
	// 业务码关联到的HTTP状态码
	HTTPStatus() int
	// 业务码（和Http状态码可以相同）
	Code() int
	// 业务码描述信息
	Message() string
}

type defaultCoder struct {
	code int
	http int
	msg  string
}

func (c defaultCoder) Code() int {
	return c.code
}

func (c defaultCoder) HTTPStatus() int {
	return c.http
}

func (c defaultCoder) Message() string {
	return c.msg
}

var (
	codes sync.Map
)

func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by `github.com/Sophie/pkg/errors` as unknownCode error code")
	}

	codes.Store(coder.Code(), coder)
}

func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by `github.com/Sophie/pkg/errors` as unknownCode error code")
	}

	if _, ok := codes.LoadOrStore(coder.Code(), coder); ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}
}

func ParseCoder(namespace string, err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCodeMessage); ok {
		if coder, ok := codes.Load(v.code); ok {
			return coder.(Coder)
		}
	}
	return unknownCoder
}

// 判断error 是否目标code类型
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCodeMessage); ok {
		if v.code == code {
			return true
		}

		if v.error != nil {
			return IsCode(v.error, code)
		}

		return false
	}
	return false
}

func init() {
	codes.Store(unknownCoder.Code(), unknownCoder)
}
