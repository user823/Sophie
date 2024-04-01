package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

type SysOperLog struct {
	// 日志主键
	OperId        int64     `json:"operId,omitempty" gorm:"column:oper_id" query:"operId" xlsx:"n:操作序号"`
	Title         string    `json:"title,omitempty" gorm:"column:title" query:"title" xlsx:"n:操作模块"`
	BusinessType  *int64    `json:"businessType,omitempty" gorm:"column:business_type" query:"businessType" xlsx:"n:业务类型;exp:0=空,1=新增,2=修改,3=删除,4=授权,5=导出,6=导入,7=强退,8=生成代码,9=清空数据,10=其他"`
	BusinessTypes []int64   `json:"businessTypes,omitempty" gorm:"-" query:"businessTypes"`
	Method        string    `json:"method,omitempty" gorm:"column:method" query:"method" xlsx:"n:请求方法"`
	RequestMethod string    `json:"requestMethod,omitempty" gorm:"column:request_method" query:"requestMethod" xlsx:"n:请求方式"`
	OperatorType  *int64    `json:"operatorType,omitempty" gorm:"column:operator_type" query:"operatorType" xlsx:"n:操作类别;exp:0=其它,1=后台用户,2=手机端用户"`
	OperName      string    `json:"operName,omitempty" gorm:"column:oper_name" query:"operName" xlsx:"n:操作人员"`
	DeptName      string    `json:"deptName,omitempty" gorm:"column:dept_name" query:"deptName" xlsx:"n:部门名称"`
	OperUrl       string    `json:"operUrl,omitempty" gorm:"column:oper_url" query:"operUrl" xlsx:"n:请求地址"`
	OperIp        string    `json:"operIp,omitempty" gorm:"column:oper_ip" query:"operIp" xlsx:"n:操作地址"`
	OperParam     string    `json:"operParam,omitempty" gorm:"column:oper_param" query:"operParam" xlsx:"n:请求参数"`
	JsonResult    string    `json:"jsonResult,omitempty" gorm:"column:json_result" query:"jsonResult" xlsx:"n:返回参数"`
	Status        string    `json:"status,omitempty" gorm:"column:status" query:"status" xlsx:"n:状态;exp:0=正常,1=异常"`
	ErrorMsg      string    `json:"errorMsg,omitempty" gorm:"column:error_msg" query:"errorMsg" xlsx:"n:错误消息"`
	OperTime      time.Time `json:"operTime,omitempty" gorm:"column:oper_time" query:"operTime" xlsx:"n:操作时间"`
	// 操作耗时。毫秒为单位
	CostTime int64 `json:"costTime,omitempty" gorm:"column:cost_time" query:"costTime" xlsx:"n:消耗时间;s:毫秒"`
}

func (s *SysOperLog) TableName() string {
	return "sys_oper_log"
}

func (s *SysOperLog) Marshal() string {
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
