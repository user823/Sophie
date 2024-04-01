package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/internal/gen/store"
	"github.com/user823/Sophie/internal/gen/utils"
	"github.com/user823/Sophie/pkg/errors"
	"gorm.io/gorm"
)

type mysqlGenTableStore struct {
	db *gorm.DB
}

var _ store.GenTableStore = &mysqlGenTableStore{}

func selectGenTableVo(db *gorm.DB) *gorm.DB {
	return db.Model(&v1.GenTable{}).Preload("Columns").Table("gen_table t")
}

func (s *mysqlGenTableStore) SelectGenTableList(ctx context.Context, table *v1.GenTable, opts *api.GetOptions) ([]*v1.GenTable, int64, error) {
	query := selectGenTableVo(s.db)
	if table.Tablename != "" {
		query = query.Where("lower(table_name) like lower(?)", "%"+table.Tablename+"%")
	}
	if table.TableComment != "" {
		query = query.Where("lower(table_comment) like lower(?)", "%"+table.Tablename+"%")
	}
	query = opts.SQLCondition(query, "create_time")

	var result []*v1.GenTable
	err := query.Find(&result).Error
	if err != nil {
		return []*v1.GenTable{}, 0, err
	}
	return result, utils.CountQuery(query, opts, "create_time"), nil
}

func (s *mysqlGenTableStore) SelectDbTableList(ctx context.Context, table *v1.GenTable, opts *api.GetOptions) ([]*v1.GenTable, int64, error) {
	query := s.db.Table("information_schema.tables").Model(&v1.GenTable{}).Where("table_schema = (select database())").Where("" +
		"table_name not like 'gen_%'").Where("table_name not in (select table_name from gen_table)").Select("TABLE_NAME as table_name, TABLE_COMMENT as table_comment, CREATE_TIME as create_time, UPDATE_TIME as update_time")
	if table.Tablename != "" {
		query = query.Where("lower(table_name) like lower(?)", "%"+table.Tablename+"%")
	}
	if table.TableComment != "" {
		query = query.Where("lower(table_comment) like lower(?)", "%"+table.Tablename+"%")
	}
	query = opts.SQLCondition(query, "create_time").Order("create_time desc")

	var result []*v1.GenTable
	err := query.Find(&result).Error
	if err != nil {
		return []*v1.GenTable{}, 0, err
	}
	return result, utils.CountQuery(query, opts, "create_time"), nil
}

func (s *mysqlGenTableStore) SelectDbTableListByNames(ctx context.Context, tableNames []string, opts *api.GetOptions) ([]*v1.GenTable, int64, error) {
	query := s.db.Table("information_schema.tables").Where("table_schema = (select database())").Where("" +
		"table_name not like 'gen_%'").Where("table_name not in (select table_name from gen_table)").Select("TABLE_NAME as table_name, TABLE_COMMENT as table_comment, CREATE_TIME as create_time, UPDATE_TIME as update_time")
	query = query.Where("table_name in ?", tableNames)

	var result []*v1.GenTable
	err := query.Find(&result).Error
	if err != nil {
		return []*v1.GenTable{}, 0, err
	}
	return result, utils.CountQuery(query, opts, ""), nil
}

func (s *mysqlGenTableStore) SelectGenTableAll(ctx context.Context, opts *api.GetOptions) ([]*v1.GenTable, int64, error) {
	query := selectGenTableVo(s.db)
	query = opts.SQLCondition(query, "")

	var result []*v1.GenTable
	err := query.Find(&result).Error
	if err != nil {
		return []*v1.GenTable{}, 0, err
	}
	return result, utils.CountQuery(query, opts, ""), nil
}

func (s *mysqlGenTableStore) SelectGenTableById(ctx context.Context, id int64, opts *api.GetOptions) (*v1.GenTable, error) {
	query := selectGenTableVo(s.db).Where("t.table_id = ?", id)
	query = opts.SQLCondition(query, "")

	var result v1.GenTable
	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlGenTableStore) SelectGenTableByName(ctx context.Context, tableName string, opts *api.GetOptions) (*v1.GenTable, error) {
	query := selectGenTableVo(s.db).Where("t.table_name = ?", tableName)
	query = opts.SQLCondition(query, "")

	var result v1.GenTable
	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlGenTableStore) InsertGenTable(ctx context.Context, genTable *v1.GenTable, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(genTable)
	return create.Error
}

func (s *mysqlGenTableStore) UpdateGenTable(ctx context.Context, genTable *v1.GenTable, opts *api.UpdateOptions) error {
	if genTable.TableId == 0 {
		return errors.New("更新生成表必须指定id")
	}
	update := opts.SQLCondition(s.db).Where("table_id = ?", genTable.TableId).Model(genTable).Updates(genTable)
	return update.Error
}

func (s *mysqlGenTableStore) DeleteGenTableByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("table_id in ?", ids).Delete(&v1.GenTable{})
	return del.Error
}
