package v1

import "github.com/user823/Sophie/api"

type SysJob struct {
	api.ObjectMeta
	JobId          int64  `json:"jobId" gorm:"column:job_id"`
	JobName        string `json:"jobName" gorm:"column:job_name"`
	JobGroup       string `json:"jobGroup" gorm:"column:job_group"`
	InvokeTarget   string `json:"invokeTarget" gorm:"column:invoke_target"`
	CronExpression string `json:"cronExpression" gorm:"column:cron_expression"`
	MisfirePolicy  string `json:"misfirePolicy" gorm:"column:misfirePolicy"`
	Concurrent     string `json:"concurrent" gorm:"column:concurrent"`
	Status         string `json:"status" gorm:"status"`
}

func (s *SysJob) TableName() string {
	return "sys_job"
}

type JobList struct {
	api.ListMeta `json:",inline"`
	Items        []SysJob `json:"items"`
}
