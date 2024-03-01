package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type RoleController struct{}

func NewRoleController() *RoleController {
	return &RoleController{}
}

type deleteRoleParam struct {
	RoleIds []int64 `json:"roleIds"`
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
	v1.RoleInfo
	v1.PageInfo
	v1.DateRange
}

func (r *RoleController) List(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListSysRole(ctx, &v1.ListSysRolesRequest{
		PageInfo:  &req.PageInfo,
		DateRange: &req.DateRange,
		RoleInfo:  &req.RoleInfo,
	})
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.JSON(c, map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	})
}

func (r *RoleController) Export(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (r *RoleController) GetInfo(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetSysRoleByid(ctx, req.RoleId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.WriteResponse(c, core.ErrResponse{
		Code:    int(resp.BaseResp.Code),
		Message: resp.BaseResp.Msg,
		Data:    resp.Data,
	})
}

func (r *RoleController) Add(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysRole(ctx, &v1.CreateSysRoleRequest{
		RoleInfo: &req.RoleInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) Edit(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysRole(ctx, &v1.UpdateSysRoleRequest{
		RoleInfo: &req.RoleInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) DataScope(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DataScope(ctx, &v1.DataScopeRequest{
		RoleInfo: &req.RoleInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) ChangeStatus(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ChangeSysRoleStatus(ctx, &v1.ChangeSysRoleStatusRequest{
		RoleInfo: &req.RoleInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deleteRoleParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysRole(ctx, &v1.DeleteSysRoleRequest{
		RoleIds: req.RoleIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) OptionSelect(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.ListRoleOption(ctx)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.WriteResponse(c, core.ErrResponse{
		Code:    int(resp.BaseResp.Code),
		Message: resp.BaseResp.Msg,
		Data:    resp.Rows,
	})
}

func (r *RoleController) AllocatedList(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.AllocatedList(ctx, &v1.AllocatedListRequest{
		PageInfo: &req.PageInfo,
		UserInfo: &req.UserInfo,
	})
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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

func (r *RoleController) UnallocatedList(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UnallocatedList(ctx, &v1.UnallocatedListRequest{
		PageInfo: &req.PageInfo,
		UserInfo: &req.UserInfo,
	})
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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

func (r *RoleController) CancelAuthUser(ctx context.Context, c *app.RequestContext) {
	var req userRoleParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CancelAuthUser(ctx, &v1.CancelAuthUserRequest{
		RoleId: req.RoleId,
		UserId: req.UserId,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) CancelAuthUserAll(ctx context.Context, c *app.RequestContext) {
	var req roleUsersParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CancelAuthUserAll(ctx, &v1.CancelAuthUserAllRequest{
		RoleId:  req.RoleId,
		UserIds: req.UserIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) SelectAuthUserALl(ctx context.Context, c *app.RequestContext) {
	var req roleUsersParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.SelectAuthUserAll(ctx, &v1.SelectAuthUserAllRequest{
		RoleId:  req.RoleId,
		UserIds: req.UserIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (r *RoleController) DeptTree(ctx context.Context, c *app.RequestContext) {
	var req roleRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeptTreeByRoleId(ctx, req.RoleId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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
