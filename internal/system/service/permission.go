package service

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
)

type PermissionSrv interface {
	// 获取角色数据权限
	GetRolePermission(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) []string
	// 获取菜单数据权限
	GetMenuPermission(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) []string
}

type permissionService struct {
	store store.Factory
}

var _ PermissionSrv = &permissionService{}

func NewPermissions(s store.Factory) PermissionSrv {
	return &permissionService{s}
}

func (s *permissionService) GetRolePermission(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) []string {
	if user.IsAdmin() {
		return []string{"admin"}
	}
	return NewRoles(s.store).SelectRolePermissionByUserId(ctx, user.UserId, opts)
}

func (s *permissionService) GetMenuPermission(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) []string {
	if user.IsAdmin() {
		return []string{api.ALL_PERMISSIONS}
	}

	if len(user.Roles) > 0 {
		permSet := hashset.New()
		// 多角色设置permissions属性，以便数据权限匹配权限
		for i := range user.Roles {
			rolePerms, _ := s.store.Menus().SelectMenuPermsByRoleId(ctx, user.Roles[i].RoleId, opts)
			user.Roles[i].Permissions = rolePerms
			for _, perm := range rolePerms {
				permSet.Add(perm)
			}
		}

		res := make([]string, 0, permSet.Size())
		for _, v := range permSet.Values() {
			res = append(res, v.(string))
		}
		return res
	}

	result, _ := s.store.Menus().SelectMenuPermsByUserId(ctx, user.UserId, opts)
	return result
}
