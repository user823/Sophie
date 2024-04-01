package gen

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
	"github.com/user823/Sophie/internal/gen/service"
	"github.com/user823/Sophie/internal/gen/store/mysql"
	"github.com/user823/Sophie/internal/gen/utils"
	"github.com/user823/Sophie/pkg/log"
	"strings"
)

// GenServiceImpl implements the last service interface defined in the IDL.
type GenServiceImpl struct{}

// ListGenTables implements the GenServiceImpl interface.
func (s *GenServiceImpl) ListGenTables(ctx context.Context, req *v1.ListGenTablesRequest) (resp *v1.ListGenTablesResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), req.GetDateRange(), false)
	store, _ := mysql.GetMySQLFactoryOr(nil)
	sysGenTable := v1.GenTable2SysGenTable(req.GenTable)

	list := service.NewGenTables(store).SelectGenTableList(ctx, sysGenTable, getOpt)
	return &v1.ListGenTablesResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysGenTable2GenTable(list.Items),
	}, nil
}

// GetInfo implements the GenServiceImpl interface.
func (s *GenServiceImpl) GetInfo(ctx context.Context, tableId int64) (resp *v1.GetInfoResponse, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	tableSrv := service.NewGenTables(store)
	columnSrv := service.NewColumns(store)

	table := tableSrv.SelectGenTableById(ctx, tableId, &api.GetOptions{})
	tables := tableSrv.SelectGenTableAll(ctx, &api.GetOptions{})
	columns := columnSrv.SelectGenTableColumnListByTableId(ctx, tableId, &api.GetOptions{})
	return &v1.GetInfoResponse{
		BaseResp: utils.Ok("操作成功"),
		Info:     v1.SysGenTable2GenTable(table),
		Rows:     v1.MSysTableColumn2GenTableColumn(columns.Items),
		Tables:   v1.MSysGenTable2GenTable(tables.Items),
	}, nil
}

// DataList implements the GenServiceImpl interface.
func (s *GenServiceImpl) DataList(ctx context.Context, req *v1.DataListRequest) (resp *v1.ListGenTablesResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), req.GetDateRange(), false)
	store, _ := mysql.GetMySQLFactoryOr(nil)
	sysGenTable := v1.GenTable2SysGenTable(req.GenTable)
	list := service.NewGenTables(store).SelectDbTableList(ctx, sysGenTable, getOpt)
	return &v1.ListGenTablesResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysGenTable2GenTable(list.Items),
	}, nil
}

// ColumnList implements the GenServiceImpl interface.
func (s *GenServiceImpl) ColumnList(ctx context.Context, tableId int64) (resp *v1.ListGenTableColumnsResponse, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	list := service.NewColumns(store).SelectGenTableColumnListByTableId(ctx, tableId, &api.GetOptions{})
	return &v1.ListGenTableColumnsResponse{
		BaseResp: utils.Ok("操作成功"),
		Rows:     v1.MSysTableColumn2GenTableColumn(list.Items),
		Total:    list.TotalCount,
	}, nil
}

// ImportTableSave implements the GenServiceImpl interface.
func (s *GenServiceImpl) ImportTableSave(ctx context.Context, req *v1.ImportTableSaveRequest) (resp *v1.BaseResp, err error) {
	tables := strings.Split(req.Tables, ",")
	store, _ := mysql.GetMySQLFactoryOr(nil)
	tableSrv := service.NewGenTables(store)

	// 查询表信息
	tableList := tableSrv.SelectDbTableListByNames(ctx, tables, &api.GetOptions{})
	if err = tableSrv.ImportGenTable(ctx, tableList.Items, req.OperName); err != nil {
		log.Errorf("导入表失败: %s", err.Error())
		return utils.Fail("导入失败，系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// EditSave implements the GenServiceImpl interface.
func (s *GenServiceImpl) EditSave(ctx context.Context, req *v1.EditSaveRequest) (resp *v1.BaseResp, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	tableSrv := service.NewGenTables(store)
	sysTable := v1.GenTable2SysGenTable(req.GenTable)

	if err = tableSrv.ValidateEdit(ctx, sysTable, &api.UpdateOptions{}); err != nil {
		return utils.Fail(err.Error()), nil
	}
	if err = tableSrv.UpdateGenTable(ctx, sysTable, &api.UpdateOptions{}); err != nil {
		log.Errorf("更新生成表失败: %s", err.Error())
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// Remove implements the GenServiceImpl interface.
func (s *GenServiceImpl) Remove(ctx context.Context, req *v1.RemoveRequest) (resp *v1.BaseResp, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	if err = service.NewGenTables(store).DeleteGenTableByIds(ctx, req.TableIds, &api.DeleteOptions{}); err != nil {
		log.Errorf("移除表失败: %s", err.Error())
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// Preview implements the GenServiceImpl interface.
func (s *GenServiceImpl) Preview(ctx context.Context, tableId int64) (resp *v1.PreviewResponse, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	dataMp := service.NewGenTables(store).PreviewCode(ctx, tableId, &api.GetOptions{})
	return &v1.PreviewResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     dataMp,
	}, nil
}

// Download implements the GenServiceImpl interface.
func (s *GenServiceImpl) Download(ctx context.Context, tableName string) (resp *v1.DownloadResponse, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	data := service.NewGenTables(store).DownloadCode(ctx, tableName)
	return &v1.DownloadResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     data,
	}, nil
}

// GenCode implements the GenServiceImpl interface.
func (s *GenServiceImpl) GenCode(ctx context.Context, tableName string) (resp *v1.BaseResp, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	if err = service.NewGenTables(store).GeneratorCode(ctx, tableName); err != nil {
		log.Errorf("生成代码失败: %s", err.Error())
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// SynchDb implements the GenServiceImpl interface.
func (s *GenServiceImpl) SynchDb(ctx context.Context, tableName string) (resp *v1.BaseResp, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	if err = service.NewGenTables(store).SynchDb(ctx, tableName); err != nil {
		log.Errorf("数据库同步失败: %s", err.Error())
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// BatchGenCode implements the GenServiceImpl interface.
func (s *GenServiceImpl) BatchGenCode(ctx context.Context, tables string) (resp *v1.DownloadResponse, err error) {
	tableNames := strings.Split(tables, ",")
	store, _ := mysql.GetMySQLFactoryOr(nil)
	data := service.NewGenTables(store).DownloadCodes(ctx, tableNames)
	return &v1.DownloadResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     data,
	}, nil
}
