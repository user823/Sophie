package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

type SysOperLog struct {
	// 日志主键
	OperId        int64     `json:"operId,omitempty" gorm:"column:oper_id" query:"operId"`
	Title         string    `json:"title,omitempty" gorm:"column:title" query:"title"`
	BusinessType  *int64    `json:"businessType,omitempty" gorm:"column:business_type" query:"businessType"`
	BusinessTypes []int64   `json:"businessTypes,omitempty" gorm:"-" query:"businessTypes"`
	Method        string    `json:"method,omitempty" gorm:"column:method" query:"method"`
	RequestMethod string    `json:"requestMethod,omitempty" gorm:"column:request_method" query:"requestMethod"`
	OperatorType  *int64    `json:"operatorType,omitempty" gorm:"column:operator_type" query:"operatorType"`
	OperName      string    `json:"operName,omitempty" gorm:"column:oper_name" query:"operName"`
	DeptName      string    `json:"deptName,omitempty" gorm:"column:dept_name" query:"deptName"`
	OperUrl       string    `json:"operUrl,omitempty" gorm:"column:oper_url" query:"operUrl"`
	OperIp        string    `json:"operIp,omitempty" gorm:"column:oper_ip" query:"operIp"`
	OperParam     string    `json:"operParam,omitempty" gorm:"column:oper_param" query:"operParam"`
	JsonResult    string    `json:"jsonResult,omitempty" gorm:"column:json_result" query:"jsonResult"`
	Status        string    `json:"status,omitempty" gorm:"column:status" query:"status"`
	ErrorMsg      string    `json:"errorMsg,omitempty" gorm:"column:error_msg" query:"errorMsg"`
	OperTime      time.Time `json:"operTime,omitempty" gorm:"column:oper_time" query:"operTime"`
	// 操作耗时。毫秒为单位
	CostTime int64 `json:"costTime,omitempty" gorm:"column:cost_time" query:"costTime"`
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
