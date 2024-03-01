package gen

import (
	"context"
	v1 "github.com/user823/Sophie/api/thrift/gen/v1"
)

// GenServiceImpl implements the last service interface defined in the IDL.
type GenServiceImpl struct{}

// ListGenTables implements the GenServiceImpl interface.
func (s *GenServiceImpl) ListGenTables(ctx context.Context, req *v1.ListGenTablesRequest) (resp *v1.ListGenTablesResponse, err error) {
	// TODO: Your code here...
	return
}

// GetInfo implements the GenServiceImpl interface.
func (s *GenServiceImpl) GetInfo(ctx context.Context, tableId int64) (resp *v1.GetInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// DataList implements the GenServiceImpl interface.
func (s *GenServiceImpl) DataList(ctx context.Context, req *v1.DataListRequest) (resp *v1.ListGenTablesResponse, err error) {
	// TODO: Your code here...
	return
}

// ColumnList implements the GenServiceImpl interface.
func (s *GenServiceImpl) ColumnList(ctx context.Context, tableId int64) (resp *v1.ListGenTablesResponse, err error) {
	// TODO: Your code here...
	return
}

// ImportTableSave implements the GenServiceImpl interface.
func (s *GenServiceImpl) ImportTableSave(ctx context.Context, tables string) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// EditSave implements the GenServiceImpl interface.
func (s *GenServiceImpl) EditSave(ctx context.Context, req *v1.EditSaveRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// Remove implements the GenServiceImpl interface.
func (s *GenServiceImpl) Remove(ctx context.Context, req *v1.RemoveRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// Preview implements the GenServiceImpl interface.
func (s *GenServiceImpl) Preview(ctx context.Context, tableId int64) (resp *v1.PreviewResponse, err error) {
	// TODO: Your code here...
	return
}

// Download implements the GenServiceImpl interface.
func (s *GenServiceImpl) Download(ctx context.Context, tableName string) (resp *v1.DownloadResponse, err error) {
	// TODO: Your code here...
	return
}

// GenCode implements the GenServiceImpl interface.
func (s *GenServiceImpl) GenCode(ctx context.Context, tableName string) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// SynchDb implements the GenServiceImpl interface.
func (s *GenServiceImpl) SynchDb(ctx context.Context, tableName string) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// BatchGenCode implements the GenServiceImpl interface.
func (s *GenServiceImpl) BatchGenCode(ctx context.Context, tables string) (resp *v1.BatchGenCodeResponse, err error) {
	// TODO: Your code here...
	return
}
