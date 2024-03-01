package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/job/v1"
)

type JobStore interface {
	// 查询调度任务日志集合
	SelectJobList(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) ([]v1.SysJob, error)
	// 查询所有调度任务
	SelectJobAll(ctx context.Context, opts *api.GetOptions) ([]v1.SysJob, error)
	// 通过调度ID查询调度任务信息
	SelectJobById(ctx context.Context, jobId int64, opts *api.GetOptions) (v1.SysJob, error)
	// 批量删除调度任务信息
	DeleteJobByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
	// 修改调度任务信息
	UpdateJob(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error
	// 新增调度任务信息
	InsertJob(ctx context.Context, job *v1.SysJob, opts *api.CreateOptions) error
}
