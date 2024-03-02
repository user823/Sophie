package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/utils/strutil"
)

type OpelogController struct{}

func NewOperlogController() *OpelogController {
	return &OpelogController{}
}

type operLogParam struct {
	v12.SysOperLog
	api.GetOptions
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

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListSysOperLogs(ctx, &v1.ListSysOperLogsRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		OperLog: v1.SysOperLog2OperLog(&req.SysOperLog),
		User:    v1.LoginUserTrans(&info),
	})
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
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
	operIdsStr := c.Param("operIds")
	operIds := strutil.Strs2Int64(operIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysOperLog(ctx, &v1.DeleteSysOperLogRequest{
		OperIds: operIds,
		User:    v1.LoginUserTrans(&info),
	})
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (o *OpelogController) Clean(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.OperLogClean(ctx)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.Code != code.SUCCESS {
		core.Fail(c, resp.Msg, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
