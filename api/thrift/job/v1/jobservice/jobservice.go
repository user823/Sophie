// Code generated by Kitex v0.8.0. DO NOT EDIT.

package jobservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	v1 "github.com/user823/Sophie/api/thrift/job/v1"
)

func serviceInfo() *kitex.ServiceInfo {
	return jobServiceServiceInfo
}

var jobServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "JobService"
	handlerType := (*v1.JobService)(nil)
	methods := map[string]kitex.MethodInfo{
		"ListJobs":          kitex.NewMethodInfo(listJobsHandler, newJobServiceListJobsArgs, newJobServiceListJobsResult, false),
		"Export":            kitex.NewMethodInfo(exportHandler, newJobServiceExportArgs, newJobServiceExportResult, false),
		"GetJobInfoById":    kitex.NewMethodInfo(getJobInfoByIdHandler, newJobServiceGetJobInfoByIdArgs, newJobServiceGetJobInfoByIdResult, false),
		"CreateJob":         kitex.NewMethodInfo(createJobHandler, newJobServiceCreateJobArgs, newJobServiceCreateJobResult, false),
		"UpdateJob":         kitex.NewMethodInfo(updateJobHandler, newJobServiceUpdateJobArgs, newJobServiceUpdateJobResult, false),
		"ChangeJobStatus":   kitex.NewMethodInfo(changeJobStatusHandler, newJobServiceChangeJobStatusArgs, newJobServiceChangeJobStatusResult, false),
		"Run":               kitex.NewMethodInfo(runHandler, newJobServiceRunArgs, newJobServiceRunResult, false),
		"DeleteJob":         kitex.NewMethodInfo(deleteJobHandler, newJobServiceDeleteJobArgs, newJobServiceDeleteJobResult, false),
		"ListJobLogs":       kitex.NewMethodInfo(listJobLogsHandler, newJobServiceListJobLogsArgs, newJobServiceListJobLogsResult, false),
		"ExportJobLog":      kitex.NewMethodInfo(exportJobLogHandler, newJobServiceExportJobLogArgs, newJobServiceExportJobLogResult, false),
		"GetJobLogInfoById": kitex.NewMethodInfo(getJobLogInfoByIdHandler, newJobServiceGetJobLogInfoByIdArgs, newJobServiceGetJobLogInfoByIdResult, false),
		"DeleteJobLog":      kitex.NewMethodInfo(deleteJobLogHandler, newJobServiceDeleteJobLogArgs, newJobServiceDeleteJobLogResult, false),
		"Clean":             kitex.NewMethodInfo(cleanHandler, newJobServiceCleanArgs, newJobServiceCleanResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "v1",
		"ServiceFilePath": `api/thrift/job/v1/job.thrift`,
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.8.0",
		Extra:           extra,
	}
	return svcInfo
}

func listJobsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceListJobsArgs)
	realResult := result.(*v1.JobServiceListJobsResult)
	success, err := handler.(v1.JobService).ListJobs(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceListJobsArgs() interface{} {
	return v1.NewJobServiceListJobsArgs()
}

func newJobServiceListJobsResult() interface{} {
	return v1.NewJobServiceListJobsResult()
}

func exportHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceExportArgs)
	realResult := result.(*v1.JobServiceExportResult)
	success, err := handler.(v1.JobService).Export(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceExportArgs() interface{} {
	return v1.NewJobServiceExportArgs()
}

func newJobServiceExportResult() interface{} {
	return v1.NewJobServiceExportResult()
}

func getJobInfoByIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceGetJobInfoByIdArgs)
	realResult := result.(*v1.JobServiceGetJobInfoByIdResult)
	success, err := handler.(v1.JobService).GetJobInfoById(ctx, realArg.Id)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceGetJobInfoByIdArgs() interface{} {
	return v1.NewJobServiceGetJobInfoByIdArgs()
}

func newJobServiceGetJobInfoByIdResult() interface{} {
	return v1.NewJobServiceGetJobInfoByIdResult()
}

func createJobHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceCreateJobArgs)
	realResult := result.(*v1.JobServiceCreateJobResult)
	success, err := handler.(v1.JobService).CreateJob(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceCreateJobArgs() interface{} {
	return v1.NewJobServiceCreateJobArgs()
}

func newJobServiceCreateJobResult() interface{} {
	return v1.NewJobServiceCreateJobResult()
}

func updateJobHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceUpdateJobArgs)
	realResult := result.(*v1.JobServiceUpdateJobResult)
	success, err := handler.(v1.JobService).UpdateJob(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceUpdateJobArgs() interface{} {
	return v1.NewJobServiceUpdateJobArgs()
}

func newJobServiceUpdateJobResult() interface{} {
	return v1.NewJobServiceUpdateJobResult()
}

func changeJobStatusHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceChangeJobStatusArgs)
	realResult := result.(*v1.JobServiceChangeJobStatusResult)
	success, err := handler.(v1.JobService).ChangeJobStatus(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceChangeJobStatusArgs() interface{} {
	return v1.NewJobServiceChangeJobStatusArgs()
}

func newJobServiceChangeJobStatusResult() interface{} {
	return v1.NewJobServiceChangeJobStatusResult()
}

func runHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceRunArgs)
	realResult := result.(*v1.JobServiceRunResult)
	success, err := handler.(v1.JobService).Run(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceRunArgs() interface{} {
	return v1.NewJobServiceRunArgs()
}

func newJobServiceRunResult() interface{} {
	return v1.NewJobServiceRunResult()
}

func deleteJobHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceDeleteJobArgs)
	realResult := result.(*v1.JobServiceDeleteJobResult)
	success, err := handler.(v1.JobService).DeleteJob(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceDeleteJobArgs() interface{} {
	return v1.NewJobServiceDeleteJobArgs()
}

func newJobServiceDeleteJobResult() interface{} {
	return v1.NewJobServiceDeleteJobResult()
}

func listJobLogsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceListJobLogsArgs)
	realResult := result.(*v1.JobServiceListJobLogsResult)
	success, err := handler.(v1.JobService).ListJobLogs(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceListJobLogsArgs() interface{} {
	return v1.NewJobServiceListJobLogsArgs()
}

func newJobServiceListJobLogsResult() interface{} {
	return v1.NewJobServiceListJobLogsResult()
}

func exportJobLogHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceExportJobLogArgs)
	realResult := result.(*v1.JobServiceExportJobLogResult)
	success, err := handler.(v1.JobService).ExportJobLog(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceExportJobLogArgs() interface{} {
	return v1.NewJobServiceExportJobLogArgs()
}

func newJobServiceExportJobLogResult() interface{} {
	return v1.NewJobServiceExportJobLogResult()
}

func getJobLogInfoByIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceGetJobLogInfoByIdArgs)
	realResult := result.(*v1.JobServiceGetJobLogInfoByIdResult)
	success, err := handler.(v1.JobService).GetJobLogInfoById(ctx, realArg.Id)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceGetJobLogInfoByIdArgs() interface{} {
	return v1.NewJobServiceGetJobLogInfoByIdArgs()
}

