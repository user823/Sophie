package obs

import (
	"bytes"
	"context"
	"github.com/eleven26/goss/core"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/file/v1"
	"github.com/user823/Sophie/internal/file/store"
	"github.com/user823/Sophie/internal/pkg/obs"
)

type obsFileStore struct {
	storage core.Storage
}

var _ store.FileStore = &obsFileStore{}

func (s *obsFileStore) Upload(ctx context.Context, file string, data []byte, opts *api.CreateOptions) (*v1.SysFile, error) {
	freader := bytes.NewReader(data)
	err := s.storage.Put(file, freader)
	if err != nil {
		return nil, err
	}
	url := obs.GetURL(file)
	return &v1.SysFile{
		Name: file,
		Url:  url,
	}, nil
}
