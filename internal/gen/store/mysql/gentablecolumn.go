package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/internal/gen/store"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/errors"
	"gorm.io/gorm"
)

type mysqlGenTableColumnStore struct {
	db *gorm.DB
}

var _ store.GenTableColumnStore = &mysqlGenTableColumnStore{}

func selectColumnVo(db *gorm.DB) *gorm.DB {
	return db.Table("gen_table_column")
}

func (s *mysqlGenTableColumnStore) SelectDbTableColumnsByName(ctx context.Context, tableName string, opts *api.GetOptions) ([]*v1.GenTableColumn, int64, error) {
	// 同时排除扩展字段
	query := s.db.Table("information_schema.columns").
		Select("COLUMN_NAME as column_name, "+
			"CASE WHEN (is_nullable = 'no' AND column_key != 'PRI') THEN '1' ELSE NULL END AS 'is_required', "+
			"CASE WHEN column_key = 'PRI' THEN '1' ELSE '0' END AS 'is_pk', "+
			"ordinal_position AS 'sort', "+
			"column_comment as column_comment, "+
			"CASE WHEN extra = 'auto_increment' THEN '1' ELSE '0' END AS 'is_increment', "+
			"column_type as column_type").Where("table_schema = (select database()) and table_name = ?", tableName).Where("COLUMN_NAME != ?", "extend_shadow")
	query = opts.SQLCondition(query, "").Order("ordinal_position")

	var result []*v1.GenTableColumn
	err := query.Find(&result).Error
	if err != nil {
		return []*v1.GenTableColumn{}, 0, err
	}

	return result, utils.CountQuery(query, opts, ""), nil
}

func (s *mysqlGenTableColumnStore) SelectGenTableColumnListByTableId(ctx context.Context, tableId int64, opts *api.GetOptions) ([]*v1.GenTableColumn, int64, error) {
	query := selectColumnVo(s.db).Where("table_id = ?", tableId)
	query = opts.SQLCondition(query, "").Order("sort")

	var result []*v1.GenTableColumn
	err := query.Find(&result).Error
	if err != nil {
		return []*v1.GenTableColumn{}, 0, err
	}
	return result, utils.CountQuery(query, opts, ""), nil
}

func (s *mysqlGenTableColumnStore) InsertGenTableColumn(ctx context.Context, column *v1.GenTableColumn, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(column)
	return create.Error
}

func (s *mysqlGenTableColumnStore) UpdateGenTableColumn(ctx context.Context, column *v1.GenTableColumn, opts *api.UpdateOptions) error {
	if column.ColumnId == 0 {
		return errors.New("更新列必须指定id")
	}
	update := opts.SQLCondition(s.db).Where("column_id = ?", column.ColumnId).Updates(column)
	return update.Error
}

func (s *mysqlGenTableColumnStore) DeleteGenTableColumns(ctx context.Context, columns []*v1.GenTableColumn, opts *api.DeleteOptions) error {
	var ids []int64
	for i := range columns {
		ids = append(ids, columns[i].ColumnId)
	}
	return s.DeleteGenTableColumnByIds(ctx, ids, opts)
}

func (s *mysqlGenTableColumnStore) DeleteGenTableColumnByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("table_id in ?", ids).Delete(&v1.GenTableColumn{})
	return del.Error
}