func newJobServiceGetJobLogInfoByIdResult() interface{} {
	return v1.NewJobServiceGetJobLogInfoByIdResult()
}

func deleteJobLogHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceDeleteJobLogArgs)
	realResult := result.(*v1.JobServiceDeleteJobLogResult)
	success, err := handler.(v1.JobService).DeleteJobLog(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceDeleteJobLogArgs() interface{} {
	return v1.NewJobServiceDeleteJobLogArgs()
}

func newJobServiceDeleteJobLogResult() interface{} {
	return v1.NewJobServiceDeleteJobLogResult()
}

func cleanHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {

	realResult := result.(*v1.JobServiceCleanResult)
	success, err := handler.(v1.JobService).Clean(ctx)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceCleanArgs() interface{} {
	return v1.NewJobServiceCleanArgs()
}

func newJobServiceCleanResult() interface{} {
	return v1.NewJobServiceCleanResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) ListJobs(ctx context.Context, req *v1.ListJobsRequest) (r *v1.ListJobsResponse, err error) {
	var _args v1.JobServiceListJobsArgs
	_args.Req = req
	var _result v1.JobServiceListJobsResult
	if err = p.c.Call(ctx, "ListJobs", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Export(ctx context.Context, req *v1.ExportJobRequest) (r *v1.ExportJobResponse, err error) {
	var _args v1.JobServiceExportArgs
	_args.Req = req
	var _result v1.JobServiceExportResult
	if err = p.c.Call(ctx, "Export", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetJobInfoById(ctx context.Context, id int64) (r *v1.JobInfoResponse, err error) {
	var _args v1.JobServiceGetJobInfoByIdArgs
	_args.Id = id
	var _result v1.JobServiceGetJobInfoByIdResult
	if err = p.c.Call(ctx, "GetJobInfoById", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CreateJob(ctx context.Context, req *v1.CreateJobRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceCreateJobArgs
	_args.Req = req
	var _result v1.JobServiceCreateJobResult
	if err = p.c.Call(ctx, "CreateJob", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateJob(ctx context.Context, req *v1.UpdateJobRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceUpdateJobArgs
	_args.Req = req
	var _result v1.JobServiceUpdateJobResult
	if err = p.c.Call(ctx, "UpdateJob", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ChangeJobStatus(ctx context.Context, req *v1.ChangeJobStatusRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceChangeJobStatusArgs
	_args.Req = req
	var _result v1.JobServiceChangeJobStatusResult
	if err = p.c.Call(ctx, "ChangeJobStatus", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Run(ctx context.Context, req *v1.RunJobRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceRunArgs
	_args.Req = req
	var _result v1.JobServiceRunResult
	if err = p.c.Call(ctx, "Run", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteJob(ctx context.Context, req *v1.DeleteJobRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceDeleteJobArgs
	_args.Req = req
	var _result v1.JobServiceDeleteJobResult
	if err = p.c.Call(ctx, "DeleteJob", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListJobLogs(ctx context.Context, req *v1.ListJobLogsRequest) (r *v1.ListJobLogsResponse, err error) {
	var _args v1.JobServiceListJobLogsArgs
	_args.Req = req
	var _result v1.JobServiceListJobLogsResult
	if err = p.c.Call(ctx, "ListJobLogs", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ExportJobLog(ctx context.Context, req *v1.ExportJobLogRequest) (r *v1.ExportJobLogResponse, err error) {
	var _args v1.JobServiceExportJobLogArgs
	_args.Req = req
	var _result v1.JobServiceExportJobLogResult
	if err = p.c.Call(ctx, "ExportJobLog", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetJobLogInfoById(ctx context.Context, id int64) (r *v1.JobLogInfoResponse, err error) {
	var _args v1.JobServiceGetJobLogInfoByIdArgs
	_args.Id = id
	var _result v1.JobServiceGetJobLogInfoByIdResult
	if err = p.c.Call(ctx, "GetJobLogInfoById", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteJobLog(ctx context.Context, req *v1.DeleteJobLogRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceDeleteJobLogArgs
	_args.Req = req
	var _result v1.JobServiceDeleteJobLogResult
	if err = p.c.Call(ctx, "DeleteJobLog", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Clean(ctx context.Context) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceCleanArgs
	var _result v1.JobServiceCleanResult
	if err = p.c.Call(ctx, "Clean", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
