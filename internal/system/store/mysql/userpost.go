package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"gorm.io/gorm"
)

type mysqlUserPostStore struct {
	db *gorm.DB
}

var _ store.UserPostStore = &mysqlUserPostStore{}

func (s *mysqlUserPostStore) DeleteUserPostByUserId(ctx context.Context, userid int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("user_id = ?", userid).Delete(&v1.SysUserPost{})
	return del.Error
}

func (s *mysqlUserPostStore) CountUserPostById(ctx context.Context, postid int64, opts *api.GetOptions) int {
	query := s.db.Table("sys_user_post").Where("post_id = ?", postid)
	query = opts.SQLCondition(query, "")

	var result int64
	query.Count(&result)
	return int(result)
}

func (s *mysqlUserPostStore) DeleteUserPost(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := s.db.Table("sys_user_post").Where("user_id in ?", ids).Delete(&v1.SysUserPost{})
	return del.Error
}

func (s *mysqlUserPostStore) BatchUserPost(ctx context.Context, userPostList []*v1.SysUserPost, opts *api.CreateOptions) error {
	for i := range userPostList {
		s.db.Table("sys_user_post").Create(userPostList[i])
	}
	return nil
}
