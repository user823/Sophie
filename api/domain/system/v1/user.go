package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/validators"
)

type SysUser struct {
	api.ObjectMeta `json:",inline,omitempty"`
	UserId         int64     `json:"userId,omitempty" gorm:"column:user_id" query:"userId" xlsx:"n:用户编号"`
	DeptId         int64     `json:"deptId,omitempty" gorm:"column:dept_id" query:"deptId" xlsx:"n:部门编号"`
	Username       string    `json:"userName,omitempty" gorm:"column:user_name" query:"userName" validate:"required,min=0,max=30,xss" xlsx:"n:登录名称"`
	Nickname       string    `json:"nickName,omitempty" gorm:"column:nick_name" query:"nickName" validate:"min=0,max=30,xss" xlsx:"n:用户名称"`
	Email          string    `json:"email,omitempty" gorm:"column:email" query:"email" validate:"email,min=0,max=50" xlsx:"n:用户邮箱"`
	Phonenumber    string    `json:"phonenumber,omitempty" gorm:"column:phonenumber" query:"phonenumber" validate:"min=0,max=11" xlsx:"n:手机号码"`
	Sex            string    `json:"sex,omitempty" gorm:"column:sex" query:"sex" xlsx:"n:用户性别;exp:0=男,1=女,2=未知"`
	Avatar         string    `json:"avatar,omitempty" gorm:"column:avatar" query:"avatar" `
	Password       string    `json:"password,omitempty" gorm:"column:password" query:"password"`
	Status         string    `json:"status,omitempty" gorm:"column:status" query:"status" xlsx:"n:账号状态;exp:0=正常,1=停用"`
	DelFlag        string    `json:"delFlag,omitempty" gorm:"column:del_flag" query:"delFlag"`
	Dept           SysDept   `json:"dept,omitempty" gorm:"foreignKey:DeptId;references:DeptId" query:"dept"`
	Roles          []SysRole `json:"roles,omitempty" gorm:"many2many:sys_user_role;foreignKey:UserId;joinForeignKey:UserId;references:RoleId;joinReferences:RoleId" query:"roles"`
	RoleIds        []int64   `json:"roleIds,omitempty" gorm:"-" query:"roleIds"`
	PostIds        []int64   `json:"postIds,omitempty" gorm:"-" query:"postIds"`
	// 仅用于查询角色分配的用户，不参与存储
	RoleId int64 `json:"roleId,omitempty" gorm:"-" query:"roleId"`
}

func (u *SysUser) TableName() string {
	return "sys_user"
}

func (u *SysUser) IsAdmin() bool {
	return u.UserId == ROOT_ID
}

func (u *SysUser) Marshal() string {
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
			return buildUserErrMsg(e)
		}
	}
	return nil
}

// 过滤掉敏感数据
func (u *SysUser) Filter() *SysUser {
	u.Password = ""
	u.Dept.Filter()
	for i := range u.Roles {
		u.Roles[i].Filter()
	}
	return u
}

func buildUserErrMsg(err validator.FieldError) error {
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
