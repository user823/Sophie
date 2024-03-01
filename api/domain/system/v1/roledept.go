package v1

import "github.com/user823/Sophie/api"

type SysRoleDept struct {
	RoleId int64 `json:"roleId" gorm:"column:role_id"`
	DeptId int64 `json:"deptId" gorm:"column:dept_id"`
}

func (s *SysRoleDept) TableName() string {
	return "sys_role_dept"
}

type RoleDeptList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysRoleDept `json:"items"`
}
