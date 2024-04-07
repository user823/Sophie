package job

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	v12 "github.com/user823/Sophie/api/thrift/schedule/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strconv"
)

type JobController struct{}

func NewJobController() *JobController {
	return &JobController{}
}

type jobRequestParam struct {
	v1.SysJob
	api.GetOptions
}

type jobLogRequestParam struct {
	v1.SysJobLog
	api.GetOptions
}

// JobList godoc
// @Summary 定时任务列表
// @Description 根据条件查询定时任务列表
// @Description 权限：monitor:job:list
// @Param jobName formData string false "定时任务名"
// @Param jobGroup formData string false "定时任务组名"
// @Param status formData string false "任务状态"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/list [GET]
func (j *JobController) List(ctx context.Context, c *app.RequestContext) {
	var req jobRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListJobs(ctx, &v12.ListJobsRequest{
		PageInfo: &v12.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		JobInfo: v12.SysJob2JobInfo(&req.SysJob),
	})

	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	res := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	}
	core.JSON(c, res)
}

// JobExport godoc
// @Summary 导出定时任务列表
// @Description 根据条件导出定时任务列表
// @Description 权限：monitor:job:export
// @Param jobName formData string false "定时任务名"
// @Param jobGroup formData string false "定时任务组名"
// @Param status formData string false "任务状态"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/export [POST]
func (j *JobController) Export(ctx context.Context, c *app.RequestContext) {
	var req jobRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ExportJobs(ctx, &v12.ExportJobsRequest{
		JobInfo: v12.SysJob2JobInfo(&req.SysJob),
	})

	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	utils.ExportExcel(c, resp.SheetName, v12.MJobInfo2SysJob(resp.List))
}

// JobInfo godoc
// @Summary 获取任务详情
// @Description 权限：monitor:job:query
// @Param jobId query int false "定时任务id"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/:jobId [GET]
func (j *JobController) GetInfo(ctx context.Context, c *app.RequestContext) {
	jobIdStr := c.Param("jobId")
	jobId, _ := strconv.ParseInt(jobIdStr, 10, 64)

	resp, err := rpc.Remoting.GetJobInfo(ctx, jobId)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

// JobAdd godoc
// @Summary 新增任务
// @Description 权限：monitor:job:add
// @Param jobName formData string true "定时任务名称"
// @Param jobGroup formData string false "任务分组"
// @Param invokeTarget formData string true "调用方法"
// @Param cronExpression formData string true "cron表达式"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job [PUT]
func (j *JobController) Add(ctx context.Context, c *app.RequestContext) {
	var req jobRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		log.Errorf("请求参数错误: %s", err.Error())
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateJob(ctx, &v12.CreateJobRequest{
		Jobinfo: v12.SysJob2JobInfo(&req.SysJob),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

// JobEdit godoc
// @Summary 修改定时任务
// @Description 权限：monitor:job:edit
// @Param jobName formData string true "定时任务名称"
// @Param jobGroup formData string false "任务分组"
// @Param invokeTarget formData string true "调用方法"
// @Param cronExpression formData string true "cron表达式"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job [PUT]
func (j *JobController) Edit(ctx context.Context, c *app.RequestContext) {
	var req jobRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateJob(ctx, &v12.UpdateJobRequest{
		JobInfo: v12.SysJob2JobInfo(&req.SysJob),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

// ChangeStatus godoc
// @Summary 获取任务详情
// @Description 权限：monitor:job:changeStatus
// @Param jobId formData int true "定时任务id"
// @Param status formData string true "状态"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/changeStatus [PUT]
func (j *JobController) ChangeStatus(ctx context.Context, c *app.RequestContext) {
	var req jobRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ChangeStatus(ctx, &v12.ChangeStatusRequest{
		JobInfo: v12.SysJob2JobInfo(&req.SysJob),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

// Run godoc
// @Summary 运行任务
// @Description 权限：monitor:job:changeStatus
// @Param jobId formData int true "定时任务id"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/run [PUT]
func (j *JobController) Run(ctx context.Context, c *app.RequestContext) {
	var req jobRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.Run(ctx, &v12.RunRequest{
		JobInfo: v12.SysJob2JobInfo(&req.SysJob),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

// Remove godoc
// @Summary 移除任务
// @Description 权限：monitor:job:remove
// @Param jobIds formData string true "定时任务id"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/:jobIds [DELETE]
func (j *JobController) Remove(ctx context.Context, c *app.RequestContext) {
	jobIdsStr := c.Param("jobIds")
	jobIds := strutil.Strs2Int64(jobIdsStr)

	resp, err := rpc.Remoting.RemoveJobs(ctx, &v12.RemoveJobsRequest{
		JobIds: jobIds,
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

// JobLogList godoc
// @Summary 定时任务日志列表
// @Description 根据条件查询定时任务日志列表
// @Description 权限：monitor:job:list
// @Param jobName formData string false "定时任务名"
// @Param jobGroup formData string false "定时任务组名"
// @Param status formData string false "任务状态"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/log/list [GET]
func (k *JobController) ListJobLog(ctx context.Context, c *app.RequestContext) {
	var req jobLogRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListJobLogs(ctx, &v12.ListJobLogsRequest{
		PageInfo: &v12.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		JobLog: v12.SysJobLog2JobLog(&req.SysJobLog),
	})

	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	res := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	}
	core.JSON(c, res)
}

// JobLogExport godoc
// @Summary 导出定时任务日志列表
// @Description 根据条件导出定时任务日志列表
// @Description 权限：monitor:job:export
// @Param jobName formData string false "定时任务名"
// @Param jobGroup formData string false "定时任务组名"
// @Param status formData string false "任务状态"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/log/export [POST]
func (j *JobController) ExportJobLog(ctx context.Context, c *app.RequestContext) {
	var req jobLogRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ExportJobLogs(ctx, &v12.ExportJobLogsRequest{
		JobLog: v12.SysJobLog2JobLog(&req.SysJobLog),
	})

	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	utils.ExportExcel(c, resp.SheetName, v12.MJobLog2SysJob(resp.List))
}

// JobLogInfo godoc
// @Summary 获取任务详情
// @Description 权限：monitor:job:query
// @Param jobLogId query int false "定时任务id"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/log/:jobLogId [GET]
func (j *JobController) JobLogInfo(ctx context.Context, c *app.RequestContext) {
	jobIdStr := c.Param("jobLogId")
	jobId, _ := strconv.ParseInt(jobIdStr, 10, 64)

	resp, err := rpc.Remoting.GetJobLogInfo(ctx, jobId)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

// RemoveJobLogs godoc
// @Summary 移除任务日志
// @Description 权限：monitor:job:remove
// @Param jobLogIds formData string true "定时任务id"
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/log/:jobLogIds [DELETE]
func (j *JobController) RemoveJobLog(ctx context.Context, c *app.RequestContext) {
	jobIdsStr := c.Param("jobLogIds")
	jobLogIds := strutil.Strs2Int64(jobIdsStr)

	resp, err := rpc.Remoting.RemoveJobLogs(ctx, &v12.RemoveJobLogsRequest{JobIds: jobLogIds})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, "ok", nil)
}

// CleanJobLogs godoc
// @Summary 移除任务日志
// @Description 权限：monitor:job:remove
// @Accept application/json
// @Produce application/json
// @Router /schedule/job/log/clean [DELETE]
func (j *JobController) CleanJobLog(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.Clean(ctx)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, "ok", nil)
}
