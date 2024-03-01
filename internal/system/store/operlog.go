package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type OperLogStore interface {
	// 新增操作日志
	InsertOperLog(ctx context.Context, operlog *v1.SysOperLog, opts *api.CreateOptions) error
	// 查询系统操作日志集合
	SelectOperLogList(ctx context.Context, operlog *v1.SysOperLog, opts *api.GetOptions) ([]*v1.SysOperLog, error)
	// 批量删除系统操作日志
	DeleteOperLogByIds(ctx context.Context, operids []int64, opts *api.DeleteOptions) error
	// 查询操作日志详情
	SelectOperLogById(ctx context.Context, operid int64, opts *api.GetOptions) (*v1.SysOperLog, error)
	// 清空操作日志
	CleanOperLog(ctx context.Context, opts *api.DeleteOptions) error
}
