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

// MenuList godoc
// @Summary 列出菜单列表
// @Description 根据条件查询菜单列表
// @Description 权限：system:menu:list
// @Param menuName formData string false "菜单名称"
// @Param status formData string false "状态"
// @Accept application/json
// @Produce application/json
// @Router /system/menu/list [GET]
func (m *MenuController) List(ctx context.Context, c *app.RequestContext) {
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListSysMenus(ctx, &v1.ListSysMenusRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
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

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

// GetInfo godoc
// @Summary 菜单详情
// @Description 根据目标菜单详情信息
// @Description 权限：system:menu:query
// @Param menuId query int true "菜单id"
// @Accept application/json
// @Produce application/json
// @Router /system/menu/:menuId [GET]
func (m *MenuController) GetInfo(ctx context.Context, c *app.RequestContext) {
	menuIdStr := c.Param("menuId")
	menuId, _ := strconv.ParseInt(menuIdStr, 10, 64)

	resp, err := rpc.Remoting.GetSysMenuById(ctx, menuId)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

// TreeSelect godoc
// @Summary 菜单树
// @Description 查询目标菜单树
// @Param menuId formData int true "菜单id"
// @Param menuName formData string false "菜单名"
// @Accept application/json
// @Produce application/json
// @Router /system/menu/treeselect [GET]
func (m *MenuController) TreeSelect(ctx context.Context, c *app.RequestContext) {
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListTreeMenu(ctx, &v1.ListTreeMenuRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
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

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

// RoleMenuTreeselect godoc
// @Summary 角色菜单树
// @Description 根据角色查询菜单树
// @Param roleId formData int true "角色id"
// @Accept application/json
// @Produce application/json
// @Router /system/menu/roleMenuTreeselect/:roleId [GET]
func (m *MenuController) RoleMenuTreeselect(ctx context.Context, c *app.RequestContext) {
	roleIdStr := c.Param("roleId")
	roleId, _ := strconv.ParseInt(roleIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	resp, err := rpc.Remoting.ListTreeMenuByRoleid(ctx, &v1.ListTreeMenuByRoleidRequest{
		Id:   roleId,
		User: &info,
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
		"code":        resp.BaseResp.Code,
		"msg":         resp.BaseResp.Msg,
		"checkedKeys": resp.CheckedKeys,
		"menus":       resp.Menus,
	}
	core.JSON(c, result)
}

// MenuAdd godoc
// @Summary 添加菜单
// @Param menuName formData string true "菜单名称"
// @Param sort formData string true "显示排序"
// @Param path formData string true "路由地址"
// @Param parentName formData string false "上级菜单"
// @Param menuType formData string false "菜单类型"
// @Param icon formData string false "菜单图标"
// @Param status formData string true "状态"
// @Param isFrame formData string true "是否外链"
// @Accept application/json
// @Produce application/json
// @Router /system/menu/add [POST]
func (m *MenuController) Add(ctx context.Context, c *app.RequestContext) {
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateMenu(ctx, &v1.CreateMenuRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
		User:     &info,
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

// MenuEdit godoc
// @Summary 添加菜单
// @Description 修改菜单
// @Description 权限: system:menu:add
// @Param menuName formData string true "菜单名称"
// @Param sort formData string true "显示排序"
// @Param path formData string true "路由地址"
// @Param parentName formData string false "上级菜单"
// @Param menuType formData string false "菜单类型"
// @Param icon formData string false "菜单图标"
// @Param status formData string true "状态"
// @Param isFrame formData string true "是否外链"
// @Accept application/json
// @Produce application/json
// @Router /system/menu [PUT]
func (m *MenuController) Edit(ctx context.Context, c *app.RequestContext) {
	var req v12.SysMenu
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateMenu(ctx, &v1.UpdateMenuRequest{
		MenuInfo: v1.SysMenu2MenuInfo(&req),
		User:     &info,
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

// MenuRemove godoc
// @Summary 移除菜单
// @Description 移除目标菜单
// @Description 权限: system:menu:remove
// @Param menuId formData int true "菜单id"
// @Accept application/json
// @Produce application/json
// @Router /system/menu/:menuId [DELETE]
func (m *MenuController) Remove(ctx context.Context, c *app.RequestContext) {
	menuIdStr := c.Param("menuId")
	menuId, _ := strconv.ParseInt(menuIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteMenu(ctx, &v1.DeleteMenuRequest{
		MenuId: menuId,
		User:   &info,
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

// GetRouters godoc
// @Summary 获取路由信息
// @Description 用户访问时获取菜单路由信息
// @Accept application/json
// @Produce application/json
// @Router /system/menu/getRouters [GET]
func (m *MenuController) GetRouters(ctx context.Context, c *app.RequestContext) {
	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	resp, err := rpc.Remoting.GetRouters(ctx, &v1.GetRoutersRequest{
		User: &info,
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}
