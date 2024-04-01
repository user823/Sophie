package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/internal/gen/store"
	"github.com/user823/Sophie/pkg/utils/strutil"
)

type GenTableColumnSrv interface {
	// 查询业务字段列表
	SelectGenTableColumnListByTableId(ctx context.Context, tableId int64, opts *api.GetOptions) *v1.GenTableColumnList
	// 新增业务字段
	InsertGenTableColumn(ctx context.Context, genTableColumn *v1.GenTableColumn, opts *api.CreateOptions) error
	// 修改业务字段
	UpdateTableColumn(ctx context.Context, genTableColumn *v1.GenTableColumn, opts *api.UpdateOptions) error
	// 删除业务字段信息
	DeleteGenTableColumnByIds(ctx context.Context, ids string, opts *api.DeleteOptions) error
}

type genTableColumnService struct {
	store store.Factory
}

func NewColumns(s store.Factory) GenTableColumnSrv {
	return &genTableColumnService{s}
}

func (s *genTableColumnService) SelectGenTableColumnListByTableId(ctx context.Context, tableId int64, opts *api.GetOptions) *v1.GenTableColumnList {
	list, total, err := s.store.GenTableColumns().SelectGenTableColumnListByTableId(ctx, tableId, opts)
	if err != nil {
		return &v1.GenTableColumnList{ListMeta: api.ListMeta{0}}
	}
	return &v1.GenTableColumnList{
		ListMeta: api.ListMeta{total},
		Items:    list,
	}
}

func (s *genTableColumnService) InsertGenTableColumn(ctx context.Context, genTableColumn *v1.GenTableColumn, opts *api.CreateOptions) error {
	return s.store.GenTableColumns().InsertGenTableColumn(ctx, genTableColumn, opts)
}

func (s *genTableColumnService) UpdateTableColumn(ctx context.Context, genTableColumn *v1.GenTableColumn, opts *api.UpdateOptions) error {
	return s.store.GenTableColumns().UpdateGenTableColumn(ctx, genTableColumn, opts)
}

func (s *genTableColumnService) DeleteGenTableColumnByIds(ctx context.Context, ids string, opts *api.DeleteOptions) error {
	columnIds := strutil.Strs2Int64(ids)
	return s.store.GenTableColumns().DeleteGenTableColumnByIds(ctx, columnIds, opts)
}
