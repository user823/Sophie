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

type authRoleParam struct {
	UserId  int64   `json:"userId"`
	RoleIds []int64 `json:"roleIds"`
}

type userRequestParam struct {
	api.GetOptions
	v12.SysUser
}

func (u *UserController) List(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListSysUsers(ctx, &v1.ListSysUsersRequest{
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
		core.WriteResponseE(c, rpc.ErrRPC, nil)
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
	loginUser, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, "获取用户登录信息失败，请重试", nil)
		return
	}

	resp, err := rpc.Remoting.GetUserInfo(ctx, loginUser.User.UserId)
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	res := map[string]any{
		"user":        resp.Data,
		"roles":       resp.Roles,
		"permissions": resp.Permissions,
	}
	core.JSON(c, res)
}

func (u *UserController) GetInfoWithName(ctx context.Context, c *app.RequestContext) {
	username := c.Param("username")

	resp, err := rpc.Remoting.GetUserInfoByName(ctx, username)
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	data := map[string]any{
		"sysUser":     resp.Data,
		"permissions": resp.Permissions,
		"roles":       resp.Roles,
	}
	core.OK(c, resp.BaseResp.Msg, data)
}

// 匹配/system/user/:userid
func (u *UserController) GetInfoWithId(ctx context.Context, c *app.RequestContext) {
	userIdStr := c.Param("userId")
	userid, _ := strconv.ParseInt(userIdStr, 10, 64)

	loginUser, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, "获取用户登录信息失败，请重试", nil)
	}

	// 绑定的userid 可以为空
	resp, err := rpc.Remoting.GetUserInfoById(ctx, &v1.GetUserInfoByIdRequest{
		Id: userid,
		User: &v1.LoginUser{
			Roles:       loginUser.Roles,
			Permissions: loginUser.Permissions,
			User:        loginUser.User,
		},
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	result := map[string]any{
		"code":    resp.BaseResp.Code,
		"msg":     resp.BaseResp.Msg,
		"roles":   resp.Roles,
		"posts":   resp.Posts,
		"postIds": resp.PostIds,
		"roleIds": resp.RoleIds,
	}
	core.JSON(c, result)
}

// 匹配/system/user
func (u *UserController) GetInfoWithId2(ctx context.Context, c *app.RequestContext) {
	loginUser, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, "获取用户登录信息失败，请重试", nil)
	}

	resp, err := rpc.Remoting.GetUserInfoById(ctx, &v1.GetUserInfoByIdRequest{
		Id: -1,
		User: &v1.LoginUser{
			Roles:       loginUser.Roles,
			Permissions: loginUser.Permissions,
			User:        loginUser.User,
		},
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	data := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"roles": resp.GetRoles(),
		"posts": resp.GetPosts(),
	}
	core.JSON(c, data)
}

func (u *UserController) Add(ctx context.Context, c *app.RequestContext) {
	var req userRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysUser(ctx, &v1.CreateSysUserRequest{
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
		User:     &info,
	})

	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
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

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysUser(ctx, &v1.UpdateSysUserRequest{
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
		User:     &info,
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) Remove(ctx context.Context, c *app.RequestContext) {
	userIdsStr := c.Param("userIds")

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysUser(ctx, &v1.DeleteSysUserRequest{
		UserIds: strutil.Strs2Int64(userIdsStr),
		User:    &info,
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
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

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ResetPassword(ctx, &v1.ResetPasswordRequest{
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
		User:     &info,
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
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

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ChangeSysUserStatus(ctx, &v1.ChangeSysUserStatus{
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
		User:     &info,
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (u *UserController) AuthRoleWithId(ctx context.Context, c *app.RequestContext) {
	userIdStr := c.Param("userId")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	resp, err := rpc.Remoting.GetAuthRoleById(ctx, userId)
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
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

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.AuthRole(ctx, &v1.AuthRoleRequest{
		UserId:  req.UserId,
		RoleIds: req.RoleIds,
		User:    &info,
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
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

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListDeptsTree(ctx, &v1.ListDeptsTreeRequest{
		DeptInfo: &req,
		User:     &info,
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
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
