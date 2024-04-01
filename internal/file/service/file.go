package service

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/file/v1"
	"github.com/user823/Sophie/internal/file/store"
	"github.com/user823/Sophie/internal/pkg/obs"
	"github.com/user823/Sophie/pkg/errors"
)

type FileSrv interface {
	// 文件上传接口
	UploadFile(ctx context.Context, path string, userid int64, file []byte) (*v1.SysFile, error)
}

type fileService struct {
	store store.Factory
}

func NewFileService(store store.Factory) FileSrv {
	return &fileService{store}
}

func (s *fileService) UploadFile(ctx context.Context, path string, userid int64, file []byte) (*v1.SysFile, error) {
	if !obs.CheckFileName(path) {
		return nil, errors.WithMessagef(nil, "文件格式不正确, 仅支持%v 格式", obs.AllowedPrefix)
	}

	objectName := obs.GetNewFileName(path, userid)
	return s.store.Files().Upload(ctx, objectName, file, &api.CreateOptions{})
}
