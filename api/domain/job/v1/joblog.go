package v1

import (
	"github.com/user823/Sophie/api"
	"time"
)

type SysJobLog struct {
	JobLogId      int64     `json:"jobLogId" gorm:"column:job_log_id"`
	JobName       string    `json:"jobName" gorm:"column:job_name"`
	JobGroup      string    `json:"jobGroup" gorm:"column:job_group"`
	InvokeTarget  string    `json:"invokeTarget" gorm:"column:invoke_target"`
	JobMessage    string    `json:"jobMessage" gorm:"column:job_message"`
	Status        string    `json:"status" gorm:"column:status"`
	ExceptionInfo string    `json:"exceptionInfo" gorm:"column:exceptionInfo"`
	StartTime     time.Time `json:"startTime" gorm:"-"`
	StopTime      time.Time `json:"stopTime" gorm:"-"`
	CreatedAt     time.Time `json:"createTime" gorm:"column:create_time"`
}

func (s *SysJobLog) TableName() string {
	return "sys_job_log"
}

type JobLogList struct {
	api.ListMeta `json:",inline"`
	Items        []SysJobLog `json:"items"`
}
