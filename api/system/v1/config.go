package domain

import "github.com/user823/Sophie/api"

type SysConfig struct {
	api.ObjectMeta
	ConfigId    int64  `json:"configId" gorm:"column:config_id"`
	ConfigName  string `json:"configName" gorm:"column:config_name"`
	ConfigKey   string `json:"configKey" gorm:"column:config_key"`
	ConfigValue string `json:"configValue" gorm:"column:config_value"`
	// 是否系统内置
	ConfigType string `json:"configType" gorm:"column:config_type"`
}

func (s *SysConfig) TableName() string {
	return "sys_config"
}

type ConfigList struct {
	api.ListMeta `json:",inline"`
	Items        []SysConfig `json:"items"`
}
