package v1

import (
	"github.com/user823/Sophie/api/domain/system/v1"
)

// 设置登陆用户上下文信息
type LoginUser struct {
	// 用户唯一标识
	//Token string
	// 权限列表（根据角色权限获取菜单权限）
	Permissions []string `json:"permissions,omitempty"`
	// 角色列表
	Roles []string `json:"roles,omitempty"`
	// 用户信息
	User v1.SysUser `json:"user,omitempty"`
}
