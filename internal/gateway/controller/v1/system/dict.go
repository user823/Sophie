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

type DictController struct{}

func NewDictController() *DictController {
	return &DictController{}
}

type dictTypeRequestParam struct {
	api.GetOptions
	v12.SysDictType
}

type deleteTypeParam struct {
	DictIds []int64 `json:"dictIds"`
}

type dictDataRequestParam struct {
	api.GetOptions
	v12.SysDictData
}

type deleteDictDataParam struct {
	DictCodes []int64 `json:"dictCodes"`
}

// DictType Controller

func (d *DictController) ListType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListDictTypes(ctx, &v1.ListDictTypesRequest{
		DictType: v1.SysDictType2DictType(&req.SysDictType),
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
		DateRange: &v1.DateRange{
			BeginTime: req.BeginTime,
			EndTime:   req.EndTime,
		},
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

	result := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	}
	core.JSON(c, result)
}

func (d *DictController) ExportType(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (d *DictController) GetInfoType(ctx context.Context, c *app.RequestContext) {
	dictIdStr := c.Param("dictId")
	dictId, _ := strconv.ParseInt(dictIdStr, 10, 64)

	resp, err := rpc.Remoting.GetDictTypeById(ctx, dictId)
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

func (d *DictController) AddType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.CreateDictType(ctx, &v1.CreateDictTypeRequest{
		DictType: v1.SysDictType2DictType(&req.SysDictType),
		User:     &info,
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

func (d *DictController) EditType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.UpdateDictType(ctx, &v1.UpdateDictTypeRequest{
		DictType: v1.SysDictType2DictType(&req.SysDictType),
		User:     &info,
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

func (d *DictController) RemoveType(ctx context.Context, c *app.RequestContext) {
	dictIdsStr := c.Param("dictIds")
	dictIds := strutil.Strs2Int64(dictIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDictType(ctx, &v1.DeleteDictTypeRequest{
		DictIds: dictIds,
		User:    &info,
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

func (d *DictController) RefreshCache(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.RefreshDictType(ctx)
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

func (d *DictController) OptionSelect(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.DictTypeOptionSelect(ctx)
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

// DictData Controller

func (d *DictController) ListData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.ListDictDatas(ctx, &v1.ListDictDatasRequest{
		DictData: v1.SysDictData2DictData(&req.SysDictData),
		PageInfo: &v1.PageInfo{
			PageNum:       req.PageNum,
			PageSize:      req.PageSize,
			OrderByColumn: req.OrderByColumn,
			IsAsc:         req.QIsAsc,
		},
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

	result := map[string]any{
		"code":  resp.BaseResp.Code,
		"msg":   resp.BaseResp.Msg,
		"total": resp.Total,
		"rows":  resp.Rows,
	}
	core.JSON(c, result)
}

func (d *DictController) ExportData(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (d *DictController) GetInfoData(ctx context.Context, c *app.RequestContext) {
	dictCodeStr := c.Param("dictCode")
	dictCode, _ := strconv.ParseInt(dictCodeStr, 10, 64)

	resp, err := rpc.Remoting.GetDictDataByCode(ctx, dictCode)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.DictData)
}

func (d *DictController) DictType(ctx context.Context, c *app.RequestContext) {
	dictType := c.Param("dictType")

	resp, err := rpc.Remoting.ListDictDataByType(ctx, dictType)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Rows)
}

func (d *DictController) AddData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.CreateDictData(ctx, &v1.CreateDictDataRequest{
		DictData: v1.SysDictData2DictData(&req.SysDictData),
		User:     &info,
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

func (d *DictController) EditData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.UpdateDictData(ctx, &v1.UpdateDictDataRequest{
		DictData: v1.SysDictData2DictData(&req.SysDictData),
		User:     &info,
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

func (d *DictController) RemoveData(ctx context.Context, c *app.RequestContext) {
	dictCodeStr := c.Param("dictCodes")
	dictCodes := strutil.Strs2Int64(dictCodeStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDictData(ctx, &v1.DeleteDictDataRequest{
		DictCodes: dictCodes,
		User:      &info,
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
