package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/schedule/v1"
)

type JobLogStore interface {
	// 获取gcron调度器日志的计划任务
	SelectJobLogList(ctx context.Context, jobLog *v1.SysJobLog, opts *api.GetOptions) ([]*v1.SysJobLog, int64, error)
	// 查询所有调度任务日志
	SelectJobLogAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysJobLog, int64, error)
	// 通过调度任务日志ID查询调度信息
	SelectJobLogById(ctx context.Context, jobLogId int64, opts *api.GetOptions) (*v1.SysJobLog, error)
	// 新增任务日志
	InsertJobLog(ctx context.Context, jobLog *v1.SysJobLog, opts *api.CreateOptions) error
	// 批量删除调度日志信息
	DeleteJobLogByIds(ctx context.Context, logIds []int64, opts *api.DeleteOptions) error
	// 删除任务日志
	DeleteJobLogById(ctx context.Context, jobId int64, opts *api.DeleteOptions) error
	// 清空任务日志
	CleanJobLog(ctx context.Context, opts *api.DeleteOptions) error
}
