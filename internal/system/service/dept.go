package service

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/api/domain/system/v1/vo"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/utils/intutil"
	"strconv"
	"strings"
)

type DeptSrv interface {
	// 查询部门管理数据
	SelectDeptList(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) *v1.DeptList
	// 查询部门树结构信息
	SelectDeptTreeList(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) []vo.TreeSelect
	// 构建前端所需树结构
	BuildDeptTree(ctx context.Context, depts []*v1.SysDept, opts *api.GetOptions) *v1.DeptList
	// 构建前端所需要下拉树结构
	BuildDeptTreeSelect(ctx context.Context, depts []*v1.SysDept, opts *api.GetOptions) []vo.TreeSelect
	// 根据角色ID查询部门树信息
	SelectDeptListByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) []int64
	// 根据部门ID查询信息
	SelectDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) *v1.SysDept
	// 根据ID查询所有子部门（正常状态）
	SelectNormalChildrenDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) int
	// 是否存在部门子节点
	HasChildByDeptId(ctx context.Context, deptId int64, opts *api.GetOptions) bool
	// 查询部门是否存在用户
	CheckDeptExistUser(ctx context.Context, deptId int64, opts *api.GetOptions) bool
	// 校验部门名称是否唯一
	CheckDeptNameUnique(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) bool
	// 校验部门是否有数据权限
	CheckDeptDataScope(ctx context.Context, deptId int64, opts *api.GetOptions) bool
	// 新增保存部门信息
	InsertDept(ctx context.Context, dept *v1.SysDept, opts *api.CreateOptions) error
	// 修改保存部门信息
	UpdateDept(ctx context.Context, dept *v1.SysDept, opts *api.UpdateOptions) error
	// 删除保存部门管理信息
	DeleteDeptById(ctx context.Context, deptId int64, opts *api.DeleteOptions) error
}

type deptService struct {
	store store.Factory
}

var _ DeptSrv = &deptService{}

func NewDepts(s store.Factory) DeptSrv {
	return &deptService{s}
}

func (s *deptService) SelectDeptList(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) *v1.DeptList {
	result, err := s.store.Depts().SelectDeptList(ctx, dept, opts)
	if err != nil {
		return &v1.DeptList{ListMeta: api.ListMeta{0}}
	}
	return &v1.DeptList{
		ListMeta: api.ListMeta{int64(len(result))},
		Items:    result,
	}
}

func (s *deptService) SelectDeptTreeList(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) []vo.TreeSelect {
	depts, err := s.store.Depts().SelectDeptList(ctx, dept, opts)
	if err != nil {
		return []vo.TreeSelect{}
	}

	return s.BuildDeptTreeSelect(ctx, depts, opts)
}

func (s *deptService) BuildDeptTree(ctx context.Context, depts []*v1.SysDept, opts *api.GetOptions) *v1.DeptList {
	list := make([]*v1.SysDept, 0, len(depts))
	deptIds := make([]int64, 0, len(depts))
	for i := range depts {
		deptIds = append(deptIds, depts[i].DeptId)
	}

	for i := range depts {
		// 如果是顶级节点，遍历该父节点的所有子节点
		if !intutil.ContainsAnyInt64(depts[i].ParentId, deptIds...) {
			deptRecursionFn(depts, depts[i])
			list = append(list, depts[i])
		}
	}
	if len(list) == 0 {
		list = depts
	}
	return &v1.DeptList{
		ListMeta: api.ListMeta{int64(len(list))},
		Items:    list,
	}
}

func (s *deptService) BuildDeptTreeSelect(ctx context.Context, depts []*v1.SysDept, opts *api.GetOptions) []vo.TreeSelect {
	deptTrees := s.BuildDeptTree(ctx, depts, opts)
	res := make([]vo.TreeSelect, 0, deptTrees.TotalCount)
	for i := range deptTrees.Items {
		res = append(res, deptTrees.Items[i].BuildTreeSelect())
	}
	return res
}

func (s *deptService) SelectDeptListByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) []int64 {
	role, err := s.store.Roles().SelectRoleById(ctx, roleId, opts)
	if err != nil {
		return []int64{}
	}
	result, err := s.store.Depts().SelectDeptListByRoleId(ctx, roleId, role.DeptCheckStrictly, opts)
	if err != nil {
		return []int64{}
	}
	return result
}

func (s *deptService) SelectDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) *v1.SysDept {
	result, err := s.store.Depts().SelectDeptById(ctx, deptId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *deptService) SelectNormalChildrenDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) int {
	return s.store.Depts().SelectNormalChildrenDeptById(ctx, deptId, opts)
}

func (s *deptService) HasChildByDeptId(ctx context.Context, deptId int64, opts *api.GetOptions) bool {
	return s.store.Depts().HasChildByDeptId(ctx, deptId, opts)
}

