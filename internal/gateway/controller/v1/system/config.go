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

type ConfigController struct{}

func NewConfigController() *ConfigController {
	return &ConfigController{}
}

type configRequestParam struct {
	v12.SysConfig
	api.GetOptions
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

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListConfigs(ctx, &v1.ListConfigsRequest{
		ConfigInfo: v1.SysConfig2ConfigInfo(&req.SysConfig),
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
		LoginUser: v1.LoginUserTrans(&info),
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
	configIdStr := c.Param("configId")
	configId, _ := strconv.ParseInt(configIdStr, 10, 64)

	resp, err := rpc.Remoting.GetConfigById(ctx, configId)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (f *ConfigController) GetConfigKey(ctx context.Context, c *app.RequestContext) {
	configKey := c.Param("configKey")

	resp, err := rpc.Remoting.GetConfigByKey(ctx, configKey)
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

func (f *ConfigController) Add(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.CreateConfig(ctx, &v1.CreateConfigRequest{
		ConfigInfo: v1.SysConfig2ConfigInfo(&req.SysConfig),
		User:       v1.LoginUserTrans(&info),
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

func (f *ConfigController) Edit(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.UpdateConfig(ctx, &v1.UpdateConfigReqeust{
		ConfigInfo: v1.SysConfig2ConfigInfo(&req.SysConfig),
		User:       v1.LoginUserTrans(&info),
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

func (f *ConfigController) Remove(ctx context.Context, c *app.RequestContext) {
	configIdsStr := c.Param("configIds")
	configIds := strutil.Strs2Int64(configIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.DeleteConfig(ctx, &v1.DeleteConfigReqeust{
		ConfigIds: configIds,
		User:      v1.LoginUserTrans(&info),
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

func (f *ConfigController) RefreshCache(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.RefreshConfig(ctx)
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
