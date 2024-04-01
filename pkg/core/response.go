package core

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"net/http"
)

type ErrResponse struct {
	Code     int           `json:"code"`
	Message  string        `json:"msg,omitempty"`
	Data     any           `json:"data,omitempty"`
	BaseResp *BaseResponse `json:"baseResp"`
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg,omitempty"`
}

func JSON(c *app.RequestContext, data any) {
	c.JSON(200, data)
}

func WriteResponseE(c *app.RequestContext, err error, data any) {
	// 有err时不一定表示响应失败
	if err != nil {
		log.Debugf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.Code(), ErrResponse{
			Code:    coder.Code(),
			Message: coder.Message(),
			Data:    data,
		})
		return
	}
	// 没有错误时返回成功
	c.JSON(http.StatusOK, ErrResponse{
		Code: http.StatusOK,
		Data: data,
	})
}

func WriteResponse(c *app.RequestContext, response ErrResponse) {
	c.JSON(200, response)
}

func OK(c *app.RequestContext, msg string, data interface{}) {
	WriteResponse(c, ErrResponse{
		Code:    api.SUCCESS,
		Message: msg,
		Data:    data,
	})
}

func Fail(c *app.RequestContext, msg string, data interface{}) {
	WriteResponse(c, ErrResponse{
		Code:    api.FAIL,
		Message: msg,
		Data:    data,
	})
}