func (s *deptService) CheckDeptExistUser(ctx context.Context, deptId int64, opts *api.GetOptions) bool {
	return s.store.Depts().CheckDeptExistUser(ctx, deptId, opts)
}

func (s *deptService) CheckDeptNameUnique(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) bool {
	info := s.store.Depts().CheckDeptNameUnique(ctx, dept.DeptName, dept.DeptId, opts)
	if info != nil && info.DeptId == dept.DeptId {
		return false
	}
	return true
}

func (s *deptService) CheckDeptDataScope(ctx context.Context, deptId int64, opts *api.GetOptions) bool {
	logininfo, err := utils.GetLogininfoFromCtx(ctx)
	if err != nil {
		return false
	}
	if !logininfo.User.IsAdmin() {
		depts := s.SelectDeptList(ctx, &v1.SysDept{DeptId: deptId}, opts)
		if depts.TotalCount == 0 {
			// 没有权限访问部分数据
			return false
		}
	}
	return true
}

func (s *deptService) InsertDept(ctx context.Context, dept *v1.SysDept, opts *api.CreateOptions) error {
	if dept.ParentId != 0 {
		info, err := s.store.Depts().SelectDeptById(ctx, dept.ParentId, &api.GetOptions{Cache: true})
		if err != nil {
			return err
		}
		// 如果父部门不为正常状态，则不允许新增子节点
		if info.Status != v1.DEPTNORMAL {
			return fmt.Errorf("部门停用，不允许新增")
		}
		dept.Ancestors = info.Ancestors + "," + strconv.FormatInt(dept.ParentId, 10)
	}
	return s.store.Depts().InsertDept(ctx, dept, opts)
}

func (s *deptService) UpdateDept(ctx context.Context, dept *v1.SysDept, opts *api.UpdateOptions) error {
	newParentDept, _ := s.store.Depts().SelectDeptById(ctx, dept.ParentId, &api.GetOptions{Cache: true})
	oldDept, _ := s.store.Depts().SelectDeptById(ctx, dept.DeptId, &api.GetOptions{Cache: true})
	if newParentDept != nil && oldDept != nil {
		newAncestors := newParentDept.Ancestors + "," + strconv.FormatInt(newParentDept.DeptId, 10)
		oldAncestors := oldDept.Ancestors
		dept.Ancestors = newAncestors
		if err := s.UpdateDeptChildren(ctx, dept.DeptId, newAncestors, oldAncestors, opts); err != nil {
			return err
		}
	}
	err := s.store.Depts().UpdateDept(ctx, dept, opts)
	// 如果部门是启用状态则启用所有上级部门
	if dept.Status == v1.DEPTNORMAL && dept.Ancestors != "" && dept.Ancestors != "0" {
		s.UpdateParentDeptStatusNormal(ctx, dept, opts)
	}
	return err
}

func (s *deptService) DeleteDeptById(ctx context.Context, deptId int64, opts *api.DeleteOptions) error {
	return s.store.Depts().DeleteDeptById(ctx, deptId, opts)
}

func deptRecursionFn(list []*v1.SysDept, t *v1.SysDept) {
	// 得到子节点列表
	childList := getDeptChildList(list, t)
	t.Children = childList
	for i := range childList {
		if len(getDeptChildList(list, childList[i])) > 0 {
			deptRecursionFn(list, childList[i])
		}
	}
}

func getDeptChildList(list []*v1.SysDept, t *v1.SysDept) []*v1.SysDept {
	tlist := make([]*v1.SysDept, 0, len(list))
	for i := range list {
		if list[i].ParentId != 0 && list[i].ParentId == t.DeptId {
			tlist = append(tlist, list[i])
		}
	}
	return tlist
}

func (s *deptService) UpdateDeptChildren(ctx context.Context, deptId int64, newAncestors string, oldAncestors string, opts *api.UpdateOptions) error {
	children, _ := s.store.Depts().SelectChildrenDeptById(ctx, deptId, &api.GetOptions{Cache: true})
	for i := range children {
		children[i].Ancestors = strings.Replace(children[i].Ancestors, oldAncestors, newAncestors, 1)
	}
	if len(children) > 0 {
		return s.store.Depts().UpdateDeptChildren(ctx, children, opts)
	}
	return nil
}

// 修改该部门的父级部门状态
func (s *deptService) UpdateParentDeptStatusNormal(ctx context.Context, dept *v1.SysDept, opts *api.UpdateOptions) {
	strids := strings.Split(dept.Ancestors, ",")
	ids := make([]int64, len(strids))
	for i := range strids {
		id, _ := strconv.ParseInt(strids[i], 10, 64)
		ids = append(ids, id)
	}
	s.store.Depts().UpdateDeptStatusNormal(ctx, ids, opts)
}
