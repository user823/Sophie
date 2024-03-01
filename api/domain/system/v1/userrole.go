package v1

import "github.com/user823/Sophie/api"

type SysUserRole struct {
	UserId int64 `json:"userId" gorm:"column:user_id"`
	RoleId int64 `json:"roleId" gorm:"column:role_id"`
}

func (s *SysUserRole) TableName() string {
	return "sys_user_role"
}

type UserRoleList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysUserRole `json:"items"`
}
