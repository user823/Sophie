package service

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
)

type PostSrv interface {
	// 查询岗位信息集合
	SelectPostList(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) *v1.PostList
	// 查询所有岗位
	SelectPostAll(ctx context.Context, opts *api.GetOptions) *v1.PostList
	// 通过岗位ID查询岗位信息
	SelectPostById(ctx context.Context, postId int64, opts *api.GetOptions) *v1.SysPost
	// 根据用户ID获取岗位选择框列表
	SelectPostListByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []int64
	// 校验岗位名称
	CheckPostNameUnique(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) bool
	// 校验岗位编码
	CheckPostCodeUnique(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) bool
	// 通过岗位ID查询岗位使用数量
	CountUserPostById(ctx context.Context, postId int64, opts *api.GetOptions) int
	// 删除岗位信息
	DeletePostById(ctx context.Context, postId int64, opts *api.DeleteOptions) error
	// 批量删除岗位信息
	DeletePostByIds(ctx context.Context, postIds []int64, opts *api.DeleteOptions) error
	// 新增保存岗位信息
	InsertPost(ctx context.Context, post *v1.SysPost, opts *api.CreateOptions) error
	// 修改保存岗位信息
	UpdatePost(ctx context.Context, post *v1.SysPost, opts *api.UpdateOptions) error
}

type postService struct {
	store store.Factory
}

var _ PostSrv = &postService{}

func NewPosts(s store.Factory) PostSrv {
	return &postService{s}
}

func (s *postService) SelectPostList(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) *v1.PostList {
	result, err := s.store.Posts().SelectPostList(ctx, post, opts)
	if err != nil {
		return &v1.PostList{
			ListMeta: api.ListMeta{int64(len(result))},
		}
	}
	return &v1.PostList{
		ListMeta: api.ListMeta{int64(len(result))},
		Items:    result,
	}
}

func (s *postService) SelectPostAll(ctx context.Context, opts *api.GetOptions) *v1.PostList {
	result, err := s.store.Posts().SelectPostAll(ctx, opts)
	if err != nil {
		return &v1.PostList{
			ListMeta: api.ListMeta{int64(len(result))},
		}
	}
	return &v1.PostList{
		ListMeta: api.ListMeta{int64(len(result))},
		Items:    result,
	}
}

func (s *postService) SelectPostById(ctx context.Context, postId int64, opts *api.GetOptions) *v1.SysPost {
	result, err := s.store.Posts().SelectPostById(ctx, postId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *postService) SelectPostListByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []int64 {
	result, err := s.store.Posts().SelectPostListByUserId(ctx, userId, opts)
	if err != nil {
		return []int64{}
	}
	return result
}

func (s *postService) CheckPostNameUnique(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) bool {
	result := s.store.Posts().CheckPostNameUnique(ctx, post.PostName, opts)
	if result != nil && result.PostId == post.PostId {
		return false
	}
	return true
}

func (s *postService) CheckPostCodeUnique(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) bool {
	result := s.store.Posts().CheckPostCodeUnique(ctx, post.PostCode, opts)
	if result != nil && result.PostId == post.PostId {
		return false
	}
	return true
}

func (s *postService) CountUserPostById(ctx context.Context, postId int64, opts *api.GetOptions) int {
	return s.store.UserPosts().CountUserPostById(ctx, postId, opts)
}

func (s *postService) DeletePostById(ctx context.Context, postId int64, opts *api.DeleteOptions) error {
	return s.store.Posts().DeletePostById(ctx, postId, opts)
}

func (s *postService) DeletePostByIds(ctx context.Context, postIds []int64, opts *api.DeleteOptions) error {
	for i := range postIds {
		post, err := s.store.Posts().SelectPostById(ctx, postIds[i], &api.GetOptions{Cache: true})
		if err != nil {
			continue
		}
		if s.store.UserPosts().CountUserPostById(ctx, postIds[i], &api.GetOptions{Cache: true}) > 0 {
			return fmt.Errorf("%s 已分配，不能删除", post.PostName)
		}
	}
	return s.store.Posts().DeletePostByIds(ctx, postIds, opts)
}

func (s *postService) InsertPost(ctx context.Context, post *v1.SysPost, opts *api.CreateOptions) error {
	return s.store.Posts().InsertPost(ctx, post, opts)
}

func (s *postService) UpdatePost(ctx context.Context, post *v1.SysPost, opts *api.UpdateOptions) error {
	return s.store.Posts().UpdatePost(ctx, post, opts)
}
