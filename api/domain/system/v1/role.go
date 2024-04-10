package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysRole struct {
	api.ObjectMeta `json:",inline,omitempty"`
	RoleId         int64  `json:"roleId,omitempty" gorm:"column:role_id" query:"roleId" xlsx:"n:角色序号"`
	RoleName       string `json:"roleName,omitempty" gorm:"column:role_name" query:"roleName" xlsx:"n:角色名称"`
	RoleKey        string `json:"roleKey,omitempty" gorm:"column:role_key" query:"roleKey" xlsx:"n:角色权限"`
	// 角色排序
	RoleSort          int64  `json:"roleSort,omitempty" gorm:"column:role_sort" query:"roleSort" xlsx:"n:角色排序"`
	DataScope         string `json:"dataScope,omitempty" gorm:"column:data_scope" query:"dataScope" xlsx:"n:数据范围;exp:1=所有数据权限,2=自定义数据权限,3=本部门数据权限,4=本部门及子部门数据权限,5=仅本人数据权限"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly,omitempty" gorm:"column:menu_check_strictly" query:"menuCheckStrictly"`
	DeptCheckStrictly bool   `json:"deptCheckStrictly,omitempty" gorm:"column:dept_check_strictly" query:"deptCheckStrictly"`
	Status            string `json:"status,omitempty" gorm:"column:status" query:"status" xlsx:"n:角色状态;exp:0=正常,1=停用"`
	DelFlag           string `json:"delFlag,omitempty" gorm:"column:del_flag" query:"delFlag"`
	// 用户是否存在此角色的标识（默认不存在）, 用于鉴权
	Flag    bool    `json:"flag,omitempty" gorm:"-" query:"flag"`
	MenuIds []int64 `json:"menuIds,omitempty" gorm:"-" query:"menuIds"`
	DeptIds []int64 `json:"deptIds,omitempty" gorm:"-" query:"deptIds"`
	// 角色菜单权限
	Permissions []string `json:"permissions,omitempty" gorm:"-" query:"permissions"`
}

func (s *SysRole) TableName() string {
	return "sys_role"
}

func (s *SysRole) IsAdmin() bool {
	return s.RoleId == 1
}

func (s *SysRole) Marshal() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysRole) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

func (s *SysRole) Filter() *SysRole {
	return s
}

type RoleList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysRole `json:"items"`
}

func IsRoleAdmin(roleId int64) bool {
	return roleId == 1
}
