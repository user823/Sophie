package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysConfig struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	ConfigId       int64  `json:"configId,omitempty" gorm:"column:config_id" query:"configId"`
	ConfigName     string `json:"configName,omitempty" gorm:"column:config_name" query:"configName"`
	ConfigKey      string `json:"configKey,omitempty" gorm:"column:config_key" query:"configKey"`
	ConfigValue    string `json:"configValue,omitempty" gorm:"column:config_value" query:"configValue"`
	// 是否系统内置
	ConfigType string `json:"configType,omitempty" gorm:"column:config_type" query:"configType"`
}

func (s *SysConfig) TableName() string {
	return "sys_config"
}

func (s *SysConfig) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysConfig) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type ConfigList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysConfig `json:"items"`
}
