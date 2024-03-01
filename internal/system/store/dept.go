package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type DeptStore interface {
	// 查询部门管理数据
	SelectDeptList(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) ([]*v1.SysDept, error)
	// 根据角色id查询部门树信息
	SelectDeptListByRoleId(ctx context.Context, roleid int64, deptCheckStrictly bool, opts *api.GetOptions) ([]int64, error)
	// 根据部门id查询信息
	SelectDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) (*v1.SysDept, error)
	// 根据id查询所有子部门
	SelectChildrenDeptById(ctx context.Context, deptid int64, opts *api.GetOptions) ([]*v1.SysDept, error)
	// 根据id查询所有子部门（正常状态）
	SelectNormalChildrenDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) int
	// 是否存在子节点
	HasChildByDeptId(ctx context.Context, deptid int64, opts *api.GetOptions) bool
	// 查询部门是否存在用户
	CheckDeptExistUser(ctx context.Context, deptid int64, opts *api.GetOptions) bool
	// 检查部门名称是否唯一
	CheckDeptNameUnique(ctx context.Context, name string, deptid int64, opts *api.GetOptions) *v1.SysDept
	// 新增部门信息
	InsertDept(ctx context.Context, dept *v1.SysDept, opts *api.CreateOptions) error
	// 修改部门信息
	UpdateDept(ctx context.Context, dept *v1.SysDept, opts *api.UpdateOptions) error
	// 修改所在部门状态
	UpdateDeptStatusNormal(ctx context.Context, deptids []int64, opts *api.UpdateOptions) error
	// 修改子元素关系
	UpdateDeptChildren(ctx context.Context, depts []*v1.SysDept, opts *api.UpdateOptions) error
	// 删除部门管理信息
	DeleteDeptById(ctx context.Context, deptid int64, opts *api.DeleteOptions) error
}
