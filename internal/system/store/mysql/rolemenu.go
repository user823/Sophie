package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"gorm.io/gorm"
)

type mysqlRoleMenuStore struct {
	db *gorm.DB
}

var _ store.RoleMenuStore = &mysqlRoleMenuStore{}

func (s *mysqlRoleMenuStore) CheckMenuExistRole(ctx context.Context, menuid int64, opts *api.GetOptions) int {
	query := s.db.Table("sys_role_menu").Where("menu_id = ?", menuid)
	query = opts.SQLCondition(query, "")

	var result int64
	query.Count(&result)
	return int(result)
}

func (s *mysqlRoleMenuStore) DeleteRoleMenuByRoleId(ctx context.Context, roleid int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("role_id = ?", roleid).Delete(&v1.SysRoleMenu{})
	return del.Error
}

func (s *mysqlRoleMenuStore) DeleteRoleMenu(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("role_id in ?", ids).Delete(&v1.SysRoleMenu{})
	return del.Error
}

func (s *mysqlRoleMenuStore) BatchRoleMenu(ctx context.Context, roleMenuList []*v1.SysRoleMenu, opts *api.CreateOptions) error {
	for i := range roleMenuList {
		opts.SQLCondition(s.db).Table("sys_role_menu").Create(roleMenuList[i])
	}
	return nil
}
