package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

type SysLogininfor struct {
	InfoId     int64     `json:"infoId,omitempty" gorm:"column:info_id"`
	UserName   string    `json:"userName,omitempty" gorm:"column:user_name"`
	Ipaddr     string    `json:"ipaddr,omitempty" gorm:"column:ipaddr"`
	Status     string    `json:"status,omitempty" gorm:"column:status"`
	Msg        string    `json:"msg,omitempty" gorm:"column:msg"`
	AccessTime time.Time `json:"accessTime,omitempty" gorm:"access_time"`
}

func (l *SysLogininfor) TableName() string {
	return "sys_logininfor"
}

func (s *SysLogininfor) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysLogininfor) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type LogininforList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysLogininfor `json:"items"`
}
