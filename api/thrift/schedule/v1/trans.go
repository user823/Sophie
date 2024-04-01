package v1

import (
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/pkg/utils"
)

func JobInfo2SysJob(job *JobInfo) *v1.SysJob {
	if job == nil {
		return nil
	}

	return &v1.SysJob{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  job.CreateBy,
			CreatedAt: utils.Str2Time(job.CreateTime),
			UpdateBy:  job.UpdateBy,
			UpdatedAt: utils.Str2Time(job.UpdateTime),
			Remark:    job.Remark,
			Extend:    map[string]interface{}{"Params": job.Params},
		},
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

func MJobInfo2SysJob(jobs []*JobInfo) []*v1.SysJob {
	sjobs := make([]*v1.SysJob, 0, len(jobs))
	for i := range jobs {
		sjobs = append(sjobs, JobInfo2SysJob(jobs[i]))
	}
	return sjobs
}

func JobLog2SysJobLog(jobLog *JobLog) *v1.SysJobLog {
	if jobLog == nil {
		return nil
	}

	return &v1.SysJobLog{
		JobLogId:      jobLog.JobLogId,
		JobName:       jobLog.JobName,
		JobGroup:      jobLog.JobGroup,
		InvokeTarget:  jobLog.InvokeTarget,
		JobMessage:    jobLog.JobMessage,
		Status:        jobLog.Status,
		ExceptionInfo: jobLog.ExceptionInfo,
		StartTime:     utils.Str2Time(jobLog.StartTime),
		StopTime:      utils.Str2Time(jobLog.StopTime),
		CreatedAt:     utils.Str2Time(jobLog.CreateTime),
	}
}

func MJobLog2SysJob(jobLogs []*JobLog) []*v1.SysJobLog {
	sjobs := make([]*v1.SysJobLog, 0, len(jobLogs))
	for i := range jobLogs {
		sjobs = append(sjobs, JobLog2SysJobLog(jobLogs[i]))
	}
	return sjobs
}
