package service

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils"
	"strings"
)

type RoleSrv interface {
	// 根据条件分页查询角色数据
	SelectRoleList(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) *v1.RoleList
	// 根据用户ID查询角色列表
	SelectRolesByUserId(ctx context.Context, userId int64, opts *api.GetOptions) *v1.RoleList
	// 根据用户ID查询角色权限
	SelectRolePermissionByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []string
	// 查询所有角色
	SelectRoleAll(ctx context.Context, opts *api.GetOptions) *v1.RoleList
	// 根据用户ID获取角色选择框列表
	SelectRoleListByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []int64
	// 通过角色ID查询角色
	SelectRoleById(ctx context.Context, roleId int64, opts *api.GetOptions) *v1.SysRole
	// 校验角色名称是否唯一
	CheckRoleNameUnique(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) bool
	// 校验角色权限是否唯一
	CheckRoleKeyUnique(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) bool
	// 校验角色是否允许操作
	CheckRoleAllowed(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) bool
	// 校验角色是否有数据权限
	CheckRoleDataScope(ctx context.Context, roleId int64, opts *api.GetOptions) bool
	// 通过角色ID查询角色使用数量
	CountUserRoleByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) int
	// 新增保存角色信息
	InsertRole(ctx context.Context, role *v1.SysRole, opts *api.CreateOptions) error
	// 修改保存角色信息
	UpdateRole(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error
	// 修改角色状态
	UpdateRoleStatus(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error
	// 修改数据权限信息
	AuthDataScope(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error
	// 通过角色id删除角色
	DeleteRoleById(ctx context.Context, roleId int64, opts *api.DeleteOptions) error
	// 批量删除角色信息
	DeleteRoleByIds(ctx context.Context, roleIds []int64, opts *api.DeleteOptions) error
	// 取消授权用户角色
	DeleteAuthUser(ctx context.Context, userRole *v1.SysUserRole, opts *api.DeleteOptions) error
	// 批量取消用户授权角色
	DeleteAuthUsers(ctx context.Context, roleId int64, userIds []int64, opts *api.DeleteOptions) error
	// 批量选择授权用户角色
	InsertAuthUsers(ctx context.Context, roleId int64, userIds []int64, opts *api.CreateOptions) error
}

type roleService struct {
	store store.Factory
}

var _ RoleSrv = &roleService{}

var (
	ErrRoleNotAllowed = fmt.Errorf("不允许操作超级管理员角色")
	ErrRoleDataScope  = fmt.Errorf("没有权限访问角色数据")
)

func NewRoles(s store.Factory) RoleSrv {
	return &roleService{s}
}

func (s *roleService) SelectRoleList(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) *v1.RoleList {
	result, total, err := s.store.Roles().SelectRoleList(ctx, role, opts)
	if err != nil {
		return &v1.RoleList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.RoleList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *roleService) SelectRolesByUserId(ctx context.Context, userId int64, opts *api.GetOptions) *v1.RoleList {
	userRoles, _ := s.store.Roles().SelectRolePermissionByUserId(ctx, userId, opts)
	roles := s.SelectRoleAll(ctx, opts).Items
	for i := range roles {
		for _, userRole := range userRoles {
			if roles[i].RoleId == userRole.RoleId {
				roles[i].Flag = true
				break
			}
		}
	}
	return &v1.RoleList{
		ListMeta: api.ListMeta{int64(len(roles))},
		Items:    roles,
	}
}

func (s *roleService) SelectRolePermissionByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []string {
	perms, _ := s.store.Roles().SelectRolePermissionByUserId(ctx, userId, opts)
	permSet := hashset.New()
	for _, perm := range perms {
		if perm != nil {
			for _, str := range strings.Split(strings.TrimSpace(perm.RoleKey), ",") {
				permSet.Add(str)
			}
		}
	}
	res := make([]string, 0, permSet.Size())
	for _, v := range permSet.Values() {
		res = append(res, v.(string))
	}
	return res
}

func (s *roleService) SelectRoleAll(ctx context.Context, opts *api.GetOptions) *v1.RoleList {
	result, total, err := s.store.Roles().SelectRoleList(ctx, &v1.SysRole{}, opts)
	if err != nil {
		return &v1.RoleList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.RoleList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *roleService) SelectRoleListByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []int64 {
	result, _ := s.store.Roles().SelectRoleListByUserId(ctx, userId, opts)
	return result
}

func (s *roleService) SelectRoleById(ctx context.Context, roleId int64, opts *api.GetOptions) *v1.SysRole {
	result, _ := s.store.Roles().SelectRoleById(ctx, roleId, opts)
	return result
}

func (s *roleService) CheckRoleNameUnique(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) bool {
	result := s.store.Roles().CheckRoleNameUnique(ctx, role.RoleName, opts)
	if result != nil && result.RoleId != role.RoleId {
		return false
	}
	return true
}

func (s *roleService) CheckRoleKeyUnique(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) bool {
	result := s.store.Roles().CheckRoleKeyUnique(ctx, role.RoleKey, opts)
	if result != nil && result.RoleId != role.RoleId {
		return false
	}
	return true
}

func (s *roleService) CheckRoleAllowed(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) bool {
	if role != nil && role.IsAdmin() {
		return false
	}
	return true
}

func (s *roleService) CheckRoleDataScope(ctx context.Context, roleId int64, opts *api.GetOptions) bool {
	logininfor := utils.GetLogininfoFromCtx(ctx)
	if !v1.IsUserAdmin(logininfor.User.GetUserId()) {
		_, total, err := s.store.Roles().SelectRoleList(ctx, &v1.SysRole{RoleId: roleId}, opts)
		if err != nil || total == 0 {
			return false
		}
	}
	return true
}

func (s *roleService) CountUserRoleByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) int {
	return s.store.UserRoles().CountUserRoleByRoleId(ctx, roleId, opts)
}

func (s *roleService) InsertRole(ctx context.Context, role *v1.SysRole, opts *api.CreateOptions) error {
	if err := s.store.Roles().InsertRole(ctx, role, opts); err != nil {
		return err
	}
	s.InsertRoleMenu(ctx, role, opts)
	return nil
}

func (s *roleService) UpdateRole(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error {
	tx := s.store.Begin()
	if err := tx.Roles().UpdateRole(ctx, role, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除用户与菜单关联
	if err := tx.RoleMenus().DeleteRoleMenuByRoleId(ctx, role.RoleId, &api.DeleteOptions{}); err != nil {
		tx.Rollback()
		return err
	}
	// 新增用户与菜单关联
	if len(role.MenuIds) > 0 {
		list := make([]*v1.SysRoleDept, 0, len(role.DeptIds))
		for i := range role.DeptIds {
			list = append(list, &v1.SysRoleDept{RoleId: role.RoleId, DeptId: role.DeptIds[i]})
		}
		if err := tx.RoleDepts().BatchRoleDept(ctx, list, &api.CreateOptions{}); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *roleService) UpdateRoleStatus(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error {
	return nil
}

func (s *roleService) AuthDataScope(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error {
	return s.store.Roles().UpdateRole(ctx, role, opts)
}

func (s *roleService) DeleteRoleById(ctx context.Context, roleId int64, opts *api.DeleteOptions) error {
	tx := s.store.Begin()
	// 删除角色菜单关联
	if err := tx.RoleMenus().DeleteRoleMenuByRoleId(ctx, roleId, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除角色部门关联
	if err := tx.RoleDepts().DeleteRoleDeptByRoleId(ctx, roleId, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除角色信息
	if err := tx.Roles().DeleteRoleById(ctx, roleId, opts); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *roleService) DeleteRoleByIds(ctx context.Context, roleIds []int64, opts *api.DeleteOptions) error {
	for i := range roleIds {
		if !s.CheckRoleAllowed(ctx, &v1.SysRole{RoleId: roleIds[i]}, &api.GetOptions{Cache: true}) {
			return ErrRoleNotAllowed
		}
		if !s.CheckRoleDataScope(ctx, roleIds[i], &api.GetOptions{Cache: true}) {
			return ErrRoleDataScope
		}
	}
	tx := s.store.Begin()
	// 删除角色与菜单关联
	if err := tx.RoleMenus().DeleteRoleMenu(ctx, roleIds, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除角色与部门
	if err := tx.RoleDepts().DeleteRoleDept(ctx, roleIds, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除角色信息
	if err := tx.Roles().DeleteRoleByIds(ctx, roleIds, opts); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *roleService) DeleteAuthUser(ctx context.Context, userRole *v1.SysUserRole, opts *api.DeleteOptions) error {
	return s.store.UserRoles().DeleteUserRoleInfo(ctx, userRole, opts)
}

func (s *roleService) DeleteAuthUsers(ctx context.Context, roleId int64, userIds []int64, opts *api.DeleteOptions) error {
	return s.store.UserRoles().DeleteUserRoleInfos(ctx, roleId, userIds, opts)
}

func (s *roleService) InsertAuthUsers(ctx context.Context, roleId int64, userIds []int64, opts *api.CreateOptions) error {
	if len(userIds) > 0 {
		list := make([]*v1.SysUserRole, 0, len(userIds))
		for i := range userIds {
			list = append(list, &v1.SysUserRole{RoleId: roleId, UserId: userIds[i]})
		}
		return s.store.UserRoles().BatchUserRole(ctx, list, opts)
	}
	return nil
}

func (s *roleService) InsertRoleMenu(ctx context.Context, role *v1.SysRole, opts *api.CreateOptions) {
	if len(role.MenuIds) > 0 {
		list := make([]*v1.SysRoleMenu, 0, len(role.MenuIds))
		for i := range role.MenuIds {
			list = append(list, &v1.SysRoleMenu{RoleId: role.RoleId, MenuId: role.MenuIds[i]})
		}
		s.store.RoleMenus().BatchRoleMenu(ctx, list, opts)
	}
}

func (s *roleService) InsertRoleDept(ctx context.Context, role *v1.SysRole, opts *api.CreateOptions) {
	if len(role.DeptIds) > 0 {
		list := make([]*v1.SysRoleDept, 0, len(role.DeptIds))
		for i := range role.DeptIds {
			list = append(list, &v1.SysRoleDept{RoleId: role.RoleId, DeptId: role.DeptIds[i]})
		}
		s.store.RoleDepts().BatchRoleDept(ctx, list, opts)
	}
}
