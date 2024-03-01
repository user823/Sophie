package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/file/v1"
)

type FileStore interface {
	// 上传文件服务
	Upload(ctx context.Context, data []byte, opts *api.CreateOptions) (v1.SysFile, error)
}
