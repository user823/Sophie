package service

import "context"

type FileSrv interface {
	// 文件上传接口
	UploadFile(ctx context.Context, file []byte) string
}
