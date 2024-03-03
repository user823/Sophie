package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/gateway/utils"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"strconv"
)

type DeptController struct{}

func NewDeptController() *DeptController {
	return &DeptController{}
}

func (d *DeptController) List(ctx context.Context, c *app.RequestContext) {
	var req v12.SysDept
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListDepts(ctx, &v1.ListDeptsRequest{
		DeptInfo:  v1.SysDept2DeptInfo(&req),
		LoginUser: &info,
	})
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (d *DeptController) ExcludeChild(ctx context.Context, c *app.RequestContext) {
	deptIdStr := c.Param("deptId")
	deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)

	resp, err := rpc.Remoting.ListDeptsExcludeChild(ctx, deptId)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (d *DeptController) GetInfo(ctx context.Context, c *app.RequestContext) {
	deptIdStr := c.Param("deptId")
	deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.GetDeptById(ctx, &v1.GetDeptByIdReq{
		Id:   deptId,
		User: &info,
	})
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

func (d *DeptController) Add(ctx context.Context, c *app.RequestContext) {
	var req v12.SysDept
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.CreateDept(ctx, &v1.CreateDeptRequest{
		Dept: v1.SysDept2DeptInfo(&req),
		User: &info,
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

func (d *DeptController) Edit(ctx context.Context, c *app.RequestContext) {
	var req v12.SysDept
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.UpdateDept(ctx, &v1.UpdateDeptRequest{
		Dept: v1.SysDept2DeptInfo(&req),
		User: &info,
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

func (d *DeptController) Remove(ctx context.Context, c *app.RequestContext) {
	deptIdStr := c.Param("deptId")
	deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDept(ctx, &v1.DeleteDeptRequest{
		DeptId: deptId,
		User:   &info,
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
