package v1

import (
	"github.com/user823/Sophie/api"
	"time"
)

type GenTableColumn struct {
	ColumnId      int64     `json:"columnId" gorm:"column:column_id"`
	TableId       int64     `json:"tableId" gorm:"column:table_id"`
	ColumnName    string    `json:"columnName" gorm:"column:column_name"`
	ColumnComment string    `json:"columnComment" gorm:"column:column_comment"`
	ColumnType    string    `json:"columnType" gorm:"column:column_type"`
	GoType        string    `json:"goType" gorm:"column:go_type"`
	GoField       string    `json:"goField" gorm:"column:go_field"`
	IsPk          string    `json:"isPk" gorm:"column:is_pk"`
	IsIncrement   string    `json:"isIncrement" gorm:"column:is_increment"`
	IsRequired    string    `json:"isRequired" gorm:"column:is_required"`
	IsInsert      string    `json:"isInsert" gorm:"column:is_insert"`
	IsEdit        string    `json:"isEdit" gorm:"column:is_edit"`
	IsList        string    `json:"isList" gorm:"column:is_list"`
	IsQuery       string    `json:"isQuery" gorm:"column:is_query"`
	QueryType     string    `json:"queryType" gorm:"column:query_type"`
	HtmlType      string    `json:"htmlType" gorm:"column:html_type"`
	DictType      string    `json:"dictType" gorm:"column:dict_type"`
	Sort          int64     `json:"sort" gorm:"column:sort"`
	CreateBy      string    `json:"createBy" gorm:"column:create_by"`
	CreatedAt     time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateBy      string    `json:"updateBy" gorm:"column:update_time"`
	UpdatedAt     time.Time `json:"updateTime" gorm:"column:update_time"`
	ExtendShadow  string    `json:"extendShadow" gorm:"column:extend_shadow"`
}

func (g *GenTableColumn) TableName() string {
	return "gen_table_column"
}

type GenTableColumnList struct {
	api.ListMeta `json:",inline"`
	Items        []GenTableColumn `json:"items"`
}
