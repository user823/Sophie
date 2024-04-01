package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/validators"
	"regexp"
)

type SysDictType struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	DictId         int64  `json:"dictId,omitempty" gorm:"column:dict_id" query:"dictId" xlsx:"n:字典主键"`
	DictName       string `json:"dictName,omitempty" gorm:"column:dict_name" query:"dictName" validate:"required,min=0,max=30" xlsx:"n:字典名称"`
	DictType       string `json:"dictType,omitempty" gorm:"column:dict_type" query:"dictType" validate:"required,min=0,max=30" xlsx:"n:字典类型"`
	Status         string `json:"status,omitempty" gorm:"column:status" query:"status" xlsx:"n:状态;exp:0=正常,1=停用"`
}

func (s *SysDictType) TableName() string {
	return "sys_dict_type"
}

func (s *SysDictType) Marshal() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysDictType) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

func (s *SysDictType) Validate() error {
	vd := validators.GetValidatorOr()
	err := vd.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return buildDictTypeErrMsg(e)
		}
	}

	// 校验字典类型
	reg := regexp.MustCompile("^[a-z][a-z0-9_]*$")
	if !reg.MatchString(s.DictType) {
		return errors.New("字典类型必须以字母开头，且只能为（小写字母，数字，下滑线）")
	}
	return nil
}

func buildDictTypeErrMsg(err validator.FieldError) error {
	switch err.StructNamespace() {
	case "SysDictType.DictName":
		return validators.BuildErrMsgHelper(err, "字典名称")
	case "SysDictType.DictType":
		return validators.BuildErrMsgHelper(err, "字典类型")
	}
	return nil
}

type DictTypeList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysDictType `json:"items"`
}
