package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type LogininfoController struct{}

func NewLogininfoController() *LogininfoController {
	return &LogininfoController{}
}

type logininforRequestParam struct {
	v1.PageInfo
	v1.Logininfo
	v1.DateRange
}

type deleteLogininforParam struct {
	InfoIds []int64 `json:"infoIds"`
}

type unlockParam struct {
	Username string `json:"userName"`
}

func (l *LogininfoController) List(ctx context.Context, c *app.RequestContext) {
	var req logininforRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListSysLogininfos(ctx, &v1.ListSysLogininfosRequest{
		PageInfo:  &req.PageInfo,
		LoginInfo: &req.Logininfo,
		DateRange: &req.DateRange,
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

func (l *LogininfoController) Export(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (l *LogininfoController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deleteLogininforParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.RemoveSysLogininfosById(ctx, &v1.RemoveSysLogininfosByIdRequest{
		InfoIds: req.InfoIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (l *LogininfoController) Clean(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.LogininfoClean(ctx)
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (l *LogininfoController) Unlock(ctx context.Context, c *app.RequestContext) {
	var req unlockParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UnlockByUserName(ctx, req.Username)
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
