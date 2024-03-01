package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/job/v1"
)

type JobLogSrv interface {
	//获取quartz调度器日志的执行任务
	SelectJobLogList(ctx context.Context, jobLog *v1.SysJobLog, opts *api.GetOptions) []v1.SysJobLog
	// 通过调度任务日志ID查询调度信息
	SelectJobLogById(ctx context.Context, jobLogId int64, opts *api.GetOptions) *v1.SysJobLog
	// 新增任务日志
	AddJobLog(ctx context.Context, jobLog *v1.SysJobLog, opts *api.CreateOptions) error
	// 批量删除调度日志信息
	DeleteJobLogByIds(ctx context.Context, jobLog *v1.SysJobLog, opts *api.DeleteOptions) error
	// 删除任务日志
	DeleteJobLogById(ctx context.Context, jobId int64, opts *api.DeleteOptions) error
	// 清空任务日志
	CleanJobLog(ctx context.Context)
}
