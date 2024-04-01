package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
)

type OperLogSrc interface {
	// 新增操作日志
	InsertOperLog(ctx context.Context, operLog *v1.SysOperLog, opts *api.CreateOptions) error
	// 查询系统操作日志集合
	SelectOperLogList(ctx context.Context, operLog *v1.SysOperLog, opts *api.GetOptions) *v1.OperLogList
	// 批量删除系统操作日志
	DeleteOperLogByIds(ctx context.Context, operIds []int64, opts *api.DeleteOptions) error
	// 查询操作日志详细
	SelectOperLogById(ctx context.Context, operId int64, opts *api.GetOptions) *v1.SysOperLog
	// 清空操作日志
	CleanOperLog(ctx context.Context, opts *api.DeleteOptions) error
}

type operLogService struct {
	store store.Factory
}

var _ OperLogSrc = &operLogService{}

func NewOperLogs(s store.Factory) OperLogSrc {
	return &operLogService{s}
}

func (s *operLogService) InsertOperLog(ctx context.Context, operLog *v1.SysOperLog, opts *api.CreateOptions) error {
	return s.store.OperLogs().InsertOperLog(ctx, operLog, opts)
}

func (s *operLogService) SelectOperLogList(ctx context.Context, operLog *v1.SysOperLog, opts *api.GetOptions) *v1.OperLogList {
	result, total, err := s.store.OperLogs().SelectOperLogList(ctx, operLog, opts)
	if err != nil {
		return &v1.OperLogList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.OperLogList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *operLogService) DeleteOperLogByIds(ctx context.Context, operIds []int64, opts *api.DeleteOptions) error {
	return s.store.OperLogs().DeleteOperLogByIds(ctx, operIds, opts)
}

func (s *operLogService) SelectOperLogById(ctx context.Context, operId int64, opts *api.GetOptions) *v1.SysOperLog {
	result, err := s.store.OperLogs().SelectOperLogById(ctx, operId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *operLogService) CleanOperLog(ctx context.Context, opts *api.DeleteOptions) error {
	return s.store.OperLogs().CleanOperLog(ctx, opts)
}
