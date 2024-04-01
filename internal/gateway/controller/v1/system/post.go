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
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strconv"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

type postRequestParam struct {
	api.GetOptions
	v12.SysPost
}

// PostList godoc
// @Summary 岗位列表
// @Description 根据条件查询系统岗位列表
// @Description 权限：system:post:list
// @Param postCode formData string false "岗位编码"
// @Param postName formData string false "岗位名称"
// @Param status formData string false "状态"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/post/list [GET]
func (p *PostController) List(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ListSysPosts(ctx, &v1.ListSysPostsRequest{
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		PostInfo: v1.SysPost2PostInfo(&req.SysPost),
		User:     &info,
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

// PostExport godoc
// @Summary 岗位列表
// @Description 根据条件导出系统岗位列表
// @Description 权限：system:post:export
// @Param postCode formData string false "岗位编码"
// @Param postName formData string false "岗位名称"
// @Param status formData string false "状态"
// @Param pageNum formData int false "页面序号"
// @Param pageSize formData int false "页面大小"
// @Accept application/json
// @Produce application/json
// @Router /system/post/export [POST]
func (p *PostController) Export(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ExportSysPost(ctx, &v1.ExportSysPostRequest{
		User:     &info,
		PostInfo: v1.SysPost2PostInfo(&req.SysPost),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	utils.ExportExcel(c, resp.SheetName, v1.MPostInfo2SysPost(resp.List))
}

// PostInfo godoc
// @Summary 岗位详情
// @Description 查询目标岗位详情
// @Description 权限：system:post:export
// @Param postId query int true "岗位Id"
// @Accept application/json
// @Produce application/json
// @Router /system/post/:postId [GET]
func (p *PostController) GetInfo(ctx context.Context, c *app.RequestContext) {
	postIdStr := c.Param("postId")
	postId, _ := strconv.ParseInt(postIdStr, 10, 64)

	resp, err := rpc.Remoting.GetSysPostById(ctx, postId)
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
		Data:    resp.PostInfo,
	})
}

// PostAdd godoc
// @Summary 新增岗位
// @Description 增加岗位
// @Description 权限：system:post:add
// @Param postName formData string true "岗位名称"
// @Param postCode formData string true "岗位编码"
// @Param postSort formData string true "岗位顺序"
// @Param status formData string false "状态"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/post [POST]
func (p *PostController) Add(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateSysPost(ctx, &v1.CreateSysPostRequest{
		PostInfo: v1.SysPost2PostInfo(&req.SysPost),
		User:     &info,
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

// PostEdit godoc
// @Summary 修改岗位
// @Description 修改岗位
// @Description 权限：system:post:edit
// @Param postName formData string true "岗位名称"
// @Param postCode formData string true "岗位编码"
// @Param postSort formData string true "岗位顺序"
// @Param status formData string false "状态"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/post [PUT]
func (p *PostController) Edit(ctx context.Context, c *app.RequestContext) {
	var req postRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		log.Error(err)
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateSysPost(ctx, &v1.UpdateSysPostRequest{
		PostInfo: v1.SysPost2PostInfo(&req.SysPost),
		User:     &info,
	})
	log.Infof("%v", v1.SysPost2PostInfo(&req.SysPost))
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

// PostDelete godoc
// @Summary 删除岗位
// @Description 删除目标岗位
// @Description 权限：system:post:edit
// @Param postIds query string true "岗位ids"
// @Accept application/json
// @Produce application/json
// @Router /system/:postIds [DELETE]
func (p *PostController) Remove(ctx context.Context, c *app.RequestContext) {
	postIdsStr := c.Param("postIds")
	postIds := strutil.Strs2Int64(postIdsStr)
	log.Info(postIds)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteSysPost(ctx, &v1.DeleteSysPostRequest{
		PostIds: postIds,
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

// OptionSelect godoc
// @Summary 可选岗位
// @Description 列出所有岗位
// @Accept application/json
// @Produce application/json
// @Router /system/optionselect [GET]
func (p *PostController) OptionSelect(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.PostOptionSelect(ctx)
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
