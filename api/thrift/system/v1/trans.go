package v1

import (
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/system/v1"
	"github.com/user823/Sophie/pkg/utils"
)

func BaseInfo2Meta(baseInfo *BaseInfo) api.ObjectMeta {
	return api.ObjectMeta{
		CreateBy:   baseInfo.CreateBy,
		CreateTime: utils.Second2Time(baseInfo.CreateTime),
		UpdateBy:   baseInfo.UpdateBy,
		UpdateTime: utils.Second2Time(baseInfo.UpdateTime),
		Remark:     baseInfo.Remark,
		Extend:     map[string]interface{}{},
	}
}

func UserInfo2SysUser(userinfo *UserInfo) *v1.SysUser {
	roles := make([]v1.SysRole, len(userinfo.Roles))
	for i := range userinfo.Roles {
		roles = append(roles, *RoleInfo2SysRole(userinfo.Roles[i]))
	}
	return &v1.SysUser{
		ObjectMeta:  BaseInfo2Meta(userinfo.BaseInfo),
		UserId:      userinfo.UserId,
		DeptId:      int64Trans(userinfo.DeptId),
		Username:    stringTrans(userinfo.UserName),
		Nickname:    stringTrans(userinfo.NickName),
		Email:       stringTrans(userinfo.Email),
		Phonenumber: stringTrans(userinfo.Phonenumber),
		Sex:         stringTrans(userinfo.Sex),
		Avatar:      stringTrans(userinfo.Avatar),
		Password:    stringTrans(userinfo.Password),
		Status:      stringTrans(userinfo.Status),
		DelFlag:     stringTrans(userinfo.Status),
		LoginIp:     stringTrans(userinfo.LoginIp),
		LoginDate:   utils.Second2Time(int64Trans(userinfo.LoginDate)),
		Dept:        *DeptInfo2SysDept(userinfo.Dept),
		Roles:       roles,
		RoleIds:     userinfo.RoleIds,
		PostIds:     userinfo.PostIds,
		RoleId:      int64Trans(userinfo.RoleId),
	}
}

func DeptInfo2SysDept(deptinfo *DeptInfo) *v1.SysDept {
	children := make([]v1.SysDept, len(deptinfo.Children))
	for i := range deptinfo.Children {
		children = append(children, *DeptInfo2SysDept(deptinfo.Children[i]))
	}

	return &v1.SysDept{
		ObjectMeta: BaseInfo2Meta(deptinfo.BaseInfo),
		DeptId:     deptinfo.DeptId,
		ParentId:   int64Trans(deptinfo.ParentId),
		Ancestors:  stringTrans(deptinfo.Ancestors),
		DeptName:   stringTrans(deptinfo.DeptName),
		OrderNum:   int64Trans(deptinfo.OrderNum),
		Leader:     stringTrans(deptinfo.Leader),
		Phone:      stringTrans(deptinfo.Phone),
		Email:      stringTrans(deptinfo.Email),
		Status:     stringTrans(deptinfo.Status),
		DelFlag:    stringTrans(deptinfo.DelFlag),
		ParentName: stringTrans(deptinfo.ParentName),
		Children:   children,
	}
}

func RoleInfo2SysRole(roleinfo *RoleInfo) *v1.SysRole {
	return &v1.SysRole{
		ObjectMeta:        BaseInfo2Meta(roleinfo.BaseInfo),
		RoleId:            roleinfo.RoleId,
		RoleName:          stringTrans(roleinfo.RoleName),
		RoleKey:           stringTrans(roleinfo.RoleKey),
		RoleSort:          int64Trans(roleinfo.RoleSort),
		DataScope:         stringTrans(roleinfo.DataScope),
		MenuCheckStrictly: boolTrans(roleinfo.MenuCheckStrictly),
		DeptCheckStrictly: boolTrans(roleinfo.DeptCheckStrictly),
		Status:            stringTrans(roleinfo.Status),
		DelFlag:           stringTrans(roleinfo.DelFlag),
		Flag:              boolTrans(roleinfo.Flag),
		MenuIds:           roleinfo.MenuIds,
		DeptIds:           roleinfo.DeptIds,
		Permissions:       roleinfo.Permissions,
	}
}

func int64Trans(a *int64) int64 {
	if a == nil {
		return 0
	}
	return *a
}

func stringTrans(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

func boolTrans(a *bool) bool {
	if a == nil {
		return false
	}
	return *a
}
