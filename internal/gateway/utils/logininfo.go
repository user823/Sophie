package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/pkg/errors"
)

func GetLoginInfoFromCtx(c *app.RequestContext) (v1.LoginUser, error) {
	data, ok := c.Get(api.LOGIN_INFO_KEY)
	if !ok {
		return v1.LoginUser{}, errors.New("获取登录信息失败")
	}
	if info, ok := data.(v1.LoginUser); ok {
		return info, nil
	}
	return v1.LoginUser{}, errors.New("获取登录信息失败")
}
