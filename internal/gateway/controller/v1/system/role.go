package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strconv"
)

type RoleController struct{}

func NewRoleController() *RoleController {
	return &RoleController{}
}

type userRoleParam struct {
	UserId int64 `json:"userId"`
	RoleId int64 `json:"roleId"`
}

type roleUsersParam struct {
	RoleId  int64   `json:"roleId"`
	UserIds []int64 `json:"userIds"`
}

type roleRequestParam struct {
	v12.SysRole
	api.GetOptions
}

// RoleList godoc
// @Summary 获取角色
// @Description 根据条件查询系统角色列表
// @Description 权限：system:role:list
// @Param roleName formData string false "角色名"
// @Param roleKey formData string false "权限字符"
// @Param status formData string false "状态"
// @Param createdTime formData string false "创建时间"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/role/list [GET]
func (r *RoleController) List(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListSysRole(ctx, &v1.ListSysRolesRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		DateRange: &v1.DateRange{
			BeginTime: req.BeginTime,
			EndTime:   req.EndTime,
		},
		RoleInfo: v1.SysRole2RoleInfo(&req.SysRole),
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

	core.JSON(c, map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	})
}

// RoleExport godoc
// @Summary 导出角色
// @Description 根据条件导出系统角色列表
// @Description 权限：system:role:export
// @Param roleName formData string false "角色名"
// @Param roleKey formData string false "权限字符"
// @Param status formData string false "状态"
// @Param createdTime formData string false "创建时间"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/role/export [GET]
func (r *RoleController) Export(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ExportSysRole(ctx, &v1.ExportSysRoleRequest{
		User:     &info,
		RoleInfo: v1.SysRole2RoleInfo(&req.SysRole),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	utils.ExportExcel(c, resp.SheetName, v1.MRoleInfo2SysRole(resp.List))
}

// RoleInfo godoc
// @Summary 角色详情
// @Description 查询目标角色详情信息
// @Description 权限：system:role:query
// @Param roleId query int false "角色id"
// @Accept application/json
// @Produce application/json
// @Router /system/role/:roleId [GET]
func (r *RoleController) GetInfo(ctx context.Context, c *app.RequestContext) {
	roleId := c.Param("roleId")
	roleid, _ := strconv.ParseInt(roleId, 10, 64)

	resp, err := rpc.Remoting.GetSysRoleByid(ctx, roleid)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.WriteResponse(c, core.ErrResponse{
		Code:    int(resp.BaseResp.Code),
		Message: resp.BaseResp.Msg,
		Data:    resp.Data,
	})
}

// RoleAdd godoc
// @Summary 创建角色
// @Description 创建角色
// @Description 权限：system:role:add
// @Param roleName formData string true "角色名称"
// @Param roleKey formData string true "权限字符"
// @Param roleSort formData string true "角色顺序"
// @Param status formData string true "状态"
// @Param menuIds formData []int false "菜单权限"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/role/:roleId [POST]
func (r *RoleController) Add(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysRole(ctx, &v1.CreateSysRoleRequest{
		RoleInfo: v1.SysRole2RoleInfo(&req.SysRole),
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

// RoleEdit godoc
// @Summary 修改角色
// @Description 修改角色
// @Description 权限：system:role:edit
// @Param roleName formData string true "角色名称"
// @Param roleKey formData string true "权限字符"
// @Param roleSort formData string true "角色顺序"
// @Param status formData string true "状态"
// @Param menuIds formData []int false "菜单权限"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/role/:roleId [PUT]
func (r *RoleController) Edit(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysRole(ctx, &v1.UpdateSysRoleRequest{
		RoleInfo: v1.SysRole2RoleInfo(&req.SysRole),
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

// DataScope godoc
// @Summary 数据范围
// @Description 限定本角色能操作的数据范围，比如只能操作本人数据、本部门数据等
// @Description 权限：system:role:edit
// @Param dataScope formData string false "数据范围"
// @Accept application/json
// @Produce application/json
// @Router /system/role/dataScope [PUT]
func (r *RoleController) DataScope(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DataScope(ctx, &v1.DataScopeRequest{
		RoleInfo: v1.SysRole2RoleInfo(&req.SysRole),
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

// ChangeStatus godoc
// @Summary 修改状态
// @Description 修改角色的状态
// @Description 权限：system:role:edit
// @Param roleId formData int falst "角色id"
// @Param status formData string true "角色状态"
// @Accept application/json
// @Produce application/json
// @Router /system/role/changeStatus [PUT]
func (r *RoleController) ChangeStatus(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ChangeSysRoleStatus(ctx, &v1.ChangeSysRoleStatusRequest{
		RoleInfo: v1.SysRole2RoleInfo(&req.SysRole),
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

// RoleDelete godoc
// @Summary 删除角色
// @Description 删除目标角色
// @Description 权限：system:role:edit
// @Param roleIds formData string true "删除角色id"
// @Accept application/json
// @Produce application/json
// @Router /system/role/:roleIds [DELETE]
func (r *RoleController) Remove(ctx context.Context, c *app.RequestContext) {
	roleIdsStr := c.Param("roleIds")
	roleIds := strutil.Strs2Int64(roleIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysRole(ctx, &v1.DeleteSysRoleRequest{
		RoleIds: roleIds,
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

// OptionSelect godoc
// @Summary 可选角色
// @Description 列出所有角色
// @Description 权限：system:role:query
// @Accept application/json
// @Produce application/json
// @Router /system/role/optionselect [GET]
func (r *RoleController) OptionSelect(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.ListRoleOption(ctx)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.WriteResponse(c, core.ErrResponse{
		Code:    int(resp.BaseResp.Code),
		Message: resp.BaseResp.Msg,
		Data:    resp.Rows,
	})
}

// AllocatedList godoc
// @Summary 查询用户角色
// @Description 查询目标用户已经分配的角色
// @Description 权限：system:role:list
// @Param userId formData int true "用户id"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/role/authUser/allocatedList [GET]
func (r *RoleController) AllocatedList(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.AllocatedList(ctx, &v1.AllocatedListRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
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

// AllocatedList godoc
// @Summary 查询用户角色
// @Description 查询目标用户还没有分配的角色
// @Description 权限：system:role:list
// @Param userId formData int true "用户id"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/role/authUser/unallocatedList [GET]
func (r *RoleController) UnallocatedList(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UnallocatedList(ctx, &v1.UnallocatedListRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
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

// CancelAuthUser godoc
// @Summary 删除用户角色
// @Description 删除目标用户的角色
// @Description 权限：system:role:edit
// @Param userId formData int true "用户id"
// @Param roleId formData int true "角色id"
// @Accept application/json
// @Produce application/json
// @Router /system/role/authUser/cancel [PUT]
func (r *RoleController) CancelAuthUser(ctx context.Context, c *app.RequestContext) {
	var req userRoleParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CancelAuthUser(ctx, &v1.CancelAuthUserRequest{
		RoleId: req.RoleId,
		UserId: req.UserId,
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

// CancelAuthUserALL godoc
// @Summary 删除用户角色
// @Description 删除多个用户的指定角色
// @Description 权限：system:role:edit
// @Param userIds formData []int true "用户ids"
// @Param roleId formData int true "角色id"
// @Accept application/json
// @Produce application/json
// @Router /system/role/authUser/cancelAll [PUT]
func (r *RoleController) CancelAuthUserAll(ctx context.Context, c *app.RequestContext) {
	var req roleUsersParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CancelAuthUserAll(ctx, &v1.CancelAuthUserAllRequest{
		RoleId:  req.RoleId,
		UserIds: req.UserIds,
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

// SelectAuthUserALl godoc
// @Summary 授权角色
// @Description 为多个用户授权指定角色
// @Description 权限：system:role:edit
// @Param userIds formData []int true "用户ids"
// @Param roleId formData int true "角色id"
// @Accept application/json
// @Produce application/json
// @Router /system/role/authUser/selectAll [PUT]
func (r *RoleController) SelectAuthUserALl(ctx context.Context, c *app.RequestContext) {
	var req roleUsersParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.SelectAuthUserAll(ctx, &v1.SelectAuthUserAllRequest{
		RoleId:  req.RoleId,
		UserIds: req.UserIds,
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

// DeptTree godoc
// @Summary 角色部门树
// @Description 查询指定角色的部门树
// @Description 权限：system:role:query
// @Param roleId query int true "角色id"
// @Accept application/json
// @Produce application/json
// @Router /system/role/deptTree/:roleId [GET]
func (r *RoleController) DeptTree(ctx context.Context, c *app.RequestContext) {
	roleIdStr := c.Param("roleId")
	roleId, _ := strconv.ParseInt(roleIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeptTreeByRoleId(ctx, &v1.DeptTreeByRoleIdRequest{
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
		"depts":       resp.Depts,
	}
	core.JSON(c, result)
}
