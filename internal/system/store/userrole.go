package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type UserRoleStore interface {
	// 通过用户id删除用户和角色关联
	DeleteUserRoleByUserId(ctx context.Context, userid int64, opts *api.DeleteOptions) error
	// 批量删除用户和角色关联
	DeleteUserRole(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
	// 通过角色Id查询角色使用数量
	CountUserRoleByRoleId(ctx context.Context, roleid int64, opts *api.GetOptions) int
	// 批量新增用户角色信息
	BatchUserRole(ctx context.Context, userRoleList []*v1.SysUserRole, opts *api.CreateOptions) error
	// 删除用户和角色关联信息
	DeleteUserRoleInfo(ctx context.Context, userRole *v1.SysUserRole, opts *api.DeleteOptions) error
	// 批量取消授权用户角色
	DeleteUserRoleInfos(ctx context.Context, roleid int64, userids []int64, opts *api.DeleteOptions) error
}
