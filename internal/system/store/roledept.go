package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type RoleDeptStore interface {
	// 根据角色id 删除角色和部门关联
	DeleteRoleDeptByRoleId(ctx context.Context, roleid int64, opts *api.DeleteOptions) error
	// 批量删除角色部门信息
	DeleteRoleDept(ctx context.Context, ids []int64, opts *api.DeleteOptions) error
	// 查询部门使用数量
	SelectCountRoleDeptByDeptId(ctx context.Context, deptid int64, opts *api.GetOptions) int
	// 批量新增角色部门信息
	BatchRoleDept(ctx context.Context, roleDeptList []*v1.SysRoleDept, opts *api.CreateOptions) error
}
