package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
)

type LogininforSrv interface {
	// 新增系统登录日志
	InsertLogininfor(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.CreateOptions) error
	// 查询系统登录日志集合
	SelectLogininforList(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.GetOptions) *v1.LogininforList
	// 批量删除系统登录日志
	DeleteLogininforByIds(ctx context.Context, inforIds []int64, opts *api.DeleteOptions) error
	// 清空系统登录日志
	CleanLogininfor(ctx context.Context, opts *api.DeleteOptions) error
}

type logininforService struct {
	store store.Factory
}

var _ LogininforSrv = &logininforService{}

func NewLogininfors(s store.Factory) LogininforSrv {
	return &logininforService{s}
}

func (s *logininforService) InsertLogininfor(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.CreateOptions) error {
	return s.store.Logininfors().InsertLogininfor(ctx, logininfor, opts)
}

func (s *logininforService) SelectLogininforList(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.GetOptions) *v1.LogininforList {
	result, total, err := s.store.Logininfors().SelectLogininforList(ctx, logininfor, opts)
	if err != nil {
		return &v1.LogininforList{ListMeta: api.ListMeta{0}}
	}
	return &v1.LogininforList{ListMeta: api.ListMeta{total}, Items: result}
}

func (s *logininforService) DeleteLogininforByIds(ctx context.Context, inforIds []int64, opts *api.DeleteOptions) error {
	return s.store.Logininfors().DeleteLogininforByIds(ctx, inforIds, opts)
}

func (s *logininforService) CleanLogininfor(ctx context.Context, opts *api.DeleteOptions) error {
	return s.store.Logininfors().CleanLogininfor(ctx, opts)
}
