package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type NoticeStore interface {
	// 查询公告信息
	SelectNoticeById(ctx context.Context, noticeid int64, opts *api.GetOptions) (*v1.SysNotice, error)
	// 查询公告列表
	SelectNoticeList(ctx context.Context, notice *v1.SysNotice, opts *api.GetOptions) ([]*v1.SysNotice, int64, error)
	// 新增公告
	InsertNotice(ctx context.Context, notice *v1.SysNotice, opts *api.CreateOptions) error
	// 修改公告
	UpdateNotice(ctx context.Context, notice *v1.SysNotice, opts *api.UpdateOptions) error
	// 删除公告
	DeleteNoticeById(ctx context.Context, noticeid int64, opts *api.DeleteOptions) error
	// 批量删除公告
	DeleteNoticeByIds(ctx context.Context, noticeids []int64, opts *api.DeleteOptions) error
}
