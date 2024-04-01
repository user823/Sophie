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

// DeptList godoc
// @Summary 部门列表
// @Description 根据条件查询部门列表
// @Description 权限：system:dept:list
// @Param deptName formData string false "部门名称"
// @Param status formData string false "状态"
// @Accept application/json
// @Produce application/json
// @Router /system/dept [GET]
func (d *DeptController) List(ctx context.Context, c *app.RequestContext) {
	var req v12.SysDept
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListDepts(ctx, &v1.ListDeptsRequest{
		DeptInfo:  v1.SysDept2DeptInfo(&req),
		LoginUser: &info,
	})
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

// ExcludeChild godoc
// @Summary 排除子部门
// @Description 权限：system:dept:list
// @Param deptId query string true "部门id"
// @Accept application/json
// @Produce application/json
// @Router /system/dept/list/exclude/:deptId [GET]
func (d *DeptController) ExcludeChild(ctx context.Context, c *app.RequestContext) {
	deptIdStr := c.Param("deptId")
	deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)

	resp, err := rpc.Remoting.ListDeptsExcludeChild(ctx, deptId)
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

// DeptInfo godoc
// @Summary 部门详情
// @Description 权限：system:dept:query
// @Param deptId query string true "部门id"
// @Accept application/json
// @Produce application/json
// @Router /system/dept/:deptId [GET]
func (d *DeptController) GetInfo(ctx context.Context, c *app.RequestContext) {
	deptIdStr := c.Param("deptId")
	deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.GetDeptById(ctx, &v1.GetDeptByIdReq{
		Id:   deptId,
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

	core.OK(c, resp.BaseResp.Msg, resp.Data)
}

// DeptAdd godoc
// @Summary 添加部门
// @Description 权限：system:dept:add
// @Param parentName formData string true "上级部门"
// @Param deptName formData string true "部门名称"
// @Param leader formData string false "负责人"
// @Param email formData string false "部门邮箱"
// @Param status formData string true "状态"
// @Accept application/json
// @Produce application/json
// @Router /system/dept [POST]
func (d *DeptController) Add(ctx context.Context, c *app.RequestContext) {
	var req v12.SysDept
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateDept(ctx, &v1.CreateDeptRequest{
		Dept: v1.SysDept2DeptInfo(&req),
		User: &info,
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

// DeptEdit godoc
// @Summary 添加部门
// @Description 权限：system:dept:edit
// @Param parentName formData string true "上级部门"
// @Param deptName formData string true "部门名称"
// @Param leader formData string false "负责人"
// @Param email formData string false "部门邮箱"
// @Param status formData string true "状态"
// @Accept application/json
// @Produce application/json
// @Router /system/dept [PUT]
func (d *DeptController) Edit(ctx context.Context, c *app.RequestContext) {
	var req v12.SysDept
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateDept(ctx, &v1.UpdateDeptRequest{
		Dept: v1.SysDept2DeptInfo(&req),
		User: &info,
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

// DeptRemove godoc
// @Summary 移除部门
// @Description 权限：system:dept:remove
// @Param deptId query string true "部门id"
// @Accept application/json
// @Produce application/json
// @Router /system/dept/:deptId [DELETE]
func (d *DeptController) Remove(ctx context.Context, c *app.RequestContext) {
	deptIdStr := c.Param("deptId")
	deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDept(ctx, &v1.DeleteDeptRequest{
		DeptId: deptId,
		User:   &info,
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
