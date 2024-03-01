package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysDictType struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	DictId         int64  `json:"dictId,omitempty" gorm:"column:dict_id"`
	DictName       string `json:"dictName,omitempty" gorm:"column:dict_name"`
	DictType       string `json:"dictType,omitempty" gorm:"column:dict_type"`
	Status         string `json:"status,omitempty" gorm:"column:status"`
}

func (s *SysDictType) TableName() string {
	return "sys_dict_type"
}

func (s *SysDictType) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysDictType) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type DictTypeList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysDictType `json:"items"`
}
