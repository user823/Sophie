package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type OpelogController struct{}

func NewOperlogController() *OpelogController {
	return &OpelogController{}
}

type operLogParam struct {
	v1.OperLog
	v1.PageInfo
}

type deleteOperLogParam struct {
	OperIds []int64 `json:"operIds"`
}

func (o *OpelogController) List(ctx context.Context, c *app.RequestContext) {
	var req operLogParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListSysOperLogs(ctx, &v1.ListSysOperLogsRequest{
		PageInfo: &req.PageInfo,
		OperLog:  &req.OperLog,
	})
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	result := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	}
	core.JSON(c, result)
}

func (o *OpelogController) Export(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (o *OpelogController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deleteOperLogParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysOperLog(ctx, &v1.DeleteSysOperLogRequest{
		OperIds: req.OperIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (o *OpelogController) Clean(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.OperLogClean(ctx)
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
