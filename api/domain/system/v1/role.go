package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysRole struct {
	api.ObjectMeta `json:",inline,omitempty"`
	RoleId         int64  `json:"roleId,omitempty" gorm:"column:role_id"`
	RoleName       string `json:"roleName,omitempty" gorm:"column:role_name"`
	RoleKey        string `json:"roleKey,omitempty" gorm:"column:role_key"`
	// 角色排序
	RoleSort          int64  `json:"roleSort,omitempty" gorm:"column:role_sort"`
	DataScope         string `json:"dataScope,omitempty" gorm:"column:data_scope"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly,omitempty" gorm:"column:menu_check_strictly"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly,omitempty" gorm:"column:dept_check_strictly"`
	Status            string `json:"status,omitempty" gorm:"column:status"`
	DelFlag           string `json:"delFlag,omitempty" gorm:"column:del_flag"`
	// 用户是否存在此角色的标识（默认不存在）, 用于鉴权
	Flag    bool    `json:"flag,omitempty" gorm:"-"`
	MenuIds []int64 `json:"-" gorm:"-"`
	DeptIds []int64 `json:"-" gorm:"-"`
	// 角色菜单权限
	Permissions []string `json:"-" gorm:"-"`
}

func (s *SysRole) TableName() string {
	return "sys_role"
}

func (s *SysRole) IsAdmin() bool {
	return s.RoleId == 1
}

func (s *SysRole) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysRole) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type RoleList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysRole `json:"items"`
}

func IsRoleAdmin(roleId int64) bool {
	return roleId == 1
}
