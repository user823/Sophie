package system

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/core"
)

type DictController struct{}

func NewDictController() *DictController {
	return &DictController{}
}

type dictTypeRequestParam struct {
	v1.DictType
	v1.PageInfo
	v1.DateRange
}

type deleteTypeParam struct {
	DictIds []int64 `json:"dictIds"`
}

type dictDataRequestParam struct {
	v1.DictData
	v1.PageInfo
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

	resp, err := rpc.Remoting.ListDictTypes(ctx, &v1.ListDictTypesRequest{
		DictType:  &req.DictType,
		PageInfo:  &req.PageInfo,
		DateRange: &req.DateRange,
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

func (d *DictController) ExportType(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (d *DictController) GetInfoType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetDictTypeById(ctx, req.DictId)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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

	resp, err := rpc.Remoting.CreateDictType(ctx, &v1.CreateDictTypeRequest{
		DictType: &req.DictType,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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

	resp, err := rpc.Remoting.UpdateDictType(ctx, &v1.UpdateDictTypeRequest{
		DictType: &req.DictType,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (d *DictController) RemoveType(ctx context.Context, c *app.RequestContext) {
	var req deleteTypeParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDictType(ctx, &v1.DeleteDictTypeRequest{
		DictIds: req.DictIds,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (d *DictController) RefreshCache(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.RefreshDictType(ctx)
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (d *DictController) OptionSelect(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.DictTypeOptionSelect(ctx)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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

	resp, err := rpc.Remoting.ListDictDatas(ctx, &v1.ListDictDatasRequest{
		DictData: &req.DictData,
		PageInfo: &req.PageInfo,
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

func (d *DictController) ExportData(ctx context.Context, c *app.RequestContext) {
	// TODO
}

func (d *DictController) GetInfoData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.GetDictDataByCode(ctx, req.DictCode)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.DictData)
}

func (d *DictController) DictType(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.ListDictDataByType(ctx, req.DictType)
	if err = rpc.ParseRpcErr(resp.BaseResp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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

	resp, err := rpc.Remoting.CreateDictData(ctx, &v1.CreateDictDataRequest{
		DictData: &req.DictData,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
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

	resp, err := rpc.Remoting.UpdateDictData(ctx, &v1.UpdateDictDataRequest{
		DictData: &req.DictData,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}

func (d *DictController) RemoveData(ctx context.Context, c *app.RequestContext) {
	var req deleteDictDataParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDictData(ctx, &v1.DeleteDictDataRequest{
		DictCodes: req.DictCodes,
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		core.WriteResponseE(c, err, nil)
		return
	}

	core.OK(c, resp.Msg, nil)
}
