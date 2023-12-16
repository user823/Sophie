package core

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"net/http"
)

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func WriteResponse(ctx *app.RequestContext, err error, data any) {
	// 有err时不一定表示响应失败
	if err != nil {
		log.Debugf("%#+v", err)
		coder := errors.ParseCoder("HttpResponse", err)
		ctx.JSON(coder.Code(), ErrResponse{
			Code:    coder.Code(),
			Message: coder.Message(),
			Data:    data,
		})
		return
	}
	// 没有错误时返回成功
	ctx.JSON(http.StatusOK, ErrResponse{
		Code: http.StatusOK,
		Data: data,
	})
}
