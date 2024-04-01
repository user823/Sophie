package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/validators"
)

type SysJob struct {
	api.ObjectMeta
	JobId          int64  `json:"jobId" gorm:"column:job_id" query:"jobId" xlsx:"n:任务序号"`
	JobName        string `json:"jobName" gorm:"column:job_name" query:"jobName" validate:"required,min=0,max=64" xlsx:"n:任务名称"`
	JobGroup       string `json:"jobGroup" gorm:"column:job_group" query:"jobGroup" xlsx:"n:任务组名"`
	InvokeTarget   string `json:"invokeTarget" gorm:"column:invoke_target" query:"invokeTarget" validate:"required,min=10,max=500" xlsx:"n:调用目标字符串"`
	CronExpression string `json:"cronExpression" gorm:"column:cron_expression" query:"cronExpression" validate:"required,min=10,max=255" xlsx:"n:执行表达式"`
	MisfirePolicy  string `json:"misfirePolicy" gorm:"column:misfire_policy" query:"misfirePolicy" xlsx:"n:计划策略;exp:0=默认,1=立即触发执行,2=触发一次执行,3=不触发立即执行"`
	Concurrent     string `json:"concurrent" gorm:"column:concurrent" query:"concurrent" xlsx:"n:并发执行;exp:0=允许,1=禁止"`
	Status         string `json:"status" gorm:"status" query:"status" xlsx:"n:任务状态;exp:0=正常,1=暂停"`
}

func (s *SysJob) TableName() string {
	return "sys_job"
}

func (s *SysJob) Validate() error {
	vd := validators.GetValidatorOr()
	err := vd.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return buildJobErrMsg(e)
		}
	}
	return nil
}

func buildJobErrMsg(err validator.FieldError) error {
	switch err.StructNamespace() {
	case "SysJob.JobName":
		return validators.BuildErrMsgHelper(err, "任务名称")
	case "SysJob.InvokeTarget":
		return validators.BuildErrMsgHelper(err, "调用目标字符串")
	case "SysJob.CronExpression":
		return validators.BuildErrMsgHelper(err, "cron执行表达式")
	}
	return nil
}

func (s *SysJob) Marshal() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysJob) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type JobList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysJob `json:"items"`
}
