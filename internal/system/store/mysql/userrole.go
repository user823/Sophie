package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"gorm.io/gorm"
)

type mysqlUserRoleStore struct {
	db *gorm.DB
}

var _ store.UserRoleStore = &mysqlUserRoleStore{}

func (s *mysqlUserRoleStore) DeleteUserRoleByUserId(ctx context.Context, userid int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("user_id = ?", userid).Delete(&v1.SysUserRole{})
	return del.Error
}

func (s *mysqlUserRoleStore) DeleteUserRole(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("user_id in ?", ids).Delete(&v1.SysUserRole{})
	return del.Error
}

func (s *mysqlUserRoleStore) CountUserRoleByRoleId(ctx context.Context, roleid int64, opts *api.GetOptions) int {
	query := s.db.Table("sys_user_role").Where("role_id = ?", roleid)
	query = opts.SQLCondition(query, "")

	var result int64
	query.Count(&result)
	return int(result)
}

func (s *mysqlUserRoleStore) BatchUserRole(ctx context.Context, userRoleList []*v1.SysUserRole, opts *api.CreateOptions) error {
	for i := range userRoleList {
		s.db.Table("sys_user_role").Create(userRoleList[i])
	}
	return nil
}

func (s *mysqlUserRoleStore) DeleteUserRoleInfo(ctx context.Context, userRole *v1.SysUserRole, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("user_id = ? and role_id = ?", userRole.UserId, userRole.RoleId).Delete(&v1.SysUserRole{})
	return del.Error
}

// 批量取消授权用户角色
func (s *mysqlUserRoleStore) DeleteUserRoleInfos(ctx context.Context, roleid int64, userids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("role_id = ? and user_id in ?", roleid, userids).Delete(&v1.SysUserRole{})
	return del.Error
}
