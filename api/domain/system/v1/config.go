package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysConfig struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	ConfigId       int64  `json:"configId,omitempty" gorm:"column:config_id"`
	ConfigName     string `json:"configName,omitempty" gorm:"column:config_name"`
	ConfigKey      string `json:"configKey,omitempty" gorm:"column:config_key"`
	ConfigValue    string `json:"configValue,omitempty" gorm:"column:config_value"`
	// 是否系统内置
	ConfigType string `json:"configType" gorm:"column:config_type"`
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
