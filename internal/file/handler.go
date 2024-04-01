package file

import (
	"context"
	v1 "github.com/user823/Sophie/api/thrift/file/v1"
	"github.com/user823/Sophie/internal/file/service"
	"github.com/user823/Sophie/internal/file/store/obs"
	"github.com/user823/Sophie/internal/file/utils"
)

// FileServiceImpl implements the last service interface defined in the IDL.
type FileServiceImpl struct{}

// Upload implements the FileServiceImpl interface.
func (s *FileServiceImpl) Upload(ctx context.Context, req *v1.UploadRequest) (resp *v1.UploadResponse, err error) {
	store := obs.GetOBSFactoryOr()
	result, err := service.NewFileService(store).UploadFile(ctx, req.Path, req.UserId, req.Data)
	if err != nil {
		return &v1.UploadResponse{
			BaseResp: utils.Fail("系统内部错误"),
		}, nil
	}
	return &v1.UploadResponse{
		BaseResp: utils.Ok("操作成功"),
		File: &v1.FileInfo{
			Name: result.Name,
			Url:  result.Url,
		},
	}, nil
}
