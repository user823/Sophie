package domain

import "github.com/user823/Sophie/api"

type SysUserRole struct {
	UserId int64 `json:"userId"`
	RoleId int64 `json:"roleId"`
}

func (s *SysUserRole) TableName() string {
	return "sys_user_role"
}

type UserRoleList struct {
	api.ListMeta `json:",inline"`
	Items        []SysUserRole `json:"items"`
}
