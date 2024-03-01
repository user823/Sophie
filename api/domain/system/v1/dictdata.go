package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysDictData struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	DictCode       int64  `json:"dictCode,omitempty" gorm:"column:dict_code"`
	DictSort       int64  `json:"dictSort,omitempty" gorm:"column:dict_sort"`
	DictLabel      string `json:"dictLabel,omitempty" gorm:"column:dict_label"`
	DictValue      string `json:"dictValue,omitempty" gorm:"column:dict_value"`
	DictType       string `json:"dictType,omitempty" gorm:"column:dict_type"`
	CssClass       string `json:"cssClass,omitempty" gorm:"column:css_class"`
	ListClass      string `json:"listClass,omitempty" gorm:"column:list_class"`
	IsDefault      string `json:"isDefault,omitempty" gorm:"column:is_default"`
	Status         string `json:"status,omitempty" gorm:"column:status"`
}

func (s *SysDictData) TableName() string {
	return "sys_dict_data"
}

func (s *SysDictData) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysDictData) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type DictDataList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysDictData `json:"items"`
}

func (d *DictDataList) String() string {
	data, _ := jsoniter.Marshal(d)
	return utils.B2s(data)
}

func (d *DictDataList) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, d)
}
