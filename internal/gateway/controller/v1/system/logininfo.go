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

// LogininforList godoc
// @Summary 列出登录信息列表
// @Description 根据条件查询登录信息列表
// @Description 权限：system:logininfor:list
// @Param ipaddr formData string false "登录地址"
// @Param userName formData string false "用户名称"
// @Param status formData string false "状态"
// @Param createTime formData string false "登录时间"
// @Accept application/json
// @Produce application/json
// @Router /system/logininfor/list [GET]
func (l *LogininfoController) List(ctx context.Context, c *app.RequestContext) {
	var req logininforRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
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
		User: &info,
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
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

// LogininforExport godoc
// @Summary 导出登录信息列表
// @Description 根据条件导出公告列表
// @Description 权限：system:logininfor:export
// @Param ipaddr formData string false "登录地址"
// @Param userName formData string false "用户名称"
// @Param status formData string false "状态"
// @Param createTime formData string false "登录时间"
// @Accept application/json
// @Produce application/json
// @Router /system/logininfor/export [POST]
func (l *LogininfoController) Export(ctx context.Context, c *app.RequestContext) {
	var req logininforRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ExportLogininfo(ctx, &v1.ExportLogininfoRequest{
		User:      &info,
		LoginInfo: v1.SysLogininfor2Logininfor(&req.SysLogininfor),
	})

	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	utils.ExportExcel(c, resp.SheetName, v1.MLoginInfo2SysLogininfo(resp.List))
}

// LogininforRemove godoc
// @Summary 删除登录信息列表
// @Description 权限：system:logininfor:remove
// @Param infoIds query string false "记录id"
// @Accept application/json
// @Produce application/json
// @Router /system/logininfor/:infoIds [DELETE]
func (l *LogininfoController) Remove(ctx context.Context, c *app.RequestContext) {
	infoIdsStr := c.Param("infoIds")
	infoIds := strutil.Strs2Int64(infoIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.RemoveSysLogininfosById(ctx, &v1.RemoveSysLogininfosByIdRequest{
		InfoIds: infoIds,
		User:    &info,
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

// LogininforClean godoc
// @Summary 清空登录信息列表
// @Description 权限：system:logininfor:remove
// @Accept application/json
// @Produce application/json
// @Router /system/logininfor/clean [DELETE]
func (l *LogininfoController) Clean(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.LogininfoClean(ctx)
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

// Unlock godoc
// @Summary 解锁用户
// @Description 权限：system:logininfor:unlock
// @Param userName query string false "用户名"
// @Accept application/json
// @Produce application/json
// @Router /system/logininfor/unlock/:userName [GET]
func (l *LogininfoController) Unlock(ctx context.Context, c *app.RequestContext) {
	userName := c.Param("userName")

	resp, err := rpc.Remoting.UnlockByUserName(ctx, userName)
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
