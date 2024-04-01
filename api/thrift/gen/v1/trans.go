package v1

import (
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/pkg/utils"
)

func GenTable2SysGenTable(genTable *GenTable) *v1.GenTable {
	if genTable == nil {
		return nil
	}

	return &v1.GenTable{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  genTable.CreateBy,
			CreatedAt: utils.Str2Time(genTable.CreateTime),
			UpdateBy:  genTable.UpdateBy,
			UpdatedAt: utils.Str2Time(genTable.UpdateTime),
			Remark:    genTable.Remark,
			Extend:    map[string]interface{}{"Params": genTable.Params},
		},
		TableId:        genTable.TableId,
		Tablename:      genTable.TableName,
		TableComment:   genTable.TableComment,
		SubTableName:   genTable.SubTableName,
		SubTableFkName: genTable.SubTableFkName,
		ClassName:      genTable.ClassName,
		TplCategory:    genTable.TplCategory,
		TplWebType:     genTable.TplWebType,
		PackageName:    genTable.PackageName,
		ModuleName:     genTable.ModuleName,
		BusinessName:   genTable.BusinessName,
		FunctionName:   genTable.FunctionName,
		FunctionAuthor: genTable.FunctionAuthor,
		GenType:        genTable.GenType,
		GenPath:        genTable.GenPath,
		PkColumn:       GenTableColumn2SysColumn(genTable.PkColumn), // trans
		SubTable:       GenTable2SysGenTable(genTable.SubTable),
		Columns:        MGenTableColumn2SysColumn(genTable.Columns), // trans
		Options:        genTable.Options,
		TreeCode:       genTable.TreeCode,
		TreeParentCode: genTable.TreeParentCode,
		TreeName:       genTable.TreeName,
		ParentMenuId:   genTable.ParentMenuId,
		ParentMenuName: genTable.ParentMenuName,
	}
}

func MGenTable2SysGenTable(tables []*GenTable) []*v1.GenTable {
	result := make([]*v1.GenTable, 0, len(tables))
	for i := range tables {
		result = append(result, GenTable2SysGenTable(tables[i]))
	}
	return result
}

func GenTableColumn2SysColumn(column *GenTableColumn) *v1.GenTableColumn {
	if column == nil {
		return nil
	}

	return &v1.GenTableColumn{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  column.CreateBy,
			CreatedAt: utils.Str2Time(column.CreateTime),
			UpdateBy:  column.UpdateBy,
			UpdatedAt: utils.Str2Time(column.UpdateTime),
			Extend:    map[string]interface{}{"Params": column.Params},
			Remark:    column.Remark,
		},
		ColumnId:      column.ColumnId,
		TableId:       column.TableId,
		ColumnName:    column.ColumnName,
		ColumnComment: column.ColumnComment,
		ColumnType:    column.ColumnType,
		GoType:        column.GoType,
		GoField:       column.GoField,
		IsPk:          column.IsPk,
		IsIncrement:   column.IsIncrement,
		IsRequired:    column.IsRequired,
		IsInsert:      column.IsInsert,
		IsEdit:        column.IsEdit,
		IsList:        column.IsList,
		IsQuery:       column.IsQuery,
		QueryType:     column.QueryType,
		HtmlType:      column.HtmlType,
		DictType:      column.DictType,
		Sort:          column.Sort,
	}
}

func MGenTableColumn2SysColumn(columns []*GenTableColumn) []*v1.GenTableColumn {
	result := make([]*v1.GenTableColumn, 0, len(columns))
	for i := range columns {
		result = append(result, GenTableColumn2SysColumn(columns[i]))
	}
	return result
}
