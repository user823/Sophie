package v1

import (
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/pkg/utils"
)

func SysJob2JobInfo(job *v1.SysJob) *JobInfo {
	if job == nil {
		return nil
	}

	var params map[string]string
	if data, ok := job.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &JobInfo{
		CreateBy:       job.CreateBy,
		CreateTime:     utils.Time2Str(job.CreatedAt),
		UpdateBy:       job.UpdateBy,
		UpdateTime:     utils.Time2Str(job.UpdatedAt),
		Remark:         job.Remark,
		Params:         params,
		JobId:          job.JobId,
		JobName:        job.JobName,
		JobGroup:       job.JobGroup,
		InvokeTarget:   job.InvokeTarget,
		CronExpression: job.CronExpression,
		MisfirePolicy:  job.MisfirePolicy,
		Concurrent:     job.Concurrent,
		Status:         job.Status,
	}
}

func MSysJob2JobInfo(jobs []*v1.SysJob) []*JobInfo {
	res := make([]*JobInfo, 0, len(jobs))
	for i := range jobs {
		res = append(res, SysJob2JobInfo(jobs[i]))
	}
	return res
}

func SysJobLog2JobLog(jobLog *v1.SysJobLog) *JobLog {
	return &JobLog{
		CreateTime:    utils.Time2Str(jobLog.CreatedAt),
		JobLogId:      jobLog.JobLogId,
		JobName:       jobLog.JobName,
		JobGroup:      jobLog.JobGroup,
		InvokeTarget:  jobLog.InvokeTarget,
		Status:        jobLog.Status,
		JobMessage:    jobLog.JobMessage,
		ExceptionInfo: jobLog.ExceptionInfo,
		StartTime:     utils.Time2Str(jobLog.StartTime),
		StopTime:      utils.Time2Str(jobLog.StopTime),
	}
}

func MSysJobLog2JobLog(joblogs []*v1.SysJobLog) []*JobLog {
	res := make([]*JobLog, 0, len(joblogs))
	for i := range joblogs {
		res = append(res, SysJobLog2JobLog(joblogs[i]))
	}
	return res
}
