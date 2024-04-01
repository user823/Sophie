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
	UserName string `json:"userName" query:"userName"`
	Ipaddr   string `json:"ipaddr" query:"ipaddr"`
}

// OnlineList godoc
// @Summary 列出在线用户
// @Description 权限：monitor:online:list
// @Param userName formData string false "用户名"
// @Param ipaddr formData string false "ip地址"
// @Accept application/json
// @Produce application/json
// @Router /online/list [GET]
func (o *OnlineUserController) List(ctx context.Context, c *app.RequestContext) {
	var req onlineUserParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListSysUserOnlines(ctx, &v1.ListSysUserOnlinesRequest{
		PageInfo: &req.PageInfo,
		UserName: req.UserName,
		Ipaddr:   req.Ipaddr,
		User:     &info,
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
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

// ForceLogout godoc
// @Summary 列出在线用户
// @Description 权限：monitor:online:list
// @Param tokenId query int false "访问id"
// @Accept application/json
// @Produce application/json
// @Router /online/:tokenId [GET]
func (o *OnlineUserController) ForceLogout(ctx context.Context, c *app.RequestContext) {
	tokenId := c.Param("tokenId")

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ForceLogout(ctx, &v1.ForceLogoutRequest{
		TokenId: tokenId,
		User:    &info,
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
