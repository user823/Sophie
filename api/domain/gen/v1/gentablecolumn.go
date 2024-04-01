package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"github.com/user823/Sophie/pkg/validators"
	"gorm.io/gorm"
)

type GenTableColumn struct {
	api.ObjectMeta
	ColumnId      int64  `json:"columnId" gorm:"column:column_id"`
	TableId       int64  `json:"tableId" gorm:"column:table_id"`
	ColumnName    string `json:"columnName" gorm:"column:column_name"`
	ColumnComment string `json:"columnComment" gorm:"column:column_comment"`
	ColumnType    string `json:"columnType" gorm:"column:column_type"`
	GoType        string `json:"goType" gorm:"column:go_type"`
	GoField       string `json:"goField" gorm:"column:go_field" validate:"required"`
	IsPk          string `json:"isPk" gorm:"column:is_pk"`
	IsIncrement   string `json:"isIncrement" gorm:"column:is_increment"`
	IsRequired    string `json:"isRequired" gorm:"column:is_required"`
	IsInsert      string `json:"isInsert" gorm:"column:is_insert"`
	IsEdit        string `json:"isEdit" gorm:"column:is_edit"`
	IsList        string `json:"isList" gorm:"column:is_list"`
	IsQuery       string `json:"isQuery" gorm:"column:is_query"`
	QueryType     string `json:"queryType" gorm:"column:query_type"`
	HtmlType      string `json:"htmlType" gorm:"column:html_type"`
	DictType      string `json:"dictType" gorm:"column:dict_type"`
	Sort          int64  `json:"sort" gorm:"column:sort"`
}

func (g *GenTableColumn) TableName() string {
	return "gen_table_column"
}

func (obj *GenTableColumn) BeforeCreate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// BeforeUpdate run before update database record.
func (obj *GenTableColumn) BeforeUpdate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
func (obj *GenTableColumn) AfterFind(tx *gorm.DB) error {
	if obj.ExtendShadow == "" {
		return nil
	}
	if err := jsoniter.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	return nil
}
func (g *GenTableColumn) Marshal() string {
	data, _ := jsoniter.Marshal(g)
	return utils.B2s(data)
}

func (g *GenTableColumn) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, g)
}

func (g *GenTableColumn) Validate() error {
	vd := validators.GetValidatorOr()
	err := vd.Struct(g)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return buildColumnErrMsg(e)
		}
	}
	return nil
}

func buildColumnErrMsg(err validator.FieldError) error {
	switch err.StructNamespace() {
	case "GenTableColumn.GoField":
		return validators.BuildErrMsgHelper(err, "golang 属性不能为空")
	}
	return nil
}

func (g *GenTableColumn) Pk() bool {
	return g.IsPk == "1"
}

func (g *GenTableColumn) Insert() bool {
	return g.IsInsert == "1"
}

func (g *GenTableColumn) Edit() bool {
	return g.IsEdit == "1"
}

func (g *GenTableColumn) IsSuperColumn() bool {
	return strutil.ContainsAnyIgnoreCase(g.GoField, "createBy", "createTime", "updateBy", "updateTime", "remark",
		"parentName", "parentId", "orderNum", "ancestors")
}

func (g *GenTableColumn) IsUsableColumn() bool {
	return strutil.ContainsAnyIgnoreCase(g.GoField, "parentId", "orderNum", "remark")
}

func (g *GenTableColumn) Require() bool {
	return g.IsRequired == REQUIRE
}

type GenTableColumnList struct {
	api.ListMeta `json:",inline"`
	Items        []*GenTableColumn `json:"items"`
}
