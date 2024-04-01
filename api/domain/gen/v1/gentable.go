package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/validators"
)

type GenTable struct {
	api.ObjectMeta
	TableId        int64  `json:"tableId" gorm:"column:table_id"`
	Tablename      string `json:"tableName" gorm:"column:table_name" validate:"required"`
	TableComment   string `json:"tableComment" gorm:"column:table_comment" validate:"required"`
	SubTableName   string `json:"subTableName" gorm:"column:sub_table_name"`
	SubTableFkName string `json:"subTableFkName" gorm:"column:sub_table_fk_name"`
	// 生成实体类的类名，驼峰表示法、首字母大写
	ClassName   string `json:"className" gorm:"column:class_name" validate:"required"`
	TplCategory string `json:"tplCategory" gorm:"column:tpl_category"`
	TplWebType  string `json:"tplWebType" gorm:"column:tpl_web_Type"`
	// 包路径, 根目录为internal，example: file、xx.job 的生成路径分别为internal/file、internal/xx/job 等
	PackageName string `json:"packageName" gorm:"column:package_name" validate:"required"`
	// 模块名, 项目内部的模块，example: file、system、job等，然后对应的微服务注册名为 Sophie File、Sophie System等
	ModuleName string `json:"moduleName" gorm:"column:module_name" validate:"required"`
	// 对应模块内部 具体业务名称，example: system.user, system.role
	BusinessName   string            `json:"businessName" gorm:"column:business_name" validate:"required"`
	FunctionName   string            `json:"functionName" gorm:"column:function_name" validate:"required"`
	FunctionAuthor string            `json:"functionAuthor" gorm:"column:function_author" validate:"required"`
	GenType        string            `json:"genType" gorm:"column:gen_type"`
	GenPath        string            `json:"genPath" gorm:"column:gen_path"`
	PkColumn       *GenTableColumn   `json:"pkColumn" gorm:"-"`
	SubTable       *GenTable         `json:"subTable" gorm:"-"`
	Columns        []*GenTableColumn `json:"columns" gorm:"foreignKey:TableId;references:TableId;order:sort"`
	// 代发生成的其他选项字段值，内容为json
	Options        string `json:"options" gorm:"column:options"`
	TreeCode       string `json:"treeCode" gorm:"-"`
	TreeParentCode string `json:"treeParentCode" gorm:"-"`
	TreeName       string `json:"treeName" gorm:"-"`
	ParentMenuId   string `json:"parentMenuId" gorm:"-"`
	ParentMenuName string `json:"parentMenuName" gorm:"-"`
}

func (g *GenTable) TableName() string {
	return "gen_table"
}

func (g *GenTable) Marshal() string {
	data, _ := jsoniter.Marshal(g)
	return utils.B2s(data)
}

func (g *GenTable) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, g)
}

func (g *GenTable) Validate() error {
	vd := validators.GetValidatorOr()
	err := vd.Struct(g)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return buildTableErrMsg(e)
		}
	}
	return nil
}

func buildTableErrMsg(err validator.FieldError) error {
	switch err.StructNamespace() {
	case "GenTable.Tablename":
		return validators.BuildErrMsgHelper(err, "表名称")
	case "GenTable.TableComment":
		return validators.BuildErrMsgHelper(err, "表描述")
	case "GenTable.ClassName":
		return validators.BuildErrMsgHelper(err, "实体类名称")
	case "GenTable.PackageName":
		return validators.BuildErrMsgHelper(err, "包生成路径")
	case "GenTable.ModuleName":
		return validators.BuildErrMsgHelper(err, "生成模块名")
	case "GenTable.BusinessName":
		return validators.BuildErrMsgHelper(err, "生成业务名")
	case "GenTable.FunctionName":
		return validators.BuildErrMsgHelper(err, "生成功能名")
	case "GenTable.FunctionAuthor":
		return validators.BuildErrMsgHelper(err, "作者名")
	}
	return nil
}

type GenTableList struct {
	api.ListMeta `json:",inline"`
	Items        []*GenTable `json:"items"`
}
