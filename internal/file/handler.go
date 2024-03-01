package file

import (
	"context"
	v1 "github.com/user823/Sophie/api/thrift/file/v1"
)

// FileServiceImpl implements the last service interface defined in the IDL.
type FileServiceImpl struct{}

// Upload implements the FileServiceImpl interface.
func (s *FileServiceImpl) Upload(ctx context.Context, req *v1.UploadRequest) (resp *v1.UploadResponse, err error) {
	// TODO: Your code here...
	return
}
