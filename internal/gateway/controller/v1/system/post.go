package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

type postRequestParam struct {
	v1.PostInfo
	v1.PageInfo
}

type deletePostParam struct {
	PostIds []int64 `json:"postIds"`
}

func (p *PostController) List(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListSysPosts(ctx, &v1.ListSysPostsRequest{
		PageInfo: &req.PageInfo,
		PostInfo: &req.PostInfo,
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

func (p *PostController) Export(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (p *PostController) GetInfo(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetSysPostById(ctx, req.PostId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.WriteResponse(c, core.ErrResponse{
		Code:    int(resp.BaseResp.Code),
		Message: resp.BaseResp.Msg,
		Data:    resp.PostInfo,
	})
}

func (p *PostController) Add(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysPost(ctx, &v1.CreateSysPostRequest{
		PostInfo: &req.PostInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (p *PostController) Edit(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysPost(ctx, &v1.UpdateSysPostRequest{
		PostInfo: &req.PostInfo,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (p *PostController) Remove(ctx context.Context, c *app.RequestContext) {
	var req deletePostParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysPost(ctx, &v1.DeleteSysPostRequest{
		PostIds: req.PostIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (p *PostController) OptionSelect(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.PostOptionSelect(ctx)
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
