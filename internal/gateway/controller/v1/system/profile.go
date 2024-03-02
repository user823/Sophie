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

type ProfileController struct{}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}

type updatePwsParam struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (p *ProfileController) Profile(ctx context.Context, c *app.RequestContext) {
	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.Profile(ctx, &v1.ProfileRequest{
		User: v1.LoginUserTrans(&info),
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
	}
	core.JSON(c, result)
}

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
		User:     v1.LoginUserTrans(&info),
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
		User:         v1.LoginUserTrans(&info),
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

func (p *ProfileController) Avatar(ctx context.Context, c *app.RequestContext) {
	// TODO
}
