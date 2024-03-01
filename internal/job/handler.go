package job

import (
	"context"
	v1 "github.com/user823/Sophie/api/thrift/job/v1"
)

// JobServiceImpl implements the last service interface defined in the IDL.
type JobServiceImpl struct{}

// ListJobs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ListJobs(ctx context.Context, req *v1.ListJobsRequest) (resp *v1.ListJobsResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportJobs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ExportJobs(ctx context.Context, req *v1.ExportJobsRequest) (resp *v1.ExportJobsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetJobInfo implements the JobServiceImpl interface.
func (s *JobServiceImpl) GetJobInfo(ctx context.Context, jobId int64) (resp *v1.GetJobInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateJob implements the JobServiceImpl interface.
func (s *JobServiceImpl) CreateJob(ctx context.Context, req *v1.CreateJobRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateJob implements the JobServiceImpl interface.
func (s *JobServiceImpl) UpdateJob(ctx context.Context, req *v1.UpdateJobRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ChangeStatus implements the JobServiceImpl interface.
func (s *JobServiceImpl) ChangeStatus(ctx context.Context, req *v1.ChangeStatusRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// Run implements the JobServiceImpl interface.
func (s *JobServiceImpl) Run(ctx context.Context, req *v1.RunRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// RemoveJobs implements the JobServiceImpl interface.
func (s *JobServiceImpl) RemoveJobs(ctx context.Context, req *v1.RemoveJobsRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListJobLogs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ListJobLogs(ctx context.Context, req *v1.ListJobLogsRequest) (resp *v1.ListJobLogsResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportJobLogs implements the JobServiceImpl interface.
func (s *JobServiceImpl) ExportJobLogs(ctx context.Context, req *v1.ExportJobLogsRequest) (resp *v1.ExportJobLogsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetJobLogInfo implements the JobServiceImpl interface.
func (s *JobServiceImpl) GetJobLogInfo(ctx context.Context, jobLogId int64) (resp *v1.GetJobLogInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// RemoveJobLogs implements the JobServiceImpl interface.
func (s *JobServiceImpl) RemoveJobLogs(ctx context.Context, req *v1.RemoveJobLogsRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// Clean implements the JobServiceImpl interface.
func (s *JobServiceImpl) Clean(ctx context.Context) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}
