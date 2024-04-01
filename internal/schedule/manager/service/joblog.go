package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/store"
	"github.com/user823/Sophie/pkg/log"
)

type JobLogSrv interface {
	//获取quartz调度器日志的执行任务
	SelectJobLogList(ctx context.Context, jobLog *v1.SysJobLog, opts *api.GetOptions) *v1.JobLogList
	// 通过调度任务日志ID查询调度信息
	SelectJobLogById(ctx context.Context, jobLogId int64, opts *api.GetOptions) *v1.SysJobLog
	// 新增任务日志
	AddJobLog(ctx context.Context, jobLog *v1.SysJobLog, opts *api.CreateOptions) error
	// 批量删除调度日志信息
	DeleteJobLogByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
	// 删除任务日志
	DeleteJobLogById(ctx context.Context, jobId int64, opts *api.DeleteOptions) error
	// 清空任务日志
	CleanJobLog(ctx context.Context)
}

type jobLogService struct {
	store store.Factory
}

func NewJobLogs(s store.Factory) JobLogSrv {
	return &jobLogService{s}
}

func (s *jobLogService) SelectJobLogList(ctx context.Context, jobLog *v1.SysJobLog, opts *api.GetOptions) *v1.JobLogList {
	result, total, _ := s.store.JobLogs().SelectJobLogList(ctx, jobLog, opts)
	return &v1.JobLogList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *jobLogService) SelectJobLogById(ctx context.Context, jobLogId int64, opts *api.GetOptions) *v1.SysJobLog {
	result, err := s.store.JobLogs().SelectJobLogById(ctx, jobLogId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *jobLogService) AddJobLog(ctx context.Context, jobLog *v1.SysJobLog, opts *api.CreateOptions) error {
	return s.store.JobLogs().InsertJobLog(ctx, jobLog, opts)
}

func (s *jobLogService) DeleteJobLogByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	return s.store.JobLogs().DeleteJobLogByIds(ctx, ids, opts)
}

func (s *jobLogService) DeleteJobLogById(ctx context.Context, jobId int64, opts *api.DeleteOptions) error {
	return s.store.JobLogs().DeleteJobLogById(ctx, jobId, opts)
}

func (s *jobLogService) CleanJobLog(ctx context.Context) {
	if err := s.store.JobLogs().CleanJobLog(ctx, &api.DeleteOptions{}); err != nil {
		log.Errorf("清除调用日志失败: %s", err.Error())
	}
}
