package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	v12 "github.com/user823/Sophie/api/domain/gateway/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

type uploadParam struct {
	// 是否禁用上传
	IsUploading bool `json:"isUploading"`
	// 是否更新已存在的用户数据
	UpdateSupport int `json:"updateSupport"`
}

type deleteUserParam struct {
	// 需要删除的用户id
	UserIds []int64 `json:"userIds"`
}

type authRoleParam struct {
	UserId  int64   `json:"userId"`
	RoleIds []int64 `json:"roleIds"`
}

type userRequestParam struct {
	v1.UserInfo
	v1.PageInfo
}

func (u *UserController) List(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListSysUsers(ctx, &v1.ListSysUsersRequest{
		PageInfo: &req.PageInfo,
		UserInfo: &req.UserInfo,
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

func (u *UserController) Export(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (u *UserController) ImportData(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (u *UserController) ImportTemplate(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (u *UserController) GetInfo(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}
	userName := req.UserName
	// 没有传递请求参数的情况下, 则自动获取当前登录用户信息
	if userName == "" {
		data, _ := c.Get(api.LOGIN_INFO_KEY)
		loginUser := data.(v12.LoginUser)
		userName = loginUser.User.Username
	}

	resp, err := rpc.Remoting.GetUserInfoByName(ctx, userName)
	if err = rpc.ParseRpcErr(resp.GetBaseResp(), err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	res := map[string]any{
		"code":        resp.GetBaseResp().Code,
		"msg":         resp.GetBaseResp().Msg,
		"roles":       resp.GetRoles(),
		"permissions": resp.GetPermissions(),
	}

	if resp.Data != nil {
		res["user"] = *v1.UserInfo2SysUser(resp.Data)
	}

	core.JSON(c, res)
}

func (u *UserController) GetInfoWithId(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	// 绑定的userid 可以为空
	resp, err := rpc.Remoting.GetUserInfoById(ctx, req.GetUserId())
	if err = rpc.ParseRpcErr(resp.GetBaseResp(), err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	result := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"roles": resp.Roles,
		"posts": resp.Posts,
	}

	if req.GetUserId() != 0 {
		result["postIds"] = resp.PostIds
		result["roleIds"] = resp.RoleIds
	}
	core.JSON(c, result)
}

func (u *UserController) Add(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysUser(ctx, &v1.CreateSysUserRequest{
		UserInfo: &req.UserInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) Edit(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysUser(ctx, &v1.UpdateSysUserRequest{
		UserInfo: &req.UserInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deleteUserParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysUser(ctx, &v1.DeleteSysUserRequest{
		UserIds: req.UserIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) ResetPwd(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ResetPassword(ctx, &v1.ResetPasswordRequest{
		UserInfo: &req.UserInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) ChangeStatus(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ChangeSysUserStatus(ctx, &v1.ChangeSysUserStatus{
		UserInfo: &req.UserInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) AuthRoleWithId(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetAuthRoleById(ctx, req.UserId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	result := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"user":  resp.User,
		"roles": resp.Roles,
	}
	core.JSON(c, result)
}

func (u *UserController) InsertAuthRole(ctx context.Context, c *app.RequestContext) {
	var req authRoleParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.AuthRole(ctx, &v1.AuthRoleRequest{
		UserId:  req.UserId,
		RoleIds: req.RoleIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) DeptTree(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListDeptsTree(ctx, &v1.ListDeptsTreeRequest{
		DeptInfo: &req,
	})
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
