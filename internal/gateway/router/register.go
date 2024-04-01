package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	code2 "github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/internal/pkg/middleware/auth"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/log"
)

func Register(ctx context.Context, c *app.RequestContext) {
	var register loginInfo
	if err := c.BindAndValidate(&register); err != nil {
		core.Fail(c, validLoginErrMsg(err), nil)
		return
	}

	// 注册用户信息
	password, _ := auth.Encrypt(register.Password) // 加密
	resp, err := rpc.Remoting.RegisterSysUser(ctx, &v1.RegisterSysUserRequest{
		UserInfo: &v1.UserInfo{
			UserName: register.Username,
			Password: password,
			NickName: register.Username,
			Avatar:   v12.AVATAR_URL,
		},
	})
	if err != nil || resp.BaseResp.Code != code2.SUCCESS {
		log.Errorf("Register user error: %s", err.Error())
		return
	}

	core.OK(c, "注册成功", nil)
}
