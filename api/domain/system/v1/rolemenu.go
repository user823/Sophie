package v1

import "github.com/user823/Sophie/api"

type SysRoleMenu struct {
	RoleId int64 `json:"roleId" gorm:"column:role_id"`
	MenuId int64 `json:"menuId" gorm:"column:menu_id"`
}

func (s *SysRoleMenu) TableName() string {
	return "sys_role_menu"
}

type RoleMenuList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysRoleMenu `json:"items"`
}
