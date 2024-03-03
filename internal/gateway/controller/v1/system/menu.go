package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"strconv"
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
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListSysMenus(ctx, &v1.ListSysMenusRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
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

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (m *MenuController) GetInfo(ctx context.Context, c *app.RequestContext) {
	menuIdStr := c.Param("menuId")
	menuId, _ := strconv.ParseInt(menuIdStr, 10, 64)

	resp, err := rpc.Remoting.GetSysMenuById(ctx, menuId)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (m *MenuController) TreeSelect(ctx context.Context, c *app.RequestContext) {
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListTreeMenu(ctx, &v1.ListTreeMenuRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
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

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (m *MenuController) RoleMenuTreeselect(ctx context.Context, c *app.RequestContext) {
	roleIdStr := c.Param("roleId")
	roleId, _ := strconv.ParseInt(roleIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	resp, err := rpc.Remoting.ListTreeMenuByRoleid(ctx, &v1.ListTreeMenuByRoleidRequest{
		Id:   roleId,
		User: &info,
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
		"code":        resp.BaseResp.Code,
		"msg":         resp.BaseResp.Msg,
		"checkedKeys": resp.CheckedKeys,
		"menus":       resp.Menus,
	}
	core.JSON(c, result)
}

func (m *MenuController) Add(ctx context.Context, c *app.RequestContext) {
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.CreateMenu(ctx, &v1.CreateMenuRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
		User:     &info,
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

func (m *MenuController) Edit(ctx context.Context, c *app.RequestContext) {
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.UpdateMenu(ctx, &v1.UpdateMenuRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
		User:     &info,
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

func (m *MenuController) Remove(ctx context.Context, c *app.RequestContext) {
	menuIdStr := c.Param("menuId")
	menuId, _ := strconv.ParseInt(menuIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.DeleteMenu(ctx, &v1.DeleteMenuRequest{
		MenuId: menuId,
		User:   &info,
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

func (m *MenuController) GetRouters(ctx context.Context, c *app.RequestContext) {
	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	resp, err := rpc.Remoting.GetRouters(ctx, &v1.GetRoutersRequest{
		User: &info,
	})
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}
