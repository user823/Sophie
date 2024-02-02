package store

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/system/v1"
)

type LogininfoStore interface {
	Create(ctx context.Context, logininfo *v1.SysLogininfo, opts api.CreateOptions) error
}
