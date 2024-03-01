package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/job/v1"
)

type JobSrv interface {
	// 获取quartz调度器的计划任务
	SelectJobList(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) []v1.SysJob
	// 通过调度任务ID查询调度信息
	SelectJobById(ctx context.Context, job *v1.SysJob, opts *api.GetOptions) *v1.SysJob
	// 暂停任务
	PauseJob(ctx context.Context, job *v1.SysJob) error
	// 恢复任务
	ResumeJob(ctx context.Context, job *v1.SysJob) error
	// 删除任务后，所对应的trigger也将被删除
	DeleteJob(ctx context.Context, job *v1.SysJob, opts *api.DeleteOptions) error
	// 批量删除调度信息
	DeleteJobByIds(ctx context.Context, jobIds []int64, opts *api.DeleteOptions) error
	// 任务调度状态修改
	ChangeStatus(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error
	// 立即执行任务
	Run(ctx context.Context, job *v1.SysJob) bool
	// 新增任务
	InsertJob(ctx context.Context, job *v1.SysJob, opts *api.CreateOptions) error
	// 更新任务
	UpdateJob(ctx context.Context, job *v1.SysJob, opts *api.UpdateOptions) error
	// 校验cron表达式是否有效
	CheckCronExpressionIsValid(cronExpression string) bool
}
