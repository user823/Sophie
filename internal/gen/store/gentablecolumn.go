package store

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
)

type GenTableColumnStore interface {
	// 根据表名称查询列信息
	SelectDbTableColumnsByName(ctx context.Context, tableName string, opts *api.GetOptions) ([]v1.GenTableColumn, error)
	// 查询业务字段列表
	SelectGenTableColumnListByTableId(ctx context.Context, tableId int64, opts *api.GetOptions) ([]v1.GenTableColumn, error)
	// 新增业务字段
	InsertGenTableColumn(ctx context.Context, column *v1.GenTableColumn, opts *api.CreateOptions) error
	// 修改业务字段
	UpdateGenTableColumn(ctx context.Context, column *v1.GenTableColumn, opts *api.UpdateOptions) error
	// 删除业务字段
	DeleteGenTableColumns(ctx context.Context, columns []v1.GenTableColumn, opts *api.DeleteOptions) error
	// 批量删除业务字段
	DeleteGenTableColumnByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
}
