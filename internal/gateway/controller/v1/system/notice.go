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
	"strconv"
)

type NoticeController struct{}

func NewNoticeController() *NoticeController {
	return &NoticeController{}
}

type noticeRequestParam struct {
	v12.SysNotice
	api.GetOptions
}

type deleteNoticeParam struct {
	NoticeIds []int64 `json:"noticeIds"`
}

// NoticeList godoc
// @Summary 列出公告列表
// @Description 根据条件查询公告列表
// @Description 权限：system:notice:list
// @Param noticeTitle formData string false "公告标题"
// @Param createBy formData string false "操作人员"
// @Param noticeType formData string false "公告类型"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/notice/list [GET]
func (n *NoticeController) List(ctx context.Context, c *app.RequestContext) {
	var req noticeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListSysNotices(ctx, &v1.ListSysNoticesRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		NoticeInfo: v1.SysNotice2NoticeInfo(&req.SysNotice),
		User:       &info,
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

// GetInfo godoc
// @Summary 获取公告详情
// @Description 获取公告详情信息
// @Description 权限：system:notice:query
// @Param noticeId query string true "公告id"
// @Accept application/json
// @Produce application/json
// @Router /system/notice/:noticeId [GET]
func (n *NoticeController) GetInfo(ctx context.Context, c *app.RequestContext) {
	noticeIdStr := c.Param("noticeId")
	noticeId, _ := strconv.ParseInt(noticeIdStr, 10, 64)

	resp, err := rpc.Remoting.GetSysNoticeById(ctx, noticeId)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.WriteResponse(c, core.ErrResponse{
		Code:    int(resp.BaseResp.Code),
		Message: resp.BaseResp.Msg,
		Data:    resp.Data,
	})
}

// NoticeAdd godoc
// @Summary 创建公告
// @Description 创建公告
// @Description 权限：system:notice:add
// @Param noticeTitle formData string true "标题"
// @Param noticeType formData string true "类型"
// @Param noticeContent formData string true "内容"
// @Accept application/json
// @Produce application/json
// @Router /system/notice [POST]
func (n *NoticeController) Add(ctx context.Context, c *app.RequestContext) {
	var req noticeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysNotice(ctx, &v1.CreateSysNoticeRequest{
		NoticeInfo: v1.SysNotice2NoticeInfo(&req.SysNotice),
		User:       &info,
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

// NoticeEdit godoc
// @Summary 修改公告
// @Description 修改公告
// @Description 权限：system:notice:edit
// @Param noticeTitle formData string true "标题"
// @Param noticeType formData string true "类型"
// @Param noticeContent formData string true "内容"
// @Accept application/json
// @Produce application/json
// @Router /system/notice [PUT]
func (n *NoticeController) Edit(ctx context.Context, c *app.RequestContext) {
	var req noticeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysNotice(ctx, &v1.UpdateSysNoticeRequest{
		NoticeInfo: v1.SysNotice2NoticeInfo(&req.SysNotice),
		User:       &info,
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

// NoticeRemove godoc
// @Summary 删除公告
// @Description 删除公告
// @Description 权限：system:notice:remove
// @Param noticeIds query string false "公告id"
// @Accept application/json
// @Produce application/json
// @Router /system/notice/:noticeIds [DELETE]
func (n *NoticeController) Remove(ctx context.Context, c *app.RequestContext) {
	noticeIdsStr := c.Param("noticeIds")
	noticeIds := strutil.Strs2Int64(noticeIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysNotice(ctx, &v1.DeleteSysNoticeRequest{
		NoticeIds: noticeIds,
		User:      &info,
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
