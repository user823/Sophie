package service

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
)

type GenTableSrv interface {
	// 查询业务列表
	SelectGenTableList(ctx context.Context, genTable *v1.GenTable, opts *api.GetOptions) []v1.GenTable
	// 查询数据列表
	SelectDbTableList(ctx context.Context, genTable *v1.GenTable, opts *api.GetOptions) []v1.GenTable
	// 查询数据库列表
	SelectDbTableListByNames(ctx context.Context, tableNames []string, opts *api.GetOptions) []v1.GenTable
	// 查询所有表信息
	SelectGenTableAll(ctx context.Context, opts *api.GetOptions) []v1.GenTable
	// 查询业务信息
	SelectGenTableById(ctx context.Context, id int64, opts *api.GetOptions) v1.GenTable
	// 修改业务
	UpdateGenTable(ctx context.Context, genTable *v1.GenTable, opts *api.UpdateOptions) error
	// 删除业务信息
	DeleteGenTableByIds(ctx context.Context, tableIds []int64, opts *api.DeleteOptions) error
	// 导入表结构
	ImportGenTable(ctx context.Context, tableList []v1.GenTable) error
	// 预览代码
	PreviewCode(ctx context.Context, tableId int64, opts *api.GetOptions) map[string]string
	// 生成代码（下载方式）
	DownloadCode(ctx context.Context, tableName string) []byte
	// 生成代码（自定义路径)
	GeneratorCode(ctx context.Context, tableName string) error
	// 同步数据库
	SynchDb(ctx context.Context, tableName string) error
	// 批量生成代码(下载方式）
	DownloadCodes(ctx context.Context, tableNames []string) []byte
	// 修改保存参数校验
	ValidateEdit(ctx context.Context, genTable *v1.GenTable, opts *api.UpdateOptions) error
}
