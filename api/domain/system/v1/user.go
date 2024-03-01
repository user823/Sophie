package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/validators"
	"time"
)

type SysUser struct {
	api.ObjectMeta `json:",inline,omitempty"`
	UserId         int64      `json:"userId,omitempty" gorm:"column:user_id"`
	DeptId         int64      `json:"deptId,omitempty" gorm:"column:dept_id"`
	Username       string     `json:"userName,omitempty" gorm:"column:user_name" validate:"required,min=0,max=30,xss"`
	Nickname       string     `json:"nickName,omitempty" gorm:"column:nick_name" validate:"min=0,max=30,xss"`
	Email          string     `json:"email,omitempty" gorm:"column:email" validate:"email,min=0,max=50"`
	Phonenumber    string     `json:"phonenumber,omitempty" gorm:"column:phonenumber" validate:"min=0,max=11"`
	Sex            string     `json:"sex,omitempty" gorm:"column:sex"`
	Avatar         string     `json:"avatar,omitempty" gorm:"column:avatar"`
	Password       string     `json:"password,omitempty" gorm:"column:password"`
	Status         string     `json:"status,omitempty" gorm:"column:status"`
	DelFlag        string     `json:"delFlag,omitempty" gorm:"column:del_flag"`
	LoginIp        string     `json:"loginIp,omitempty" gorm:"column:login_ip"`
	LoginDate      *time.Time `json:"loginDate,omitempty" gorm:"column:login_date"`
	Dept           SysDept    `json:"dept,omitempty" gorm:"foreignKey:DeptId;references:DeptId"`
	Roles          []SysRole  `json:"roles,omitempty" gorm:"many2many:sys_user_role;foreignKey:UserId;joinForeignKey:UserId;references:RoleId;joinReferences:RoleId"`
	RoleIds        []int64    `json:"roleIds,omitempty" gorm:"-"`
	PostIds        []int64    `json:"postIds,omitempty" gorm:"-"`
	// 仅用于查询角色分配的用户，不参与存储
	RoleId int64 `json:"-" gorm:"-"`
}

func (u *SysUser) TableName() string {
	return "sys_user"
}

func (u *SysUser) IsAdmin() bool {
	return u.UserId == ROOT_ID
}

func (u *SysUser) String() string {
	data, _ := jsoniter.Marshal(u)
	return utils.B2s(data)
}

func (u *SysUser) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, u)
}

func (u *SysUser) Validate() error {
	vd := validators.GetValidatorOr()
	err := vd.Struct(u)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return buildErrMsg(e)
		}
	}
	return nil
}

func buildErrMsg(err validator.FieldError) error {
	switch err.StructNamespace() {
	case "SysUser.Username":
		return validators.BuildErrMsgHelper(err, "用户账号")
	case "SysUser.Nickname":
		return validators.BuildErrMsgHelper(err, "用户昵称")
	case "SysUser.Email":
		return validators.BuildErrMsgHelper(err, "邮箱")
	case "SysUser.Phonenumber":
		return validators.BuildErrMsgHelper(err, "手机号码")
	}
	return nil
}

type UserList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysUser `json:"items"`
}

// 用户状态
type Status struct {
	Code string
	Info string
}

var UserStatus = map[string]Status{
	"OK":      {USERNORMAL, "正常"},
	"DISABLE": {USERDISABLE, "停用"},
	"DELETED": {USERDELETED, "删除"},
}

func IsUserAdmin(userId int64) bool {
	return userId == ROOT_ID
}
