package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type UserPostStore interface {
	// 通过id删除用户和岗位关联
	DeleteUserPostByUserId(ctx context.Context, userid int64, opts *api.DeleteOptions) error
	// 通过岗位id查询岗位使用数量
	CountUserPostById(ctx context.Context, postid int64, opts *api.GetOptions) int
	// 批量删除用户和岗位关联
	DeleteUserPost(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
	// 批量新增用户岗位信息
	BatchUserPost(ctx context.Context, userPostList []*v1.SysUserPost, opts *api.CreateOptions) error
}
