package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type LogininforStore interface {
	// 新增系统登录日志
	InsertLogininfor(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.CreateOptions) error
	// 查询系统登录日志集合
	SelectLogininforList(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.GetOptions) ([]*v1.SysLogininfor, error)
	// 批量删除系统登录日志
	DeleteLogininforByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
	// 清空系统登录日志
	CleanLogininfor(ctx context.Context, opts *api.DeleteOptions) error
}
