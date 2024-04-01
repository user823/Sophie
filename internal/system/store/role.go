package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type RoleStore interface {
	// 根据分页条件查询角色数据
	SelectRoleList(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) ([]*v1.SysRole, int64, error)
	// 根据用户id 查询角色
	SelectRolePermissionByUserId(uctx context.Context, serid int64, opts *api.GetOptions) ([]*v1.SysRole, error)
	// 查询所有角色
	SelectRoleAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysRole, error)
	// 根据用户ID 获取角色选择框列表
	SelectRoleListByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]int64, error)
	// 根据角色id查询角色
	SelectRoleById(ctx context.Context, roleid int64, opts *api.GetOptions) (*v1.SysRole, error)
	// 根据用户名查询角色
	SelectRolesByUserName(ctx context.Context, name string, opts *api.GetOptions) ([]*v1.SysRole, error)
	// 检验角色名称是否唯一
	CheckRoleNameUnique(ctx context.Context, name string, opts *api.GetOptions) *v1.SysRole
	// 校验角色key是否唯一
	CheckRoleKeyUnique(ctx context.Context, key string, opts *api.GetOptions) *v1.SysRole
	// 修改角色信息
	UpdateRole(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error
	// 新增角色信息
	InsertRole(ctx context.Context, role *v1.SysRole, opts *api.CreateOptions) error
	// 通过角色id 删除角色
	DeleteRoleById(ctx context.Context, roleid int64, opts *api.DeleteOptions) error
	// 批量删除角色信息
	DeleteRoleByIds(ctx context.Context, roleids []int64, opts *api.DeleteOptions) error
}
