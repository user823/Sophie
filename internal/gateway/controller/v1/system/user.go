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

type authRoleParam struct {
	UserId  int64   `json:"userId"`
	RoleIds []int64 `json:"roleIds"`
}

type userRequestParam struct {
	api.GetOptions
	v12.SysUser
}

// UserList godoc
// @Summary 获取用户列表
// @Description 根据条件查询系统用户列表
// @Description 权限：system:user:list
// @Param userName formData string false "用户名"
// @Param phonenumber formData string false "手机号码"
// @Param status formData string false "状态"
// @Param createdTime formData string false "创建时间"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/user/list [GET]
func (u *UserController) List(ctx context.Context, c *app.RequestContext) {
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

// UserExport godoc
// @Summary 导出用户列表
// @Description 根据条件导出系统用户列表
// @Description 权限：system:user:export
// @Param userName formData string false "用户名"
// @Param phonenumber formData string false "手机号码"
// @Param status formData string false "状态"
// @Param createdTime formData string false "创建时间"
// @Accept application/json
// @Produce application/json
// @Router /system/user/export [POST]
func (u *UserController) Export(ctx context.Context, c *app.RequestContext) {
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

	resp, err := rpc.Remoting.ExportSysUser(ctx, &v1.ExportSysUserRequest{
		User:     &info,
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
	})
	if err != nil {
		core.WriteResponseE(c, rpc.ErrRPC, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	utils.ExportExcel(c, resp.SheetName, v1.MUserInfo2SysUser(resp.List))
}

// UserInfo godoc
// @Summary 用户详情
// @Description 用户登录后获取当前用户详情信息
// @Produce application/json
// @Router /system/user/getInfo [get]
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

// UserInfoWithName godoc
// @Summary 用户详情
// @Description 获取目标用户的详情信息
// @Param username query string true "用户名"
// @Accept application/json
// @Produce application/json
// @Router /system/user/info/:username [GET]
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

// UserInfoWithId godoc
// @Summary 用户详情
// @Description 获取目标用户的详情信息
// @Description 权限：system:user:query
// @Param userId query string true "用户id"
// @Accept application/json
// @Produce application/json
// @Router /system/user/:userid [GET]
func (u *UserController) GetInfoWithId(ctx context.Context, c *app.RequestContext) {
	userIdStr := c.Param("userId")
	userid, _ := strconv.ParseInt(userIdStr, 10, 64)

	loginUser, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, "获取用户登录信息失败，请重试", nil)
	}

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
		"data":    resp.Data,
	}
	core.JSON(c, result)
}

// 匹配/system/user

// UserInfoWithId godoc
// @Summary 用户详情
// @Description 获取目标用户（登录用户）的详情信息
// @Description 权限：system:user:query
// @Accept application/json
// @Produce application/json
// @Router /system/user/ [GET]
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

// UserAdd godoc
// @Summary 创建用户
// @Description 创建用户
// @Description 权限：system:user:add
// @Param deptId formData int true "部门id"
// @Param userName formData string true "用户名"
// @Param nickName formData string true "用户昵称"
// @Param email formData string false "用户邮箱"
// @Param phonenumber formData string false "手机号"
// @Param sex formData string false "性别"
// @Param password formData string true "密码"
// @Param status formData string true "用户状态"
// @Param roleIds formData string false "角色id"
// @Param postIds formData string false "岗位id"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/user [POST]
func (u *UserController) Add(ctx context.Context, c *app.RequestContext) {
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

// UserEdit godoc
// @Summary 修改用户
// @Description 修改用户
// @Description 权限：system:user:edit
// @Param deptId formData int true "部门id"
// @Param nickName formData string true "用户昵称"
// @Param email formData string false "用户邮箱"
// @Param phonenumber formData string false "手机号"
// @Param sex formData string false "性别"
// @Param status formData string true "用户状态"
// @Param roleIds formData string false "角色id"
// @Param postIds formData string false "岗位id"
// @Accept application/json
// @Produce application/json
// @Router /system/user [PUT]
func (u *UserController) Edit(ctx context.Context, c *app.RequestContext) {
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

// UserDelete godoc
// @Summary 删除用户
// @Description 删除用户
// @Description 权限：system:user:remove
// @Param userIds query string true "需要删除的用户id"
// @Accept application/json
// @Produce application/json
// @Router /system/user/:userIds [DELETE]
func (u *UserController) Remove(ctx context.Context, c *app.RequestContext) {
	userIdsStr := c.Param("userIds")

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
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

// ResetPassword godoc
// @Summary 重置密码
// @Description 重置密码
// @Description 权限：system:user:edit
// @Param userId formData int true "用户id"
// @Param password formData string true "新密码"
// @Accept application/json
// @Produce application/json
// @Router /system/user/resetPwd [PUT]
func (u *UserController) ResetPwd(ctx context.Context, c *app.RequestContext) {
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

// ChangeStatus godoc
// @Summary 修改用户状态
// @Description 修改目标用户状态
// @Description 权限：system:user:edit
// @Param userId formData int true "用户id"
// @Param status formData string true "状态"
// @Accept application/json
// @Produce application/json
// @Router /system/user/changeStatus [PUT]
func (u *UserController) ChangeStatus(ctx context.Context, c *app.RequestContext) {
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

// AuthRoleWithId godoc
// @Summary 查询角色
// @Description 查询目标用户的角色
// @Description 权限：system:user:query
// @Param userId query int true "用户id"
// @Accept application/json
// @Produce application/json
// @Router /system/user/authRole/:userId [GET]
func (u *UserController) AuthRoleWithId(ctx context.Context, c *app.RequestContext) {
	userIdStr := c.Param("userId")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.GetAuthRoleById(ctx, &v1.GetAuthRoleByIdRequest{
		Id:   userId,
		User: &info,
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
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"user":  resp.User,
		"roles": resp.Roles,
	}
	core.JSON(c, result)
}

// InsertAuthRole godoc
// @Summary 授权角色
// @Description 为目标用户分配角色
// @Description 权限：system:user:edit
// @Param userId formData int true "用户id"
// @Param roleIds formData []int true "授权角色的id"
// @Accept application/json
// @Produce application/json
// @Router /system/user/authRole [PUT]
func (u *UserController) InsertAuthRole(ctx context.Context, c *app.RequestContext) {
	var req authRoleParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
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

// DeptTree godoc
// @Summary 查询部门树
// @Description 根据当前登录用户，查询他的部门信息
// @Description 权限：system:user:list
// @Accept application/json
// @Produce application/json
// @Router /system/user/deptTree [GET]
func (u *UserController) DeptTree(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
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
