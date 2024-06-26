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

// ConfigList godoc
// @Summary 配置列表
// @Description 根据条件查询配置列表
// @Description 权限：system:config:list
// @Param configName formData string false "参数名称"
// @Param configKey formData string false "参数键名"
// @Param configType formData string false "系统内置"
// @Accept application/json
// @Produce application/json
// @Router /system/config/list [GET]
func (f *ConfigController) List(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
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
		LoginUser: &info,
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

// ConfigExport godoc
// @Summary 导出配置列表
// @Description 根据条件导出配置列表
// @Description 权限：system:config:export
// @Param configName formData string false "参数名称"
// @Param configKey formData string false "参数键名"
// @Param configType formData string false "系统内置"
// @Accept application/json
// @Produce application/json
// @Router /system/config/export [POST]
func (f *ConfigController) Export(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ExportConfig(ctx, &v1.ExportConfigRequest{
		LoginUser:  &info,
		ConfigInfo: v1.SysConfig2ConfigInfo(&req.SysConfig),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	utils.ExportExcel(c, resp.SheetName, resp.List)
}

// ConfigInfo godoc
// @Summary 获取配置详情
// @Param configId query int true "参数id"
// @Accept application/json
// @Produce application/json
// @Router /system/config/:configId [GET]
func (f *ConfigController) GetInfo(ctx context.Context, c *app.RequestContext) {
	configIdStr := c.Param("configId")
	configId, _ := strconv.ParseInt(configIdStr, 10, 64)

	resp, err := rpc.Remoting.GetConfigById(ctx, configId)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

// ConfigKey godoc
// @Summary 获取配置键
// @Param configKey query string true "参数键名"
// @Accept application/json
// @Produce application/json
// @Router /system/config/:configKey [GET]
func (f *ConfigController) GetConfigKey(ctx context.Context, c *app.RequestContext) {
	configKey := c.Param("configKey")

	resp, err := rpc.Remoting.GetConfigByKey(ctx, configKey)
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

// ConfigAdd godoc
// @Summary 添加配置
// @Description 权限：system:config:add
// @Param configName formData string true "参数名"
// @Param configKey formData string true "参数键"
// @Param configValue formData string true "参数值"
// @Param configType formData string true "系统内置"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/config [POST]
func (f *ConfigController) Add(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateConfig(ctx, &v1.CreateConfigRequest{
		ConfigInfo: v1.SysConfig2ConfigInfo(&req.SysConfig),
		User:       &info,
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

// ConfigEdit godoc
// @Summary 添加配置
// @Description 权限：system:config:edit
// @Param configName formData string true "参数名"
// @Param configKey formData string true "参数键"
// @Param configValue formData string true "参数值"
// @Param configType formData string true "系统内置"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/config/:configKey [PUT]
func (f *ConfigController) Edit(ctx context.Context, c *app.RequestContext) {
	var req configRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateConfig(ctx, &v1.UpdateConfigReqeust{
		ConfigInfo: v1.SysConfig2ConfigInfo(&req.SysConfig),
		User:       &info,
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

// ConfigRemove godoc
// @Summary 移除配置
// @Description 权限: system:config:remove
// @Param configIds query string true "配置id"
// @Accept application/json
// @Produce application/json
// @Router /system/config/:configIds [DELETE]
func (f *ConfigController) Remove(ctx context.Context, c *app.RequestContext) {
	configIdsStr := c.Param("configIds")
	configIds := strutil.Strs2Int64(configIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteConfig(ctx, &v1.DeleteConfigReqeust{
		ConfigIds: configIds,
		User:      &info,
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

// RefreshCache godoc
// @Summary 刷新缓存
// @Description 权限: system:config:remove
// @Accept application/json
// @Produce application/json
// @Router /system/config/refreshCache [DELETE]
func (f *ConfigController) RefreshCache(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.RefreshConfig(ctx)
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
