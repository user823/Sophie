package store

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
)

type GenTableStore interface {
	// 查询业务列表
	SelectGenTableList(ctx context.Context, table *v1.GenTable, opts *api.GetOptions) ([]v1.GenTable, error)
	// 查询数据库列表
	SelectDbTableList(ctx context.Context, table *v1.GenTable, opts *api.GetOptions) ([]v1.GenTable, error)
	// 查询数据库列表
	SelectDbTableListByNames(ctx context.Context, tableNames []string, opts *api.GetOptions) ([]v1.GenTable, error)
	// 查询所有表信息
	SelectGenTableAll(ctx context.Context, opts *api.GetOptions) ([]v1.GenTable, error)
	// 查询表id业务信息
	SelectGenTableById(ctx context.Context, id int64, opts *api.GetOptions) (v1.GenTable, error)
	// 查询表名称业务信息
	SelectGenTableByName(ctx context.Context, tableName string, opts *api.GetOptions) (v1.GenTable, error)
	// 新增业务
	InsertGenTable(ctx context.Context, genTable *v1.GenTable, opts *api.CreateOptions) error
	// 修改业务
	UpdateGenTable(ctx context.Context, genTable *v1.GenTable, opts *api.UpdateOptions) error
	// 批量删除业务
	DeleteGenTableByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
}
