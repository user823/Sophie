package service

import (
	"archive/zip"
	"bytes"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/internal/gen/engine"
	"github.com/user823/Sophie/internal/gen/store"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"io"
	"os"
	"path/filepath"
)

type GenTableSrv interface {
	// 查询业务列表
	SelectGenTableList(ctx context.Context, genTable *v1.GenTable, opts *api.GetOptions) *v1.GenTableList
	// 查询数据列表
	SelectDbTableList(ctx context.Context, genTable *v1.GenTable, opts *api.GetOptions) *v1.GenTableList
	// 查询数据库列表
	SelectDbTableListByNames(ctx context.Context, tableNames []string, opts *api.GetOptions) *v1.GenTableList
	// 查询所有表信息
	SelectGenTableAll(ctx context.Context, opts *api.GetOptions) *v1.GenTableList
	// 查询业务信息
	SelectGenTableById(ctx context.Context, id int64, opts *api.GetOptions) *v1.GenTable
	// 修改业务
	UpdateGenTable(ctx context.Context, genTable *v1.GenTable, opts *api.UpdateOptions) error
	// 删除业务信息
	DeleteGenTableByIds(ctx context.Context, tableIds []int64, opts *api.DeleteOptions) error
	// 导入表结构
	ImportGenTable(ctx context.Context, tableList []*v1.GenTable, operName string) error
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

type genTableService struct {
	store store.Factory
}

func NewGenTables(s store.Factory) GenTableSrv {
	return &genTableService{s}
}

func (s *genTableService) SelectGenTableList(ctx context.Context, genTable *v1.GenTable, opts *api.GetOptions) *v1.GenTableList {
	list, total, err := s.store.GenTables().SelectGenTableList(ctx, genTable, opts)
	if err != nil {
		return &v1.GenTableList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.GenTableList{
		ListMeta: api.ListMeta{total},
		Items:    list,
	}
}

func (s *genTableService) SelectDbTableList(ctx context.Context, genTable *v1.GenTable, opts *api.GetOptions) *v1.GenTableList {
	list, total, err := s.store.GenTables().SelectDbTableList(ctx, genTable, opts)
	if err != nil {
		return &v1.GenTableList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.GenTableList{
		ListMeta: api.ListMeta{total},
		Items:    list,
	}
}

func (s *genTableService) SelectDbTableListByNames(ctx context.Context, tableNames []string, opts *api.GetOptions) *v1.GenTableList {
	list, total, err := s.store.GenTables().SelectDbTableListByNames(ctx, tableNames, opts)
	if err != nil {
		log.Errorf("查找表失败: %s", err.Error())
		return &v1.GenTableList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.GenTableList{
		ListMeta: api.ListMeta{total},
		Items:    list,
	}
}

func (s *genTableService) SelectGenTableAll(ctx context.Context, opts *api.GetOptions) *v1.GenTableList {
	list, total, err := s.store.GenTables().SelectGenTableAll(ctx, opts)
	if err != nil {
		return &v1.GenTableList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.GenTableList{
		ListMeta: api.ListMeta{total},
		Items:    list,
	}
}

func (s *genTableService) SelectGenTableById(ctx context.Context, id int64, opts *api.GetOptions) *v1.GenTable {
	table, err := s.store.GenTables().SelectGenTableById(ctx, id, opts)
	if err != nil {
		log.Errorf("查询生成表失败: %s", err.Error())
		return nil
	}
	setTableFromOptions(table)
	return table
}

func (s *genTableService) UpdateGenTable(ctx context.Context, genTable *v1.GenTable, opts *api.UpdateOptions) error {
	// 设置table的Options
	if genTable.Extend["Params"] != nil {
		if str, ok := genTable.Extend["Params"].(string); ok {
			genTable.Options = str
		}
	}
	// 开启事务
	tx := s.store.Begin()
	if err := tx.GenTables().UpdateGenTable(ctx, genTable, opts); err != nil {
		log.Errorf("更新表失败: %s", err.Error())
		tx.Rollback()
		return err
	}
	for i := range genTable.Columns {
		if err := tx.GenTableColumns().UpdateGenTableColumn(ctx, genTable.Columns[i], opts); err != nil {
			log.Errorf("更新列失败: %s", err.Error())
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *genTableService) DeleteGenTableByIds(ctx context.Context, tableIds []int64, opts *api.DeleteOptions) error {
	// 开启事务
	tx := s.store.Begin()
	if err := tx.GenTables().DeleteGenTableByIds(ctx, tableIds, opts); err != nil {
		log.Errorf("删除表失败: %s", err.Error())
		tx.Rollback()
		return err
	}
	if err := tx.GenTableColumns().DeleteGenTableColumnByIds(ctx, tableIds, opts); err != nil {
		log.Errorf("删除列失败: %s", err.Error())
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *genTableService) ImportGenTable(ctx context.Context, tableList []*v1.GenTable, operName string) error {
	genHelper := engine.DefaultGenHelper()
	for i := range tableList {
		genHelper.InitTable(tableList[i], operName)
		// 开启事务
		tx := s.store.Begin()
		if err := tx.GenTables().InsertGenTable(ctx, tableList[i], &api.CreateOptions{}); err != nil {
			log.Errorf("插入表失败: %s", err.Error())
			tx.Rollback()
			return err
		}
		// 获取插入后的表信息
		table, _ := tx.GenTables().SelectGenTableByName(ctx, tableList[i].Tablename, &api.GetOptions{})

		// 保存列信息
		genTableColumns, _, err := s.store.GenTableColumns().SelectDbTableColumnsByName(ctx, table.Tablename, &api.GetOptions{})
		if err != nil {
			log.Errorf("根据表明获取列信息失败: %s", err.Error())
			tx.Rollback()
			continue
		}
		for j := range genTableColumns {
			genHelper.InitColumnField(table, genTableColumns[j])
			if err := tx.GenTableColumns().InsertGenTableColumn(ctx, genTableColumns[j], &api.CreateOptions{}); err != nil {
				log.Errorf("插入列失败: %s", err.Error())
				tx.Rollback()
				return err
			}
		}
		tx.Commit()
	}
	return nil
}

func (s *genTableService) PreviewCode(ctx context.Context, tableId int64, opts *api.GetOptions) map[string]string {
	dataMap := make(map[string]string)
	// 查询表信息
	table, _ := s.store.GenTables().SelectGenTableById(ctx, tableId, opts)
	// 设置主子表信息
	setSubTable(ctx, s.store, table)
	// 设置主键列信息
	setPkColumn(ctx, s.store, table)

	// 获取全局模板引擎
	tmplEngine := engine.GetTemplateEngineOr()
	tmplEngine.ParseContext(table)
	// 获取模板列表
	templates := tmplEngine.GetTemplateList(table.TplCategory, table.TplWebType)

	// 渲染模板
	for i := range templates {
		result, err := tmplEngine.ExecTemplate(templates[i], tableId)
		if err != nil {
			log.Errorf("模板渲染失败: %s", err.Error())
			continue
		}
		dataMap[templates[i]] = result
	}

	return dataMap
}

func (s *genTableService) DownloadCode(ctx context.Context, tableName string) []byte {
	var data []byte
	buffer := bytes.NewBuffer(data)
	err := generatorCode(ctx, s.store, buffer, tableName, false, []string{})
	if err != nil {
		return []byte{}
	}
	return buffer.Bytes()
}

// 自定义路径生成代码
func (s *genTableService) GeneratorCode(ctx context.Context, tableName string) error {
	excludePath := []string{"sql.template", "api.js.template", "index.vue.template", "index-tree.vue.template"}
	var data []byte
	buffer := bytes.NewBuffer(data)
	err := generatorCode(ctx, s.store, buffer, tableName, true, excludePath)
	if err != nil {
		return err
	}
	// viper获取项目保存路径地址
	genPath := viper.GetString("gen_path_root")
	if genPath == "" {
		genPath = "temp"
	}
	return os.WriteFile(genPath, buffer.Bytes(), os.ModePerm)
}

func (s *genTableService) SynchDb(ctx context.Context, tableName string) error {
	// 获取表信息
	table, _ := s.store.GenTables().SelectGenTableByName(ctx, tableName, &api.GetOptions{})
	columnList, total, _ := s.store.GenTableColumns().SelectDbTableColumnsByName(ctx, tableName, &api.GetOptions{})
	if total == 0 {
		return errors.New("同步数据库失败，原表结构不存在")
	}

	mp := make(map[string]*v1.GenTableColumn)
	for i := range table.Columns {
		mp[table.Columns[i].ColumnName] = table.Columns[i]
	}

	genHelper := engine.DefaultGenHelper()
	for i := range columnList {
		genHelper.InitColumnField(table, columnList[i])
		if prevColumn, ok := mp[columnList[i].ColumnName]; ok {
			columnList[i].ColumnId = prevColumn.ColumnId
			if columnList[i].IsList == "1" {
				// 如果是列表，继续保留查询方式/字典类型选线
				columnList[i].DictType = prevColumn.DictType
				columnList[i].QueryType = prevColumn.QueryType
			}

			if prevColumn.IsRequired != "" && !columnList[i].Pk() && (columnList[i].Insert() || columnList[i].Edit()) && (columnList[i].IsUsableColumn() || !columnList[i].IsSuperColumn()) {
				// 如果是（新增/修改&非主键/非忽略及父属性），继续保留必填/显示类型选项
				columnList[i].IsRequired = prevColumn.IsRequired
				columnList[i].HtmlType = prevColumn.HtmlType
			}
			err := s.store.GenTableColumns().UpdateGenTableColumn(ctx, columnList[i], &api.UpdateOptions{})
			if err != nil {
				log.Errorf("更新列失败: %s", err.Error())
			}
			delete(mp, columnList[i].ColumnName)
		} else {
			err := s.store.GenTableColumns().InsertGenTableColumn(ctx, columnList[i], &api.CreateOptions{})
			if err != nil {
				log.Errorf("生成列失败: %s", err.Error())
			}
		}
	}

	// 删除多余的列
	var delIds []int64
	for _, v := range mp {
		delIds = append(delIds, v.ColumnId)
	}
	return s.store.GenTableColumns().DeleteGenTableColumnByIds(ctx, delIds, &api.DeleteOptions{})
}

func (s *genTableService) DownloadCodes(ctx context.Context, tableNames []string) []byte {
	var data []byte
	buffer := bytes.NewBuffer(data)
	for i := range tableNames {
		generatorCode(ctx, s.store, buffer, tableNames[i], false, []string{})
	}
	return buffer.Bytes()
}

func (s *genTableService) ValidateEdit(ctx context.Context, genTable *v1.GenTable, opts *api.UpdateOptions) error {
	if genTable.TplCategory == v1.TPL_TREE {
		setTableFromOptions(genTable)
		if genTable.TreeCode == "" {
			return errors.New("树编码字段不能为空")
		}
		if genTable.TreeParentCode == "" {
			return errors.New("树父编码字段不能为空")
		}
		if genTable.TreeName == "" {
			return errors.New("树名称字段不能为空")
		}
	} else if genTable.TplCategory == v1.TPL_SUB {
		if genTable.SubTableName == "" {
			return errors.New("关联子表的表名不能为空")
		}
		if genTable.SubTableFkName == "" {
			return errors.New("子表关联的外键名不能为空")
		}
	}
	return nil
}

// 设置代码生成其他选项值
func setTableFromOptions(table *v1.GenTable) {
	if table.Options == "" {
		return
	}
	if err := jsoniter.Unmarshal(utils.S2b(table.Options), table); err != nil {
		log.Errorf("设置table其他选项失败: %s", err.Error())
	}
}

// 设置主子表信息
func setSubTable(ctx context.Context, store store.Factory, table *v1.GenTable) {
	if table.SubTableName != "" {
		if subTable, err := store.GenTables().SelectGenTableByName(ctx, table.SubTableName, &api.GetOptions{}); err == nil {
			table.SubTable = subTable
		}
	}
}

// 设置主键列信息
func setPkColumn(ctx context.Context, store store.Factory, table *v1.GenTable) {
	// 设置第一个主键列
	for i := range table.Columns {
		if table.Columns[i].IsPk == "1" {
			table.PkColumn = table.Columns[i]
			break
		}
	}

	// 如果主键列为空，则设置第一个列为主键列
	if table.PkColumn == nil && len(table.Columns) > 0 {
		table.PkColumn = table.Columns[0]
	}
	// 设置子表主键列
	if table.TplCategory == v1.TPL_SUB {
		for i := range table.SubTable.Columns {
			if table.SubTable.Columns[i].IsPk == "1" {
				table.SubTable.PkColumn = table.SubTable.Columns[i]
				break
			}
		}
		if table.SubTable.PkColumn == nil && len(table.SubTable.Columns) > 0 {
			table.SubTable.PkColumn = table.SubTable.Columns[0]
		}
	}
}

// 查询表信息并生成代码
func generatorCode(ctx context.Context, store store.Factory, w io.Writer, tableName string, useGenPath bool, excludePath []string) (err error) {
	// 查询表信息并生成代码
	table, err := store.GenTables().SelectGenTableByName(ctx, tableName, &api.GetOptions{})
	if err != nil {
		return err
	}
	// 设置主子表信息
	setSubTable(ctx, store, table)
	// 设置主键列信息
	setPkColumn(ctx, store, table)

	// 获取全局模板引擎
	tmplEngine := engine.GetTemplateEngineOr()
	tmplEngine.ParseContext(table)
	// 获取模板列表
	templates := tmplEngine.GetTemplateList(table.TplCategory, table.TplWebType)

	// 获取压缩器
	zw := zip.NewWriter(w)

	rootPath := v1.GEN_ROOT
	if useGenPath && table.GenPath != "" {
		rootPath = table.GenPath
	}

	// 渲染模板
	for i := range templates {
		if !strutil.ContainsAny(templates[i], excludePath...) {
			entry, err := zw.Create(filepath.Join(rootPath, engine.GetFileName(templates[i], table)))
			if err != nil {
				log.Errorf("创建zip文件失败: %s", err.Error())
				continue
			}
			err = tmplEngine.ExecTemplateW(templates[i], table.TableId, entry)
			if err != nil {
				log.Errorf("模板渲染失败: %s", err.Error())
				continue
			}
			zw.Flush()
		}
	}
	zw.Close()
	return
}
