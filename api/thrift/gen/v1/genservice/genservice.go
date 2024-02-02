// Code generated by Kitex v0.8.0. DO NOT EDIT.

package genservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
)

func serviceInfo() *kitex.ServiceInfo {
	return genServiceServiceInfo
}

var genServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "GenService"
	handlerType := (*v1.GenService)(nil)
	methods := map[string]kitex.MethodInfo{
		"ListGenTables":   kitex.NewMethodInfo(listGenTablesHandler, newGenServiceListGenTablesArgs, newGenServiceListGenTablesResult, false),
		"GetInfoById":     kitex.NewMethodInfo(getInfoByIdHandler, newGenServiceGetInfoByIdArgs, newGenServiceGetInfoByIdResult, false),
		"DataList":        kitex.NewMethodInfo(dataListHandler, newGenServiceDataListArgs, newGenServiceDataListResult, false),
		"ListColumnsById": kitex.NewMethodInfo(listColumnsByIdHandler, newGenServiceListColumnsByIdArgs, newGenServiceListColumnsByIdResult, false),
		"ImportTableSave": kitex.NewMethodInfo(importTableSaveHandler, newGenServiceImportTableSaveArgs, newGenServiceImportTableSaveResult, false),
		"EditSave":        kitex.NewMethodInfo(editSaveHandler, newGenServiceEditSaveArgs, newGenServiceEditSaveResult, false),
		"DeleteTables":    kitex.NewMethodInfo(deleteTablesHandler, newGenServiceDeleteTablesArgs, newGenServiceDeleteTablesResult, false),
		"PreviewById":     kitex.NewMethodInfo(previewByIdHandler, newGenServicePreviewByIdArgs, newGenServicePreviewByIdResult, false),
		"DownloadByName":  kitex.NewMethodInfo(downloadByNameHandler, newGenServiceDownloadByNameArgs, newGenServiceDownloadByNameResult, false),
		"GenCode":         kitex.NewMethodInfo(genCodeHandler, newGenServiceGenCodeArgs, newGenServiceGenCodeResult, false),
		"SynchDB":         kitex.NewMethodInfo(synchDBHandler, newGenServiceSynchDBArgs, newGenServiceSynchDBResult, false),
		"BatchGenCode":    kitex.NewMethodInfo(batchGenCodeHandler, newGenServiceBatchGenCodeArgs, newGenServiceBatchGenCodeResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "v1",
		"ServiceFilePath": `api/thrift/gen/v1/gen.thrift`,
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.8.0",
		Extra:           extra,
	}
	return svcInfo
}

func listGenTablesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceListGenTablesArgs)
	realResult := result.(*v1.GenServiceListGenTablesResult)
	success, err := handler.(v1.GenService).ListGenTables(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceListGenTablesArgs() interface{} {
	return v1.NewGenServiceListGenTablesArgs()
}

func newGenServiceListGenTablesResult() interface{} {
	return v1.NewGenServiceListGenTablesResult()
}

func getInfoByIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceGetInfoByIdArgs)
	realResult := result.(*v1.GenServiceGetInfoByIdResult)
	success, err := handler.(v1.GenService).GetInfoById(ctx, realArg.Id)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceGetInfoByIdArgs() interface{} {
	return v1.NewGenServiceGetInfoByIdArgs()
}

func newGenServiceGetInfoByIdResult() interface{} {
	return v1.NewGenServiceGetInfoByIdResult()
}

func dataListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceDataListArgs)
	realResult := result.(*v1.GenServiceDataListResult)
	success, err := handler.(v1.GenService).DataList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceDataListArgs() interface{} {
	return v1.NewGenServiceDataListArgs()
}

func newGenServiceDataListResult() interface{} {
	return v1.NewGenServiceDataListResult()
}

func listColumnsByIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceListColumnsByIdArgs)
	realResult := result.(*v1.GenServiceListColumnsByIdResult)
	success, err := handler.(v1.GenService).ListColumnsById(ctx, realArg.Id)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceListColumnsByIdArgs() interface{} {
	return v1.NewGenServiceListColumnsByIdArgs()
}

func newGenServiceListColumnsByIdResult() interface{} {
	return v1.NewGenServiceListColumnsByIdResult()
}

func importTableSaveHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceImportTableSaveArgs)
	realResult := result.(*v1.GenServiceImportTableSaveResult)
	success, err := handler.(v1.GenService).ImportTableSave(ctx, realArg.Tables)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceImportTableSaveArgs() interface{} {
	return v1.NewGenServiceImportTableSaveArgs()
}

func newGenServiceImportTableSaveResult() interface{} {
	return v1.NewGenServiceImportTableSaveResult()
}

func editSaveHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceEditSaveArgs)
	realResult := result.(*v1.GenServiceEditSaveResult)
	success, err := handler.(v1.GenService).EditSave(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceEditSaveArgs() interface{} {
	return v1.NewGenServiceEditSaveArgs()
}

func newGenServiceEditSaveResult() interface{} {
	return v1.NewGenServiceEditSaveResult()
}

func deleteTablesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceDeleteTablesArgs)
	realResult := result.(*v1.GenServiceDeleteTablesResult)
	success, err := handler.(v1.GenService).DeleteTables(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceDeleteTablesArgs() interface{} {
	return v1.NewGenServiceDeleteTablesArgs()
}

func newGenServiceDeleteTablesResult() interface{} {
	return v1.NewGenServiceDeleteTablesResult()
}

func previewByIdHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServicePreviewByIdArgs)
	realResult := result.(*v1.GenServicePreviewByIdResult)
	success, err := handler.(v1.GenService).PreviewById(ctx, realArg.Id)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServicePreviewByIdArgs() interface{} {
	return v1.NewGenServicePreviewByIdArgs()
}

