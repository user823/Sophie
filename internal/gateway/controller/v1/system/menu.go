package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type MenuController struct{}

func NewMenuController() *MenuController {
	return &MenuController{}
}

type roleMenuTreeParam struct {
	RoleId int64 `json:"roleId"`
}

type deleteMenuParam struct {
	MenuId int64 `json:"MenuId"`
}

func (m *MenuController) List(ctx context.Context, c *app.RequestContext) {
	var req v1.MenuInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListSysMenus(ctx, &v1.ListSysMenusRequest{
		MenuInfo: &req,
	})
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (m *MenuController) GetInfo(ctx context.Context, c *app.RequestContext) {
	var req v1.MenuInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetSysMenuById(ctx, req.MenuId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (m *MenuController) TreeSelect(ctx context.Context, c *app.RequestContext) {
	var req v1.MenuInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListTreeMenu(ctx, &v1.ListTreeMenuRequest{
		MenuInfo: &req,
	})
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (m *MenuController) RoleMenuTreeselect(ctx context.Context, c *app.RequestContext) {
	var req roleMenuTreeParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListTreeMenuByRoleid(ctx, req.RoleId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	result := map[string]any{
		"code":        resp.BaseResp.Code,
		"msg":         resp.BaseResp.Msg,
		"checkedKeys": resp.CheckedKeys,
		"menus":       resp.Menus,
	}
	core.JSON(c, result)
}

func (m *MenuController) Add(ctx context.Context, c *app.RequestContext) {
	var req v1.MenuInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateMenu(ctx, &v1.CreateMenuRequest{
		MenuInfo: &req,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (m *MenuController) Edit(ctx context.Context, c *app.RequestContext) {
	var req v1.MenuInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateMenu(ctx, &v1.UpdateMenuRequest{
		MenuInfo: &req,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (m *MenuController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deleteMenuParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteMenu(ctx, &v1.DeleteMenuRequest{
		MenuId: req.MenuId,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (m *MenuController) GetRouters(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.GetRouters(ctx)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}
