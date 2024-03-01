package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"gorm.io/gorm"
)

type mysqlRoleDeptStore struct {
	db *gorm.DB
}

var _ store.RoleDeptStore = &mysqlRoleDeptStore{}

func (s *mysqlRoleDeptStore) DeleteRoleDeptByRoleId(ctx context.Context, roleid int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("role_id = ?", roleid).Delete(v1.SysRoleDept{})
	return del.Error
}

func (s *mysqlRoleDeptStore) DeleteRoleDept(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("role_id in ?", ids).Delete(v1.SysRoleDept{})
	return del.Error
}

func (s *mysqlRoleDeptStore) SelectCountRoleDeptByDeptId(ctx context.Context, deptid int64, opts *api.GetOptions) int {
	query := s.db.Table("sys_role_dept").Where("dept_id = ?", deptid)
	query = opts.SQLCondition(query, "")

	var result int64
	query.Count(&result)
	return int(result)
}

func (s *mysqlRoleDeptStore) BatchRoleDept(ctx context.Context, roleDeptList []*v1.SysRoleDept, opts *api.CreateOptions) error {
	for i := range roleDeptList {
		opts.SQLCondition(s.db).Table("sys_role_dept").Create(roleDeptList[i])
	}
	return nil
}
