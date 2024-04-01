package gen

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	v12 "github.com/user823/Sophie/api/thrift/gen/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strconv"
)

type GenController struct{}

func NewGenController() *GenController {
	return &GenController{}
}

type genRequestParam struct {
	v1.GenTable
	api.GetOptions
	Tables string `json:"tables" query:"tables"`
}

// TableList godoc
// @Summary 代码生成
// @Description 获取导入的数据库表列表
// @Description 权限：tool:gen:list
// @Param tableName formData string false "表名称"
// @Param tableComment formData string false "表描述"
// @Param createTime formData string false "创建时间"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/list [GET]
func (g *GenController) List(ctx context.Context, c *app.RequestContext) {
	var req genRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListGenTables(ctx, &v12.ListGenTablesRequest{
		GenTable: v12.SysGenTable2GenTable(&req.GenTable),
		DateRange: &v12.DateRange{
			BeginTime: req.BeginTime,
			EndTime:   req.EndTime,
		},
		PageInfo: &v12.PageInfo{
			PageSize:      req.PageSize,
			PageNum:       req.PageNum,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	res := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"rows":  resp.Rows,
		"total": resp.Total,
	}
	core.JSON(c, res)
}

// TableInfo godoc
// @Summary 数据表详情
// @Description 获取导入的数据库表详情
// @Description 权限：tool:gen:query
// @Param tableId query string false "数据库表id"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/:tableId [GET]
func (g *GenController) GetInfo(ctx context.Context, c *app.RequestContext) {
	tableIdStr := c.Param("tableId")
	tableId, _ := strconv.ParseInt(tableIdStr, 10, 64)

	resp, err := rpc.Remoting.GetInfo(ctx, tableId)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	res := map[string]any{
		"code": resp.BaseResp.Code,
		"msg":  resp.BaseResp.Msg,
		"data": map[string]any{
			"info":   resp.Info,
			"rows":   resp.Rows,
			"tables": resp.Tables,
		},
	}
	core.JSON(c, res)
}

// ImportTableList godoc
// @Summary 获取可导入的数据表
// @Description 权限：tool:gen:list
// @Param tableName formData string false "表名称"
// @Param tableComment formData string false "表描述"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/db/list [GET]
func (g *GenController) DataList(ctx context.Context, c *app.RequestContext) {
	var req genRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DataList(ctx, &v12.DataListRequest{
		GenTable: v12.SysGenTable2GenTable(&req.GenTable),
		DateRange: &v12.DateRange{
			BeginTime: req.BeginTime,
			EndTime:   req.EndTime,
		},
		PageInfo: &v12.PageInfo{
			PageSize:      req.PageSize,
			PageNum:       req.PageNum,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	res := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"rows":  resp.Rows,
		"total": resp.Total,
	}
	core.JSON(c, res)
}

// TableColumnList godoc
// @Summary	数据库列字段列表
// @Description 权限：tool:gen:list
// @Param tableId query int true "表id"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/column/:tableId [GET]
func (g *GenController) ColumnList(ctx context.Context, c *app.RequestContext) {
	tableIdStr := c.Param("tableId")
	tableId, _ := strconv.ParseInt(tableIdStr, 10, 64)

	resp, err := rpc.Remoting.ColumnList(ctx, tableId)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	res := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"rows":  resp.Rows,
		"total": resp.Total,
	}
	core.JSON(c, res)
}

// ImportTable godoc
// @Summary	导入数据库表
// @Description 权限：tool:gen:import
// @Param tableName query string true "数据库表名称"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/importTable [POST]
func (g *GenController) ImportTableSave(ctx context.Context, c *app.RequestContext) {
	var req genRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ImportTableSave(ctx, &v12.ImportTableSaveRequest{
		Tables:   req.Tables,
		OperName: info.User.UserName,
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

// EditSave godoc
// @Summary	导入数据库表
// @Description 权限：tool:gen:edit
// @Param tableName query string true "数据库表名称"
// @Accept application/json
// @Produce application/json
// @Router /code/gen [PUT]
func (g *GenController) EditSave(ctx context.Context, c *app.RequestContext) {
	var req genRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.EditSave(ctx, &v12.EditSaveRequest{
		GenTable: v12.SysGenTable2GenTable(&req.GenTable),
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

// Remove godoc
// @Summary	删除数据库表
// @Description 权限：tool:gen:remove
// @Param tableIds query string true "数据库表id"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/:tableIds [DELETE]
func (g *GenController) Remove(ctx context.Context, c *app.RequestContext) {
	tableIdsStr := c.Param("tableIds")
	tableIds := strutil.Strs2Int64(tableIdsStr)

	resp, err := rpc.Remoting.Remove(ctx, &v12.RemoveRequest{
		TableIds: tableIds,
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

// Preview godoc
// @Summary	预览数据库表
// @Description 权限：tool:gen:preview
// @Param tableId query string true "数据库表id"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/preview/:tableId [GET]
func (g *GenController) Preview(ctx context.Context, c *app.RequestContext) {
	tableIdStr := c.Param("tableId")
	tableId, _ := strconv.ParseInt(tableIdStr, 10, 64)

	resp, err := rpc.Remoting.Preview(ctx, tableId)
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

// Preview godoc
// @Summary	预览数据库表
// @Description 权限：tool:gen:code
// @Param tableName query string true "数据库表名称"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/download/:tableName [GET]
func (g *GenController) Download(ctx context.Context, c *app.RequestContext) {
	tableName := c.Param("tableName")

	resp, err := rpc.Remoting.Download(ctx, tableName)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	utils.GenCode(c, resp.Data)
}

// GenCode godoc
// @Summary	生成代码
// @Description 权限：tool:gen:code
// @Param tableName query string true "数据库表名称"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/genCode/:tableName [GET]
func (g *GenController) GenCode(ctx context.Context, c *app.RequestContext) {
	tableName := c.Param("tableName")

	resp, err := rpc.Remoting.GenCode(ctx, tableName)
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

// SyncDb godoc
// @Summary	同步数据库
// @Description 权限：tool:gen:edit
// @Param tableName query string true "数据库表名称"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/synchDb/:tableName [GET]
func (g *GenController) SynchDb(ctx context.Context, c *app.RequestContext) {
	tableName := c.Param("tableName")

	resp, err := rpc.Remoting.SynchDb(ctx, tableName)
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

// BatchGenCode godoc
// @Summary	批量下载
// @Description 权限：tool:gen:code
// @Param tables query string true "数据库表id"
// @Accept application/json
// @Produce application/json
// @Router /code/gen/batchGenCode [GET]
func (g *GenController) BatchGenCode(ctx context.Context, c *app.RequestContext) {
	tables := c.Query("tables")

	resp, err := rpc.Remoting.BatchGenCode(ctx, tables)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	utils.GenCode(c, resp.Data)
}
