package v1

import (
	"github.com/user823/Sophie/api"
	"time"
)

type SysUser struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	UserId         int64     `json:"user_id" gorm:"column:user_id"`
	DeptId         int64     `json:"dept_id" gorm:"column:dept_id"`
	Username       string    `json:"user_name" gorm:"column:user_name"`
	Nickname       int64     `json:"nick_name" gorm:"column:nick_name"`
	Email          string    `json:"email" gorm:"column:email"`
	Phonenumber    string    `json:"phonenumber" gorm:"column:phonenumber"`
	Sex            string    `json:"sex" gorm:"column:sex"`
	Avatar         string    `json:"avatar" gorm:"column:avatar"`
	Password       string    `json:"password" gorm:"column:password"`
	Status         string    `json:"status" gorm:"column:status"`
	DelFlag        string    `json:"del_flag" gorm:"column:del_flag"`
	LoginIp        string    `json:"login_ip" gorm:"column:login_ip"`
	LoginDate      time.Time `json:"login_date" gorm:"column:login_date"`
	Dept           SysDept
	Roles          []SysRole
	RoleIds        []int64
	PostIds        []int64
	RoleId         int64
}

type UserList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysUser `json:"items"`
}

func (u *SysUser) TableName() string {
	return "sys_user"
}
