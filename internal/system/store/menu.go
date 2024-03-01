package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type MenuStore interface {
	// 查询系统菜单列表
	SelectMenuList(ctx context.Context, menu *v1.SysMenu, opts *api.GetOptions) ([]*v1.SysMenu, error)
	// 根据用户所有权限
	SelectMenuPerms(ctx context.Context, opts *api.GetOptions) ([]string, error)
	// 根据用户查询系统菜单列表
	SelectMenuListByUserId(ctx context.Context, menu *v1.SysMenu, userid int64, opts *api.GetOptions) ([]*v1.SysMenu, error)
	// 根据角色id查询权限
	SelectMenuPermsByRoleId(ctx context.Context, roleid int64, opts *api.GetOptions) ([]string, error)
	// 根据用户id查询权限
	SelectMenuPermsByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]string, error)
	// 根据用户id查询菜单
	SelectMenuTreeAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysMenu, error)
	// 根据用户id查询菜单
	SelectMenuTreeByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]*v1.SysMenu, error)
	// 根据角色id查询菜单树信息
	SelectMenuListByRoleId(ctx context.Context, roleid int64, menuCheckStrictly bool, opts *api.GetOptions) ([]int64, error)
	// 根据菜单id查询信息
	SelectMenuById(ctx context.Context, menuid int64, opts *api.GetOptions) (*v1.SysMenu, error)
	// 是否存在菜单子节点
	HasChildByMenuId(ctx context.Context, menuid int64, opts *api.GetOptions) bool
	// 新增菜单信息
	InsertMenu(ctx context.Context, menu *v1.SysMenu, opts *api.CreateOptions) error
	// 修改菜单信息
	UpdateMenu(ctx context.Context, menu *v1.SysMenu, opts *api.UpdateOptions) error
	// 删除菜单管理信息
	DeleteMenuById(ctx context.Context, menuid int64, opts *api.DeleteOptions) error
	// 检验菜单名称是否唯一
	CheckMenuNameUnique(ctx context.Context, menuName string, parentid int64, options *api.GetOptions) *v1.SysMenu
}
