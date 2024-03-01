package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

type SysOperLog struct {
	// 日志主键
	OperId        int64     `json:"operId,omitempty" gorm:"column:oper_id"`
	Title         string    `json:"title,omitempty" gorm:"column:title"`
	BusinessType  *int64    `json:"businessType,omitempty" gorm:"column:business_type"`
	BusinessTypes []int64   `json:"-" gorm:"-"`
	Method        string    `json:"method,omitempty" gorm:"column:method"`
	RequestMethod string    `json:"requestMethod,omitempty" gorm:"column:request_method"`
	OperatorType  *int64    `json:"operatorType,omitempty" gorm:"column:operator_type"`
	OperName      string    `json:"operName,omitempty" gorm:"column:oper_name"`
	DeptName      string    `json:"deptName,omitempty" gorm:"column:dept_name"`
	OperUrl       string    `json:"operUrl,omitempty" gorm:"column:oper_url"`
	OperIp        string    `json:"operIp,omitempty" gorm:"column:oper_ip"`
	OperParam     string    `json:"operParam,omitempty" gorm:"column:oper_param"`
	JsonResult    string    `json:"jsonResult,omitempty" gorm:"column:json_result"`
	Status        string    `json:"status,omitempty" gorm:"column:status"`
	ErrorMsg      string    `json:"errorMsg,omitempty" gorm:"column:error_msg"`
	OperTime      time.Time `json:"operTime,omitempty" gorm:"column:oper_time"`
	// 操作耗时。毫秒为单位
	CostTime int64 `json:"costTime,omitempty" gorm:"column:cost_time"`
}

func (s *SysOperLog) TableName() string {
	return "sys_oper_log"
}

func (s *SysOperLog) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysOperLog) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type OperLogList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysOperLog `json:"items"`
}
