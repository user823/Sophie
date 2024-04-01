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

// DictType Controller

// DictTypeList godoc
// @Summary 字典类型列表
// @Description 根据条件查询字典类型列表
// @Description 权限：system:dict:list
// @Param dictName formData string false "字典名称"
// @Param dictType formData string false "字典类型"
// @Param status formData string false "状态"
// @Param createTime formData string false "创建时间"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type/list [GET]
func (d *DictController) ListType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
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

// DictTypeExport godoc
// @Summary 字典类型列表
// @Description 根据条件导出字典类型列表
// @Description 权限：system:dict:export
// @Param dictName formData string false "字典名称"
// @Param dictType formData string false "字典类型"
// @Param status formData string false "状态"
// @Param createTime formData string false "创建时间"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type/export [POST]
func (d *DictController) ExportType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ExportDictType(ctx, &v1.ExportDictTypeRequest{
		User:     &info,
		DictType: v1.SysDictType2DictType(&req.SysDictType),
	})

	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}
	utils.ExportExcel(c, resp.SheetName, v1.MDictType2SysDictType(resp.List))
}

// DictTypeInfo godoc
// @Summary 字典类型详情
// @Description 获取目标字典类型详情
// @Description 权限：system:dict:query
// @Param dictId formData int false "字典id"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type/:dictId [GET]
func (d *DictController) GetInfoType(ctx context.Context, c *app.RequestContext) {
	dictIdStr := c.Param("dictId")
	dictId, _ := strconv.ParseInt(dictIdStr, 10, 64)

	resp, err := rpc.Remoting.GetDictTypeById(ctx, dictId)
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

// DictTypeAdd godoc
// @Summary 创建字典类型
// @Description 权限：system:dict:add
// @Param dictName formData string true "字典名称"
// @Param dictType formData string true "字典类型"
// @Param status formData string false "状态"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type [POST]
func (d *DictController) AddType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateDictType(ctx, &v1.CreateDictTypeRequest{
		DictType: v1.SysDictType2DictType(&req.SysDictType),
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

// DictTypeEdit godoc
// @Summary 修改字典类型
// @Description 权限：system:dict:edit
// @Param dictName formData string true "字典名称"
// @Param dictType formData string true "字典类型"
// @Param status formData string false "状态"
// @Param remark formData string false "备注"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type [PUT]
func (d *DictController) EditType(ctx context.Context, c *app.RequestContext) {
	var req dictTypeRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateDictType(ctx, &v1.UpdateDictTypeRequest{
		DictType: v1.SysDictType2DictType(&req.SysDictType),
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

// DictTypeRemove godoc
// @Summary 移除字典类型
// @Description 权限：system:dict:remove
// @Param dictIds query string true "字典id"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type/:dictIds [DELETE]
func (d *DictController) RemoveType(ctx context.Context, c *app.RequestContext) {
	dictIdsStr := c.Param("dictIds")
	dictIds := strutil.Strs2Int64(dictIdsStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDictType(ctx, &v1.DeleteDictTypeRequest{
		DictIds: dictIds,
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

// DictTypeRemove godoc
// @Summary 移除字典类型
// @Description 权限：system:dict:remove
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type/refreshCache [DELETE]
func (d *DictController) RefreshCache(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.RefreshDictType(ctx)
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

// DictTypeOptionSelect godoc
// @Summary 列出可选字典类型
// @Accept application/json
// @Produce application/json
// @Router /system/dict/type/optionselect [GET]
func (d *DictController) OptionSelect(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.Remoting.DictTypeOptionSelect(ctx)
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

// DictData Controller

// DictDataList godoc
// @Summary 列出字典数据列表
// @Description 根据条件查询字典数据列表
// @Description 权限：system:dict:list
// @Param dictName formData string false "字典名称"
// @Param dictLabel formData string false "字典标签"
// @Param status formData string false "状态"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/data/list [GET]
func (d *DictController) ListData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
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

// DictDataExport godoc
// @Summary 导出字典数据列表
// @Description 根据条件导出字典数据列表
// @Description 权限：system:dict:export
// @Param dictName formData string false "字典名称"
// @Param dictLabel formData string false "字典标签"
// @Param status formData string false "状态"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/data/export [POST]
func (d *DictController) ExportData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.ExportDictData(ctx, &v1.ExportDictDataRequest{
		User:     &info,
		DictData: v1.SysDictData2DictData(&req.SysDictData),
	})
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	utils.ExportExcel(c, resp.SheetName, resp.List)
}

// DictDataInfo godoc
// @Summary 字典数据详情
// @Description 根据条件获取字典数据详情
// @Description 权限：system:dict:query
// @Param dictCode query string true "字典编码"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/data/:dictCode [GET]
func (d *DictController) GetInfoData(ctx context.Context, c *app.RequestContext) {
	dictCodeStr := c.Param("dictCode")
	dictCode, _ := strconv.ParseInt(dictCodeStr, 10, 64)

	resp, err := rpc.Remoting.GetDictDataByCode(ctx, dictCode)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.DictData)
}

// DictDataByTye godoc
// @Summary 获取字典数据
// @Description 根据字典类型获取字典数据列表
// @Param dictType query string true "字典类型"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/data/type/:dictType [GET]
func (d *DictController) DictType(ctx context.Context, c *app.RequestContext) {
	dictType := c.Param("dictType")

	resp, err := rpc.Remoting.ListDictDataByType(ctx, dictType)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}
	if resp.BaseResp.Code != code.SUCCESS {
		core.Fail(c, resp.BaseResp.Msg, nil)
		return
	}

	core.OK(c, resp.BaseResp.Msg, resp.Rows)
}

// DictDataAdd godoc
// @Summary 添加字典数据
// @Param dictType formData string true "字典类型"
// @Param dictLabel formData string true "数据标签"
// @Param dictValue formData string true "字典键值"
// @Param cssType formData string false "样式属性"
// @Param sort formData string true "显示排序"
// @Param listType formData string false "回显样式"
// @Param status formData string false "状态"
// @Param remark formData string false "评论"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/data [POST]
func (d *DictController) AddData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.CreateDictData(ctx, &v1.CreateDictDataRequest{
		DictData: v1.SysDictData2DictData(&req.SysDictData),
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

// DictDataEdit godoc
// @Summary 修改字典数据
// @Param dictType formData string true "字典类型"
// @Param dictLabel formData string true "数据标签"
// @Param dictValue formData string true "字典键值"
// @Param cssType formData string false "样式属性"
// @Param sort formData string true "显示排序"
// @Param listType formData string false "回显样式"
// @Param status formData string false "状态"
// @Param remark formData string false "评论"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/data [PUT]
func (d *DictController) EditData(ctx context.Context, c *app.RequestContext) {
	var req dictDataRequestParam
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "请求参数错误", nil)
		return
	}

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.UpdateDictData(ctx, &v1.UpdateDictDataRequest{
		DictData: v1.SysDictData2DictData(&req.SysDictData),
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

// DictDataRemove godoc
// @Summary 移除字典数据
// @Param dictCodes query string false "字典编码"
// @Accept application/json
// @Produce application/json
// @Router /system/dict/data/:dictCodes [DELETE]
func (d *DictController) RemoveData(ctx context.Context, c *app.RequestContext) {
	dictCodeStr := c.Param("dictCodes")
	dictCodes := strutil.Strs2Int64(dictCodeStr)

	info, err := utils.GetLoginInfoFromCtx(c)
	if err != nil {
		core.Fail(c, err.Error(), nil)
		return
	}

	resp, err := rpc.Remoting.DeleteDictData(ctx, &v1.DeleteDictDataRequest{
		DictCodes: dictCodes,
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
