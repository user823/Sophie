// Code generated by Kitex v0.8.0. DO NOT EDIT.

package jobservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	v1 "github.com/user823/Sophie/api/thrift/schedule/v1"
)

func serviceInfo() *kitex.ServiceInfo {
	return jobServiceServiceInfo
}

var jobServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "JobService"
	handlerType := (*v1.JobService)(nil)
	methods := map[string]kitex.MethodInfo{
		"ListJobs":      kitex.NewMethodInfo(listJobsHandler, newJobServiceListJobsArgs, newJobServiceListJobsResult, false),
		"ExportJobs":    kitex.NewMethodInfo(exportJobsHandler, newJobServiceExportJobsArgs, newJobServiceExportJobsResult, false),
		"GetJobInfo":    kitex.NewMethodInfo(getJobInfoHandler, newJobServiceGetJobInfoArgs, newJobServiceGetJobInfoResult, false),
		"CreateJob":     kitex.NewMethodInfo(createJobHandler, newJobServiceCreateJobArgs, newJobServiceCreateJobResult, false),
		"UpdateJob":     kitex.NewMethodInfo(updateJobHandler, newJobServiceUpdateJobArgs, newJobServiceUpdateJobResult, false),
		"ChangeStatus":  kitex.NewMethodInfo(changeStatusHandler, newJobServiceChangeStatusArgs, newJobServiceChangeStatusResult, false),
		"Run":           kitex.NewMethodInfo(runHandler, newJobServiceRunArgs, newJobServiceRunResult, false),
		"RemoveJobs":    kitex.NewMethodInfo(removeJobsHandler, newJobServiceRemoveJobsArgs, newJobServiceRemoveJobsResult, false),
		"ListJobLogs":   kitex.NewMethodInfo(listJobLogsHandler, newJobServiceListJobLogsArgs, newJobServiceListJobLogsResult, false),
		"ExportJobLogs": kitex.NewMethodInfo(exportJobLogsHandler, newJobServiceExportJobLogsArgs, newJobServiceExportJobLogsResult, false),
		"GetJobLogInfo": kitex.NewMethodInfo(getJobLogInfoHandler, newJobServiceGetJobLogInfoArgs, newJobServiceGetJobLogInfoResult, false),
		"RemoveJobLogs": kitex.NewMethodInfo(removeJobLogsHandler, newJobServiceRemoveJobLogsArgs, newJobServiceRemoveJobLogsResult, false),
		"Clean":         kitex.NewMethodInfo(cleanHandler, newJobServiceCleanArgs, newJobServiceCleanResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "v1",
		"ServiceFilePath": `api/thrift/schedule/v1/job.thrift`,
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

func exportJobsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceExportJobsArgs)
	realResult := result.(*v1.JobServiceExportJobsResult)
	success, err := handler.(v1.JobService).ExportJobs(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceExportJobsArgs() interface{} {
	return v1.NewJobServiceExportJobsArgs()
}

func newJobServiceExportJobsResult() interface{} {
	return v1.NewJobServiceExportJobsResult()
}

func getJobInfoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceGetJobInfoArgs)
	realResult := result.(*v1.JobServiceGetJobInfoResult)
	success, err := handler.(v1.JobService).GetJobInfo(ctx, realArg.JobId)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceGetJobInfoArgs() interface{} {
	return v1.NewJobServiceGetJobInfoArgs()
}

func newJobServiceGetJobInfoResult() interface{} {
	return v1.NewJobServiceGetJobInfoResult()
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

func changeStatusHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceChangeStatusArgs)
	realResult := result.(*v1.JobServiceChangeStatusResult)
	success, err := handler.(v1.JobService).ChangeStatus(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceChangeStatusArgs() interface{} {
	return v1.NewJobServiceChangeStatusArgs()
}

func newJobServiceChangeStatusResult() interface{} {
	return v1.NewJobServiceChangeStatusResult()
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

func removeJobsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceRemoveJobsArgs)
	realResult := result.(*v1.JobServiceRemoveJobsResult)
	success, err := handler.(v1.JobService).RemoveJobs(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceRemoveJobsArgs() interface{} {
	return v1.NewJobServiceRemoveJobsArgs()
}

func newJobServiceRemoveJobsResult() interface{} {
	return v1.NewJobServiceRemoveJobsResult()
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

func exportJobLogsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceExportJobLogsArgs)
	realResult := result.(*v1.JobServiceExportJobLogsResult)
	success, err := handler.(v1.JobService).ExportJobLogs(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceExportJobLogsArgs() interface{} {
	return v1.NewJobServiceExportJobLogsArgs()
}

func newJobServiceExportJobLogsResult() interface{} {
	return v1.NewJobServiceExportJobLogsResult()
}

func getJobLogInfoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceGetJobLogInfoArgs)
	realResult := result.(*v1.JobServiceGetJobLogInfoResult)
	success, err := handler.(v1.JobService).GetJobLogInfo(ctx, realArg.JobLogId)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceGetJobLogInfoArgs() interface{} {
	return v1.NewJobServiceGetJobLogInfoArgs()
}

func newJobServiceGetJobLogInfoResult() interface{} {
	return v1.NewJobServiceGetJobLogInfoResult()
}

func removeJobLogsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.JobServiceRemoveJobLogsArgs)
	realResult := result.(*v1.JobServiceRemoveJobLogsResult)
	success, err := handler.(v1.JobService).RemoveJobLogs(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newJobServiceRemoveJobLogsArgs() interface{} {
	return v1.NewJobServiceRemoveJobLogsArgs()
}

func newJobServiceRemoveJobLogsResult() interface{} {
	return v1.NewJobServiceRemoveJobLogsResult()
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

func (p *kClient) ExportJobs(ctx context.Context, req *v1.ExportJobsRequest) (r *v1.ExportJobsResponse, err error) {
	var _args v1.JobServiceExportJobsArgs
	_args.Req = req
	var _result v1.JobServiceExportJobsResult
	if err = p.c.Call(ctx, "ExportJobs", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetJobInfo(ctx context.Context, jobId int64) (r *v1.GetJobInfoResponse, err error) {
	var _args v1.JobServiceGetJobInfoArgs
	_args.JobId = jobId
	var _result v1.JobServiceGetJobInfoResult
	if err = p.c.Call(ctx, "GetJobInfo", &_args, &_result); err != nil {
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

func (p *kClient) ChangeStatus(ctx context.Context, req *v1.ChangeStatusRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceChangeStatusArgs
	_args.Req = req
	var _result v1.JobServiceChangeStatusResult
	if err = p.c.Call(ctx, "ChangeStatus", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Run(ctx context.Context, req *v1.RunRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceRunArgs
	_args.Req = req
	var _result v1.JobServiceRunResult
	if err = p.c.Call(ctx, "Run", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RemoveJobs(ctx context.Context, req *v1.RemoveJobsRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceRemoveJobsArgs
	_args.Req = req
	var _result v1.JobServiceRemoveJobsResult
	if err = p.c.Call(ctx, "RemoveJobs", &_args, &_result); err != nil {
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

func (p *kClient) ExportJobLogs(ctx context.Context, req *v1.ExportJobLogsRequest) (r *v1.ExportJobLogsResponse, err error) {
	var _args v1.JobServiceExportJobLogsArgs
	_args.Req = req
	var _result v1.JobServiceExportJobLogsResult
	if err = p.c.Call(ctx, "ExportJobLogs", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetJobLogInfo(ctx context.Context, jobLogId int64) (r *v1.GetJobLogInfoResponse, err error) {
	var _args v1.JobServiceGetJobLogInfoArgs
	_args.JobLogId = jobLogId
	var _result v1.JobServiceGetJobLogInfoResult
	if err = p.c.Call(ctx, "GetJobLogInfo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RemoveJobLogs(ctx context.Context, req *v1.RemoveJobLogsRequest) (r *v1.BaseResp, err error) {
	var _args v1.JobServiceRemoveJobLogsArgs
	_args.Req = req
	var _result v1.JobServiceRemoveJobLogsResult
	if err = p.c.Call(ctx, "RemoveJobLogs", &_args, &_result); err != nil {
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
