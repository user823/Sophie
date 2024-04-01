package system

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	v12 "github.com/user823/Sophie/api/thrift/file/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/log"
	utils2 "github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"io"
)

type ProfileController struct{}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}

type updatePwsParam struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Profile godoc
// @Summary 用户个人详情
// @Description 查询登录用户详情信息
// @Accept application/json
// @Produce application/json
// @Router /system/profile [GET]
func (p *ProfileController) Profile(ctx context.Context, c *app.RequestContext) {
	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.Profile(ctx, &v1.ProfileRequest{
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
		"code":      resp.BaseResp.Code,
		"msg":       resp.BaseResp.Msg,
		"roleGroup": resp.RoleGroup,
		"postGroup": resp.PostGroup,
		"data":      resp.UserInfo,
	}
	core.JSON(c, result)
}

// UpdateProfile godoc
// @Summary 更新用户个人详情
// @Description 更新登录用户详情信息
// @Param nickName formData string true "用户昵称"
// @Param phonenumber formData string true "手机号码"
// @Param email formData string true "用户邮箱"
// @Param sex formData string false "性别"
// @Accept application/json
// @Produce application/json
// @Router /system/profile [PUT]
func (p *ProfileController) UpdateProfile(ctx context.Context, c *app.RequestContext) {
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

	resp, err := rpc.Remoting.UpdateProfile(ctx, &v1.UpdateProfileRequest{
		UserInfo: v1.SysUser2UserInfo(&req.SysUser),
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

// UpdatePwd godoc
// @Summary 更新密码
// @Description 更新登录用户密码
// @Param oldPassword formData string true "旧密码"
// @Param newPassword formData string true "新密码"
// @Accept application/json
// @Produce application/json
// @Router /system/updatePwd [PUT]
func (p *ProfileController) UpdatePwd(ctx context.Context, c *app.RequestContext) {
	var req updatePwsParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.UpdatePassword(ctx, &v1.UpdatePasswordRequest{
		OldPassword:  req.OldPassword,
		NewPassword_: req.NewPassword,
		User:         &info,
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

// UpdatePwd godoc
// @Summary 更新密码
// @Description 更新登录用户密码
// @Param avatarfile formData file true "用户头像"
// @Accept application/json
// @Produce application/json
// @Router /system/avatar [POST]
func (p *ProfileController) Avatar(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("avatarfile")
	if err != nil {
		core.Fail(c, "上传图片失败，请联系管理员", nil)
		return
	}
	filename := string(c.FormValue("fileName"))
	if filename == "" {
		filename = "blob"
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	// 检查文件后缀
	ext := utils2.GetExtension(file)
	if !strutil.ContainsAny(ext, strutil.IMAGE_FILE...) {
		core.Fail(c, fmt.Sprintf("文件格式不正确，仅支持%v", strutil.IMAGE_FILE), nil)
		return
	}

	freader, err := file.Open()
	if err != nil {
		core.Fail(c, "文件打开失败，请重新上传", nil)
		return
	}
	data, err := io.ReadAll(freader)
	if err != nil {
		core.Fail(c, "文件读取失败，请重新上传", nil)
		return
	}

	log.Infof("文件名: %s", filename+ext)
	resp, err := rpc.Remoting.Upload(ctx, &v12.UploadRequest{
		Path: filename + ext,
		Data: data,
	})

	if err != nil {
		log.Error("error: %s", err.Error())
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	resp1, err := rpc.Remoting.UpdateUserAvatar(ctx, &v1.UpdateUserAvatarRequest{
		User:   &info,
		Avatar: resp.File.Url,
	})
	if err != nil {
		log.Error("error: %s", err.Error())
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp1.Code != code.SUCCESS {
		core.Fail(c, resp1.Msg, nil)
		return
	}

	res := map[string]any{
		"code":   200,
		"msg":    "操作成功",
		"imgUrl": resp.File.Url,
	}
	core.JSON(c, res)
}
