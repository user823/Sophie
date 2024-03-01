package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type ConfigController struct{}

func NewConfigController() *ConfigController {
	return &ConfigController{}
}

type configRequestParam struct {
	v1.ConfigInfo
	v1.PageInfo
	v1.DateRange
}

type deleteConfigParam struct {
	ConfigIds []int64 `json:"configIds"`
}

func (f *ConfigController) List(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListConfigs(ctx, &v1.ListConfigsRequest{
		ConfigInfo: &req.ConfigInfo,
		PageInfo:   &req.PageInfo,
		DateRange:  &req.DateRange,
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

func (f *ConfigController) Export(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (f *ConfigController) GetInfo(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetConfigById(ctx, req.ConfigId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (f *ConfigController) GetConfigKey(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetConfigByKey(ctx, req.ConfigKey)
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (f *ConfigController) Add(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateConfig(ctx, &v1.CreateConfigRequest{
		ConfigInfo: &req.ConfigInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (f *ConfigController) Edit(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateConfig(ctx, &v1.UpdateConfigReqeust{
		ConfigInfo: &req.ConfigInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (f *ConfigController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deleteConfigParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteConfig(ctx, &v1.DeleteConfigReqeust{
		ConfigIds: req.ConfigIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (f *ConfigController) RefreshCache(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.RefreshConfig(ctx)
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
