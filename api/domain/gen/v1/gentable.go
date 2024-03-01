package v1

import "github.com/user823/Sophie/api"

type GenTable struct {
	api.ObjectMeta
	TableId        int64            `json:"tableId" gorm:"column:table_id"`
	Tablename      string           `json:"tableName" gorm:"column:table_name"`
	TableComment   string           `json:"tableComment" gorm:"column:table_comment"`
	SubTableName   string           `json:"subTableName" gorm:"column:sub_table_name"`
	SubTableFkName string           `json:"subTableFkName" gorm:"column:sub_table_fk_name"`
	ClassName      string           `json:"className" gorm:"column:class_name"`
	TplCategory    string           `json:"tplCategory" gorm:"column:tpl_category"`
	TplWebType     string           `json:"tplWebType" gorm:"column:tpl_webType"`
	PackageName    string           `json:"packageName" gorm:"column:package_name"`
	ModuleName     string           `json:"moduleName" gorm:"column:module_name"`
	BusinessName   string           `json:"businessName" gorm:"column:business_name"`
	FunctionName   string           `json:"functionName" gorm:"column:function_name"`
	FunctionAuthor string           `json:"functionAuthor" gorm:"column:function_author"`
	GenType        string           `json:"genType" gorm:"column:gen_type"`
	GenPath        string           `json:"genPath" gorm:"column:gen_path"`
	PkColumn       GenTableColumn   `json:"pkColumn" gorm:"-"`
	SubTable       *GenTable        `json:"subTable" gorm:"-"`
	Columns        []GenTableColumn `json:"columns" gorm:"-"`
	Options        string           `json:"options" gorm:"column:options"`
	TreeCode       string           `json:"treeCode" gorm:"-"`
	TreeParentCode string           `json:"treeParentCode" gorm:"-"`
	TreeName       string           `json:"treeName" gorm:"-"`
	ParentMenuId   string           `json:"parentMenuId" gorm:"-"`
	ParentMenuName string           `json:"parentMenuName" gorm:"-"`
}

func (g *GenTable) TableName() string {
	return "gen_table"
}

type GenTableList struct {
	api.ListMeta `json:",inline"`
	Items        []GenTable `json:"items"`
}
