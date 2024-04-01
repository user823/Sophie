package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

type SysLogininfor struct {
	InfoId     int64     `json:"infoId,omitempty" gorm:"column:info_id" query:"infoId" xlsx:"n:序号"`
	UserName   string    `json:"userName,omitempty" gorm:"column:user_name" query:"userName" xlsx:"n:用户账号"`
	Ipaddr     string    `json:"ipaddr,omitempty" gorm:"column:ipaddr" query:"ipaddr" xlsx:"n:状态;exp:0=成功,1=失败"`
	Status     string    `json:"status,omitempty" gorm:"column:status" query:"status" xlsx:"n:地址"`
	Msg        string    `json:"msg,omitempty" gorm:"column:msg" query:"msg" xlsx:"n:描述"`
	AccessTime time.Time `json:"accessTime,omitempty" gorm:"access_time" query:"accessTime" xlsx:"访问时间"`
}

func (l *SysLogininfor) TableName() string {
	return "sys_logininfor"
}

func (s *SysLogininfor) Marshal() string {
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
