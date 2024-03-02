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

type LogininfoController struct{}

func NewLogininfoController() *LogininfoController {
	return &LogininfoController{}
}

type logininforRequestParam struct {
	api.GetOptions
	v12.SysLogininfor
}

type deleteLogininforParam struct {
	InfoIds []int64 `json:"infoIds"`
}

type unlockParam struct {
	Username string `json:"userName"`
}

func (l *LogininfoController) List(ctx context.Context, c *app.RequestContext) {
	var req logininforRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListSysLogininfos(ctx, &v1.ListSysLogininfosRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		LoginInfo: v1.SysLogininfor2Logininfor(&req.SysLogininfor),
		DateRange: &v1.DateRange{
			BeginTime: req.BeginTime,
			EndTime:   req.EndTime,
		},
		User: v1.LoginUserTrans(&info),
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

func (l *LogininfoController) Export(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (l *LogininfoController) Remove(ctx context.Context, c *app.RequestContext) {
	infoIdsStr := c.Param("infoIds")
	infoIds := strutil.Strs2Int64(infoIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.RemoveSysLogininfosById(ctx, &v1.RemoveSysLogininfosByIdRequest{
		InfoIds: infoIds,
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

func (l *LogininfoController) Clean(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.LogininfoClean(ctx)
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

func (l *LogininfoController) Unlock(ctx context.Context, c *app.RequestContext) {
	userName := c.Param("userName")

	resp, err := rpc.Remoting.UnlockByUserName(ctx, userName)
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
