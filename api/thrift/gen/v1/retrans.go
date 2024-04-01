package v1

import (
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/pkg/utils"
)

func SysGenTable2GenTable(genTable *v1.GenTable) *GenTable {
	if genTable == nil {
		return nil
	}
	var params map[string]string
	if data, ok := genTable.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &GenTable{
		CreateBy:       genTable.CreateBy,
		CreateTime:     utils.Time2Str(genTable.CreatedAt),
		UpdateBy:       genTable.UpdateBy,
		UpdateTime:     utils.Time2Str(genTable.UpdatedAt),
		Remark:         genTable.Remark,
		Params:         params,
		TableId:        genTable.TableId,
		TableName:      genTable.Tablename,
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
		PkColumn:       SysColumn2GenTableColumn(genTable.PkColumn), // trans
		SubTable:       SysGenTable2GenTable(genTable.SubTable),
		Columns:        MSysTableColumn2GenTableColumn(genTable.Columns),
		Options:        genTable.Options,
		TreeCode:       genTable.TreeCode,
		TreeParentCode: genTable.TreeParentCode,
		TreeName:       genTable.TreeName,
		ParentMenuId:   genTable.ParentMenuId,
		ParentMenuName: genTable.ParentMenuName,
	}
}
func MSysGenTable2GenTable(genTables []*v1.GenTable) []*GenTable {
	result := make([]*GenTable, 0, len(genTables))
	for i := range genTables {
		result = append(result, SysGenTable2GenTable(genTables[i]))
	}
	return result
}

func SysColumn2GenTableColumn(column *v1.GenTableColumn) *GenTableColumn {
	if column == nil {
		return nil
	}

	var params map[string]string
	if data, ok := column.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &GenTableColumn{
		CreateBy:      column.CreateBy,
		CreateTime:    utils.Time2Str(column.CreatedAt),
		UpdateBy:      column.UpdateBy,
		UpdateTime:    utils.Time2Str(column.UpdatedAt),
		Params:        params,
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

func MSysTableColumn2GenTableColumn(columns []*v1.GenTableColumn) []*GenTableColumn {
	result := make([]*GenTableColumn, 0, len(columns))
	for i := range columns {
		result = append(result, SysColumn2GenTableColumn(columns[i]))
	}
	return result
}
