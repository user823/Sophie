package manager

import (
	"context"
	"github.com/user823/Sophie/api"
	v12 "github.com/user823/Sophie/api/domain/schedule/v1"
	v1 "github.com/user823/Sophie/api/thrift/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/manager/loadbalance"
	"github.com/user823/Sophie/internal/schedule/manager/service"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/internal/schedule/store/mysql"
	utils2 "github.com/user823/Sophie/internal/schedule/utils"
)

// JobServiceImpl implements the last service interface defined in the IDL.
type JobServiceImpl struct {
	etcdNodePool *models.EtcdNodePool
	lb           loadbalance.LoadBalancer
}

func (s *JobServiceImpl) buildJobSrv() service.JobSrv {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	return service.NewJobs(store, s.etcdNodePool, s.etcdNodePool, s.lb)
}

// ListJobs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ListJobs(ctx context.Context, req *v1.ListJobsRequest) (resp *v1.ListJobsResponse, err error) {
	getOpts := utils2.BuildGetOption(req.PageInfo, false)
	sysJob := v1.JobInfo2SysJob(req.JobInfo)

	list := s.buildJobSrv().SelectJobList(ctx, sysJob, getOpts)
	return &v1.ListJobsResponse{
		BaseResp: utils2.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysJob2JobInfo(list.Items),
	}, nil
}

// ExportJobs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ExportJobs(ctx context.Context, req *v1.ExportJobsRequest) (resp *v1.ExportJobsResponse, err error) {
	sysjob := v1.JobInfo2SysJob(req.JobInfo)

	list := s.buildJobSrv().SelectJobList(ctx, sysjob, &api.GetOptions{Cache: false})
	return &v1.ExportJobsResponse{
		BaseResp:  utils2.Ok("操作成功"),
		SheetName: "定时任务",
		List:      v1.MSysJob2JobInfo(list.Items),
	}, nil
}

// GetJobInfo implements the JobServiceImpl interface.
func (s *JobServiceImpl) GetJobInfo(ctx context.Context, jobId int64) (resp *v1.GetJobInfoResponse, err error) {

	sysJob := s.buildJobSrv().SelectJobById(ctx, &v12.SysJob{JobId: jobId}, &api.GetOptions{Cache: false})
	return &v1.GetJobInfoResponse{
		BaseResp: utils2.Ok("操作成功"),
		Data:     v1.SysJob2JobInfo(sysJob),
	}, nil
}

// CreateJob implements the JobServiceImpl interface.
func (s *JobServiceImpl) CreateJob(ctx context.Context, req *v1.CreateJobRequest) (resp *v1.BaseResp, err error) {
	sysJob := v1.JobInfo2SysJob(req.GetJobinfo())

	if err = s.buildJobSrv().InsertJob(ctx, sysJob, &api.CreateOptions{Validate: false}); err != nil {
		return utils2.Fail(err.Error()), nil
	}
	return utils2.Ok("操作成功"), nil
}

// UpdateJob implements the JobServiceImpl interface.
func (s *JobServiceImpl) UpdateJob(ctx context.Context, req *v1.UpdateJobRequest) (resp *v1.BaseResp, err error) {
	sysJob := v1.JobInfo2SysJob(req.JobInfo)

	if err = s.buildJobSrv().UpdateJob(ctx, sysJob, &api.UpdateOptions{Validate: false}); err != nil {
		return utils2.Fail(err.Error()), nil
	}
	return utils2.Ok("操作成功"), nil
}

// ChangeStatus implements the JobServiceImpl interface.
func (s *JobServiceImpl) ChangeStatus(ctx context.Context, req *v1.ChangeStatusRequest) (resp *v1.BaseResp, err error) {
	sysJob := v1.JobInfo2SysJob(req.JobInfo)
	jobSrv := s.buildJobSrv()

	newJob := jobSrv.SelectJobById(ctx, sysJob, &api.GetOptions{Cache: false})
	if newJob != nil {
		newJob.Status = sysJob.Status
		if err = jobSrv.ChangeStatus(ctx, newJob, &api.UpdateOptions{}); err != nil {
			return utils2.Fail(err.Error()), nil
		}
	}
	return utils2.Ok("操作成功"), nil
}

// Run implements the JobServiceImpl interface.
func (s *JobServiceImpl) Run(ctx context.Context, req *v1.RunRequest) (resp *v1.BaseResp, err error) {
	sysJob := v1.JobInfo2SysJob(req.JobInfo)

	if s.buildJobSrv().Run(ctx, sysJob) {
		return utils2.Fail("任务不存在或已过期"), nil
	}
	return utils2.Ok("操作成功"), nil
}

// RemoveJobs implements the JobServiceImpl interface.
func (s *JobServiceImpl) RemoveJobs(ctx context.Context, req *v1.RemoveJobsRequest) (resp *v1.BaseResp, err error) {
	if err = s.buildJobSrv().DeleteJobByIds(ctx, req.JobIds, &api.DeleteOptions{}); err != nil {
		return utils2.Fail(err.Error()), nil
	}
	return utils2.Ok("操作成功"), nil
}

// ListJobLogs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ListJobLogs(ctx context.Context, req *v1.ListJobLogsRequest) (resp *v1.ListJobLogsResponse, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	getOpts := utils2.BuildGetOption(req.GetPageInfo(), false)
	sysJobLog := v1.JobLog2SysJobLog(req.JobLog)

	list := service.NewJobLogs(store).SelectJobLogList(ctx, sysJobLog, getOpts)
	return &v1.ListJobLogsResponse{
		BaseResp: utils2.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysJobLog2JobLog(list.Items),
	}, nil
}

// ExportJobLogs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ExportJobLogs(ctx context.Context, req *v1.ExportJobLogsRequest) (resp *v1.ExportJobLogsResponse, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)
	sysJobLog := v1.JobLog2SysJobLog(req.JobLog)

	list := service.NewJobLogs(store).SelectJobLogList(ctx, sysJobLog, &api.GetOptions{Cache: false})
	return &v1.ExportJobLogsResponse{
		BaseResp:  utils2.Ok("操作成功"),
		List:      v1.MSysJobLog2JobLog(list.Items),
		SheetName: "任务调度日志",
	}, nil
}

// GetJobLogInfo implements the JobServiceImpl interface.
func (s *JobServiceImpl) GetJobLogInfo(ctx context.Context, jobLogId int64) (resp *v1.GetJobLogInfoResponse, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)

	sysJob := service.NewJobLogs(store).SelectJobLogById(ctx, jobLogId, &api.GetOptions{Cache: false})
	return &v1.GetJobLogInfoResponse{
		BaseResp: utils2.Ok("操作成功"),
		Data:     v1.SysJobLog2JobLog(sysJob),
	}, nil
}

// RemoveJobLogs implements the JobServiceImpl interface.
func (s *JobServiceImpl) RemoveJobLogs(ctx context.Context, req *v1.RemoveJobLogsRequest) (resp *v1.BaseResp, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)

	if err = service.NewJobLogs(store).DeleteJobLogByIds(ctx, req.GetJobIds(), &api.DeleteOptions{}); err != nil {
		return utils2.Fail(err.Error()), nil
	}
	return utils2.Ok("操作成功"), nil
}

// Clean implements the JobServiceImpl interface.
func (s *JobServiceImpl) Clean(ctx context.Context) (resp *v1.BaseResp, err error) {
	store, _ := mysql.GetMySQLFactoryOr(nil)

	service.NewJobLogs(store).CleanJobLog(ctx)
	return utils2.Ok("操作成功"), nil
}
