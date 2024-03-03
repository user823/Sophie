package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
)

type OnlineUserController struct{}

func NewOnlineUserController() *OnlineUserController {
	return &OnlineUserController{}
}

type onlineUserParam struct {
	v1.PageInfo
	v1.UserOnlineInfo
}

func (o *OnlineUserController) List(ctx context.Context, c *app.RequestContext) {
	var req onlineUserParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListSysUserOnlines(ctx, &v1.ListSysUserOnlinesRequest{
		PageInfo: &req.PageInfo,
		UserName: req.UserName,
		Ipaddr:   req.Ipaddr,
		User:     &info,
	})
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	result := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	}
	core.JSON(c, result)
}

func (o *OnlineUserController) ForceLogout(ctx context.Context, c *app.RequestContext) {
	var req onlineUserParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ForceLogout(ctx, &v1.ForceLogoutRequest{
		TokenId: req.TokenId,
		User:    &info,
	})
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
