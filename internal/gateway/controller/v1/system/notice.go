package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type NoticeController struct{}

func NewNoticeController() *NoticeController {
	return &NoticeController{}
}

type noticeRequestParam struct {
	v1.NoticeInfo
	v1.PageInfo
}

type deleteNoticeParam struct {
	NoticeIds []int64 `json:"noticeIds"`
}

func (n *NoticeController) List(ctx context.Context, c *app.RequestContext) {
	var req noticeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListSysNotices(ctx, &v1.ListSysNoticesRequest{
		PageInfo:   &req.PageInfo,
		NoticeInfo: &req.NoticeInfo,
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

func (n *NoticeController) GetInfo(ctx context.Context, c *app.RequestContext) {
	var req noticeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetSysNoticeById(ctx, req.NoticeId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.WriteResponse(c, core.ErrResponse{
		Code:    int(resp.BaseResp.Code),
		Message: resp.BaseResp.Msg,
		Data:    resp.Data,
	})
}

func (n *NoticeController) Add(ctx context.Context, c *app.RequestContext) {
	var req noticeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysNotice(ctx, &v1.CreateSysNoticeRequest{
		NoticeInfo: &req.NoticeInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (n *NoticeController) Edit(ctx context.Context, c *app.RequestContext) {
	var req noticeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysNotice(ctx, &v1.UpdateSysNoticeRequest{
		NoticeInfo: &req.NoticeInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (n *NoticeController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deleteNoticeParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysNotice(ctx, &v1.DeleteSysNoticeRequest{
		NoticeIds: req.NoticeIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