func newGenServicePreviewByIdResult() interface{} {
	return v1.NewGenServicePreviewByIdResult()
}

func downloadByNameHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceDownloadByNameArgs)
	realResult := result.(*v1.GenServiceDownloadByNameResult)
	success, err := handler.(v1.GenService).DownloadByName(ctx, realArg.Name)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceDownloadByNameArgs() interface{} {
	return v1.NewGenServiceDownloadByNameArgs()
}

func newGenServiceDownloadByNameResult() interface{} {
	return v1.NewGenServiceDownloadByNameResult()
}

func genCodeHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceGenCodeArgs)
	realResult := result.(*v1.GenServiceGenCodeResult)
	success, err := handler.(v1.GenService).GenCode(ctx, realArg.Name)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceGenCodeArgs() interface{} {
	return v1.NewGenServiceGenCodeArgs()
}

func newGenServiceGenCodeResult() interface{} {
	return v1.NewGenServiceGenCodeResult()
}

func synchDBHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceSynchDBArgs)
	realResult := result.(*v1.GenServiceSynchDBResult)
	success, err := handler.(v1.GenService).SynchDB(ctx, realArg.Name)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceSynchDBArgs() interface{} {
	return v1.NewGenServiceSynchDBArgs()
}

func newGenServiceSynchDBResult() interface{} {
	return v1.NewGenServiceSynchDBResult()
}

func batchGenCodeHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*v1.GenServiceBatchGenCodeArgs)
	realResult := result.(*v1.GenServiceBatchGenCodeResult)
	success, err := handler.(v1.GenService).BatchGenCode(ctx, realArg.Tables)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newGenServiceBatchGenCodeArgs() interface{} {
	return v1.NewGenServiceBatchGenCodeArgs()
}

func newGenServiceBatchGenCodeResult() interface{} {
	return v1.NewGenServiceBatchGenCodeResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) ListGenTables(ctx context.Context, req *v1.ListGenTablesRequest) (r *v1.ListGenTablesResponse, err error) {
	var _args v1.GenServiceListGenTablesArgs
	_args.Req = req
	var _result v1.GenServiceListGenTablesResult
	if err = p.c.Call(ctx, "ListGenTables", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetInfoById(ctx context.Context, id int64) (r *v1.GetInfoResponse, err error) {
	var _args v1.GenServiceGetInfoByIdArgs
	_args.Id = id
	var _result v1.GenServiceGetInfoByIdResult
	if err = p.c.Call(ctx, "GetInfoById", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DataList(ctx context.Context, req *v1.DataListRequest) (r *v1.ListGenTablesResponse, err error) {
	var _args v1.GenServiceDataListArgs
	_args.Req = req
	var _result v1.GenServiceDataListResult
	if err = p.c.Call(ctx, "DataList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListColumnsById(ctx context.Context, id int64) (r *v1.ListGenTablesResponse, err error) {
	var _args v1.GenServiceListColumnsByIdArgs
	_args.Id = id
	var _result v1.GenServiceListColumnsByIdResult
	if err = p.c.Call(ctx, "ListColumnsById", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ImportTableSave(ctx context.Context, tables string) (r *v1.BaseResp, err error) {
	var _args v1.GenServiceImportTableSaveArgs
	_args.Tables = tables
	var _result v1.GenServiceImportTableSaveResult
	if err = p.c.Call(ctx, "ImportTableSave", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) EditSave(ctx context.Context, req *v1.EditSaveRequest) (r *v1.BaseResp, err error) {
	var _args v1.GenServiceEditSaveArgs
	_args.Req = req
	var _result v1.GenServiceEditSaveResult
	if err = p.c.Call(ctx, "EditSave", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteTables(ctx context.Context, req *v1.DeleteTableRequest) (r *v1.BaseResp, err error) {
	var _args v1.GenServiceDeleteTablesArgs
	_args.Req = req
	var _result v1.GenServiceDeleteTablesResult
	if err = p.c.Call(ctx, "DeleteTables", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) PreviewById(ctx context.Context, id int64) (r *v1.PreviewResponse, err error) {
	var _args v1.GenServicePreviewByIdArgs
	_args.Id = id
	var _result v1.GenServicePreviewByIdResult
	if err = p.c.Call(ctx, "PreviewById", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DownloadByName(ctx context.Context, name string) (r *v1.DownloadResponse, err error) {
	var _args v1.GenServiceDownloadByNameArgs
	_args.Name = name
	var _result v1.GenServiceDownloadByNameResult
	if err = p.c.Call(ctx, "DownloadByName", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GenCode(ctx context.Context, name string) (r *v1.BaseResp, err error) {
	var _args v1.GenServiceGenCodeArgs
	_args.Name = name
	var _result v1.GenServiceGenCodeResult
	if err = p.c.Call(ctx, "GenCode", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SynchDB(ctx context.Context, name string) (r *v1.BaseResp, err error) {
	var _args v1.GenServiceSynchDBArgs
	_args.Name = name
	var _result v1.GenServiceSynchDBResult
	if err = p.c.Call(ctx, "SynchDB", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) BatchGenCode(ctx context.Context, tables string) (r *v1.BatchGenCodeResponse, err error) {
	var _args v1.GenServiceBatchGenCodeArgs
	_args.Tables = tables
	var _result v1.GenServiceBatchGenCodeResult
	if err = p.c.Call(ctx, "BatchGenCode", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}