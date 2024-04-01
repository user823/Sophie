package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/validators"
)

type SysDictData struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	DictCode       int64  `json:"dictCode,omitempty" gorm:"column:dict_code" query:"dictCode" xlsx:"n:字典编码"`
	DictSort       int64  `json:"dictSort,omitempty" gorm:"column:dict_sort" query:"dictSort" xlsx:"n:字典排序"`
	DictLabel      string `json:"dictLabel,omitempty" gorm:"column:dict_label" query:"dictLabel" validate:"required,min=0,max=100" xlsx:"n:字典标签"`
	DictValue      string `json:"dictValue,omitempty" gorm:"column:dict_value" query:"dictValue" validate:"required,min=0,max=100" xlsx:"n:字典键值"`
	DictType       string `json:"dictType,omitempty" gorm:"column:dict_type" query:"dictType" validate:"required,min=0,max=100" xlsx:"n:字典类型"`
	CssClass       string `json:"cssClass,omitempty" gorm:"column:css_class" validate:"min=0,max=100" query:"cssClass"`
	ListClass      string `json:"listClass,omitempty" gorm:"column:list_class" query:"listClass"`
	IsDefault      string `json:"isDefault,omitempty" gorm:"column:is_default" query:"isDefault" xlsx:"n:是否默认;exp:Y=是,N=否"`
	Status         string `json:"status,omitempty" gorm:"column:status" query:"status" xlsx:"n:状态;exp:0=正常,1=停用"`
}

func (s *SysDictData) TableName() string {
	return "sys_dict_data"
}

func (s *SysDictData) Marshal() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysDictData) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

func (s *SysDictData) Validate() error {
	vd := validators.GetValidatorOr()
	err := vd.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return buildDictDataErrMsg(e)
		}
	}
	return nil
}

func buildDictDataErrMsg(err validator.FieldError) error {
	switch err.StructNamespace() {
	case "SysDictData.DictLabel":
		return validators.BuildErrMsgHelper(err, "字典标签")
	case "SysDictData.DictValue":
		return validators.BuildErrMsgHelper(err, "字典键值")
	case "SysDictData.DictType":
		return validators.BuildErrMsgHelper(err, "字典类型")
	case "SysDictData.CssClass":
		return validators.BuildErrMsgHelper(err, "样式属性")
	}
	return nil
}

type DictDataList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysDictData `json:"items"`
}

func (d *DictDataList) Marshal() string {
	data, _ := jsoniter.Marshal(d)
	return utils.B2s(data)
}

func (d *DictDataList) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, d)
}
