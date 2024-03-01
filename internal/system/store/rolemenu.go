package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type RoleMenuStore interface {
	// 查询菜单使用数量
	CheckMenuExistRole(ctx context.Context, menuid int64, opts *api.GetOptions) int
	// 通过角色id 删除角色和菜单关联
	DeleteRoleMenuByRoleId(ctx context.Context, roleid int64, opts *api.DeleteOptions) error
	// 批量删除角色菜单关联信息
	DeleteRoleMenu(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
	// 批量新增角色菜单信息
	BatchRoleMenu(ctx context.Context, roleMenuList []*v1.SysRoleMenu, opts *api.CreateOptions) error
}
