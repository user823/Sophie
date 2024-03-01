package service

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
)

type GenTableColumnSrv interface {
	// 查询业务字段列表
	SelectGenTableColumnListByTableId(ctx context.Context, tableId int64, opts *api.GetOptions) []v1.GenTableColumn
	// 新增业务字段
	InsertGenTableColumn(ctx context.Context, genTableColumn *v1.GenTableColumn, opts *api.CreateOptions) error
	// 修改业务字段
	UpdateTableColumn(ctx context.Context, genTableColumn *v1.GenTableColumn, opts *api.UpdateOptions) error
	// 删除业务字段信息
	DeleteGenTableColumnByIds(ctx context.Context, ids string, opts *api.DeleteOptions) error
}
