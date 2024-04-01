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
)

type OpelogController struct{}

func NewOperlogController() *OpelogController {
	return &OpelogController{}
}

type operLogParam struct {
	v12.SysOperLog
	api.GetOptions
}

type deleteOperLogParam struct {
	OperIds []int64 `json:"operIds"`
}

// OperlogList godoc
// @Summary 列出操作列表
// @Description 根据条件查询操作日志列表
// @Description 权限：system:operlog:list
// @Param operIp formData string false "操作ip"
// @Param title formData string false "系统模块"
// @Param operName formData string false "操作人员"
// @param businessType formData string false "操作类型"
// @Param status formData string false "状态"
// @Param operTime formData string false "操作时间"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/operlog/list [GET]
func (o *OpelogController) List(ctx context.Context, c *app.RequestContext) {
	var req operLogParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListSysOperLogs(ctx, &v1.ListSysOperLogsRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		OperLog: v1.SysOperLog2OperLog(&req.SysOperLog),
		User:    &info,
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

// OperlogExport godoc
// @Summary 列出操作列表
// @Description 根据条件查询操作日志列表
// @Description 权限：system:operlog:export
// @Param operIp formData string false "操作ip"
// @Param title formData string false "系统模块"
// @Param operName formData string false "操作人员"
// @param businessType formData string false "操作类型"
// @Param status formData string false "状态"
// @Param operTime formData string false "操作时间"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/operlog/export [POST]
func (o *OpelogController) Export(ctx context.Context, c *app.RequestContext) {
	var req operLogParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ExportSysOperLog(ctx, &v1.ExportSysOperLogRequest{
		User:    &info,
		OperLog: v1.SysOperLog2OperLog(&req.SysOperLog),
	})

	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	utils.ExportExcel(c, resp.SheetName, v1.MOperLog2SysOperLog(resp.OperLogs))
}

// OperlogDelete godoc
// @Summary 删除操作列表
// @Description 删除多个操作记录
// @Description 权限：system:operlog:remove
// @Param operIps formData string false "操作ip"
// @Accept application/json
// @Produce application/json
// @Router /system/operlog/:operIds [DELETE]
func (o *OpelogController) Remove(ctx context.Context, c *app.RequestContext) {
	operIdsStr := c.Param("operIds")
	operIds := strutil.Strs2Int64(operIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysOperLog(ctx, &v1.DeleteSysOperLogRequest{
		OperIds: operIds,
		User:    &info,
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

// OperlogClean godoc
// @Summary 清空操作记录
// @Description 权限：system:operlog:remove
// @Accept application/json
// @Produce application/json
// @Router /system/operlog/clean [DELETE]
func (o *OpelogController) Clean(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.OperLogClean(ctx)
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
