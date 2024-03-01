package utils

import (
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/pkg/code"
)

func Ok(msg string) *v1.BaseResp {
	return &v1.BaseResp{
		Code: code.SUCCESS,
		Msg:  msg,
	}
}

func Fail(msg string) *v1.BaseResp {
	return &v1.BaseResp{
		Code: code.ERROR,
		Msg:  msg,
	}
}

func Warn(msg string) *v1.BaseResp {
	return &v1.BaseResp{
		Code: code.WARN,
		Msg:  msg,
	}
}
