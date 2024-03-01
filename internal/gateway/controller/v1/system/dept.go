package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type DeptController struct{}

func NewDeptController() *DeptController {
	return &DeptController{}
}

func (d *DeptController) List(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListDepts(ctx, &v1.ListDeptsRequest{
		DeptInfo: &req,
	})
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (d *DeptController) ExcludeChild(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListDeptsExcludeChild(ctx, req.DeptId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (d *DeptController) GetInfo(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetDeptById(ctx, req.GetDeptId())
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (d *DeptController) Add(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.CreateDept(ctx, &v1.CreateDeptRequest{
		Dept: &req,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (d *DeptController) Edit(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.UpdateDept(ctx, &v1.UpdateDeptRequest{
		Dept: &req,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (d *DeptController) Remove(ctx context.Context, c *app.RequestContext) {
	var req v1.DeptInfo
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDept(ctx, &v1.DeleteDeptRequest{
		DeptId: req.DeptId,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
