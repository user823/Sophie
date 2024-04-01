package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
)

type NoticeSrv interface {
	// 查询公告信息
	SelectNoticeById(ctx context.Context, noticeId int64, opts *api.GetOptions) *v1.SysNotice
	// 查询公告列表
	SelectNoticeList(ctx context.Context, notice *v1.SysNotice, opts *api.GetOptions) *v1.NoticeList
	// 新增公告
	InsertNotice(ctx context.Context, notice *v1.SysNotice, opts *api.CreateOptions) error
	// 修改公告
	UpdateNotice(ctx context.Context, notice *v1.SysNotice, opts *api.UpdateOptions) error
	// 删除公告
	DeleteNoticeById(ctx context.Context, noticeId int64, opts *api.DeleteOptions) error
	// 批量删除公告信息
	DeleteNoticeByIds(ctx context.Context, noticeIds []int64, opts *api.DeleteOptions) error
}

type noticeService struct {
	store store.Factory
}

var _ NoticeSrv = &noticeService{}

func NewNotices(s store.Factory) NoticeSrv {
	return &noticeService{s}
}

func (s *noticeService) SelectNoticeById(ctx context.Context, noticeId int64, opts *api.GetOptions) *v1.SysNotice {
	result, err := s.store.Notices().SelectNoticeById(ctx, noticeId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *noticeService) SelectNoticeList(ctx context.Context, notice *v1.SysNotice, opts *api.GetOptions) *v1.NoticeList {
	result, total, err := s.store.Notices().SelectNoticeList(ctx, notice, opts)
	if err != nil {
		return &v1.NoticeList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.NoticeList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *noticeService) InsertNotice(ctx context.Context, notice *v1.SysNotice, opts *api.CreateOptions) error {
	return s.store.Notices().InsertNotice(ctx, notice, opts)
}

func (s *noticeService) UpdateNotice(ctx context.Context, notice *v1.SysNotice, opts *api.UpdateOptions) error {
	return s.store.Notices().UpdateNotice(ctx, notice, opts)
}

func (s *noticeService) DeleteNoticeById(ctx context.Context, noticeId int64, opts *api.DeleteOptions) error {
	return s.store.Notices().DeleteNoticeById(ctx, noticeId, opts)
}

func (s *noticeService) DeleteNoticeByIds(ctx context.Context, noticeIds []int64, opts *api.DeleteOptions) error {
	return s.store.Notices().DeleteNoticeByIds(ctx, noticeIds, opts)
}
