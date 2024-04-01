package worker

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/internal/schedule/utils"
	"github.com/user823/Sophie/internal/schedule/worker/service"
)

// WorkerServiceImpl implements the last service interface defined in the IDL.
type WorkerServiceImpl struct {
	jobMp  models.JobMap
	nodeid string
}

// CreateJob implements the WorkerServiceImpl interface.
func (s *WorkerServiceImpl) CreateJob(ctx context.Context, job *v1.JobInfo) (resp *v1.BaseResp, err error) {
	sysJob := v1.JobInfo2SysJob(job)

	if err = service.NewJobs(s.jobMp, s.nodeid).CreateJob(ctx, sysJob, &api.CreateOptions{Validate: true}); err != nil {
		return utils.Fail(err.Error()), nil
	}
	return utils.Ok("操作成功"), nil
}

// RemoveJobs implements the WorkerServiceImpl interface.
func (s *WorkerServiceImpl) RemoveJobs(ctx context.Context, jobIds []int64) (resp *v1.BaseResp, err error) {
	if err = service.NewJobs(s.jobMp, s.nodeid).RemoveJobs(ctx, jobIds, &api.DeleteOptions{}); err != nil {
		return utils.Fail(err.Error()), nil
	}
	return utils.Ok("操作成功"), nil
}

// PauseJobs implements the WorkerServiceImpl interface.
func (s *WorkerServiceImpl) PauseJobs(ctx context.Context, jobIds []int64) (resp *v1.BaseResp, err error) {
	if err = service.NewJobs(s.jobMp, s.nodeid).PauseJobs(ctx, jobIds); err != nil {
		return utils.Fail(err.Error()), nil
	}
	return utils.Ok("操作成功"), nil
}

// ResumeJobs implements the WorkerServiceImpl interface.
func (s *WorkerServiceImpl) ResumeJobs(ctx context.Context, jobIds []int64) (resp *v1.BaseResp, err error) {
	if err = service.NewJobs(s.jobMp, s.nodeid).ResumeJobs(ctx, jobIds); err != nil {
		return utils.Fail(err.Error()), nil
	}
	return utils.Ok("操作成功"), nil
}

// Run implements the WorkerServiceImpl interface.
func (s *WorkerServiceImpl) Run(ctx context.Context, jobIds []int64) (resp *v1.BaseResp, err error) {
	if err = service.NewJobs(s.jobMp, s.nodeid).Run(ctx, jobIds); err != nil {
		return utils.Fail(err.Error()), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateJob implements the WorkerServiceImpl interface.
func (s *WorkerServiceImpl) UpdateJob(ctx context.Context, job *v1.JobInfo) (resp *v1.BaseResp, err error) {
	sysJob := v1.JobInfo2SysJob(job)

	if err = service.NewJobs(s.jobMp, s.nodeid).UpdateJob(ctx, sysJob, &api.UpdateOptions{Validate: true}); err != nil {
		return utils.Fail(err.Error()), nil
	}
	return utils.Ok("操作成功"), nil
}
