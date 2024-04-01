package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type UserStore interface {
	// 根据条件分页查询用户列表
	SelectUserList(ctx context.Context, sysUser *v1.SysUser, opts *api.GetOptions) ([]*v1.SysUser, int64, error)
	// 根据条件分页查询已配用户角色列表
	SelectAllocatedList(ctx context.Context, sysUser *v1.SysUser, opts *api.GetOptions) ([]*v1.SysUser, error)
	// 根据条件分页查询未分配用户角色列表
	SelectUnallocatedList(ctx context.Context, sysUser *v1.SysUser, opts *api.GetOptions) ([]*v1.SysUser, error)
	// 通过用户名查询用户
	SelectUserByUserName(ctx context.Context, name string, opts *api.GetOptions) (*v1.SysUser, error)
	// 通过用户ID 查询用户
	SelectUserById(ctx context.Context, userid int64, opts *api.GetOptions) (*v1.SysUser, error)
	// 新增用户信息
	InsertUser(ctx context.Context, sysUser *v1.SysUser, opts *api.CreateOptions) error
	// 修改用户信息
	UpdateUser(ctx context.Context, sysUser *v1.SysUser, opts *api.UpdateOptions) error
	// 修改用户头像
	UpdateUserAvatar(ctx context.Context, userName, avatar string, opts *api.UpdateOptions) error
	// 重置用户密码
	UpdateUserPwd(ctx context.Context, userName, password string, opts *api.UpdateOptions) error
	// 删除用户
	DeleteUserById(ctx context.Context, userid int64, opts *api.DeleteOptions) error
	// 批量删除用户
	DeleteUserByIds(ctx context.Context, userids []int64, opts *api.DeleteOptions) error
	// 校验用户名称是否唯一
	CheckUserNameUnique(ctx context.Context, name string, opts *api.GetOptions) *v1.SysUser
	// 检验手机号是否唯一
	CheckPhoneUnique(ctx context.Context, phonenumber string, opts *api.GetOptions) *v1.SysUser
	// 校验email 是否唯一
	CheckEmailUnique(ctx context.Context, email string, opts *api.GetOptions) *v1.SysUser
}
