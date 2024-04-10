package v1

import (
	"github.com/mitchellh/mapstructure"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	vo2 "github.com/user823/Sophie/api/domain/vo"
	"github.com/user823/Sophie/pkg/utils"
)

func UserInfo2SysUser(userinfo *UserInfo) *v1.SysUser {
	if userinfo == nil {
		return nil
	}

	var roles []v1.SysRole
	for _, r := range userinfo.Roles {
		if r != nil {
			roles = append(roles, *RoleInfo2SysRole(r))
		}
	}

	user := &v1.SysUser{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  userinfo.CreateBy,
			CreatedAt: utils.Str2Time(userinfo.CreateTime),
			UpdateBy:  userinfo.UpdateBy,
			UpdatedAt: utils.Str2Time(userinfo.UpdateTime),
			Remark:    userinfo.Remark,
			Extend:    map[string]interface{}{"Params": userinfo.Params},
		},
		UserId:      userinfo.UserId,
		DeptId:      userinfo.GetDeptId(),
		Username:    userinfo.GetUserName(),
		Nickname:    userinfo.GetNickName(),
		Email:       userinfo.GetEmail(),
		Phonenumber: userinfo.GetPhonenumber(),
		Sex:         userinfo.GetSex(),
		Avatar:      userinfo.GetAvatar(),
		Password:    userinfo.GetPassword(),
		Status:      userinfo.GetStatus(),
		DelFlag:     userinfo.GetDelFlag(),
		Roles:       roles,
		RoleIds:     userinfo.RoleIds,
		PostIds:     userinfo.PostIds,
		RoleId:      userinfo.GetRoleId(),
	}

	if userinfo.Dept != nil {
		user.Dept = *DeptInfo2SysDept(userinfo.Dept)
	}
	return user
}

func MUserInfo2SysUser(userinfos []*UserInfo) []*v1.SysUser {
	sUsers := make([]*v1.SysUser, 0, len(userinfos))
	for i := range userinfos {
		sUsers = append(sUsers, UserInfo2SysUser(userinfos[i]))
	}
	return sUsers
}

func DeptInfo2SysDept(deptinfo *DeptInfo) *v1.SysDept {
	if deptinfo == nil {
		return nil
	}
	children := make([]*v1.SysDept, len(deptinfo.Children))
	for i := range deptinfo.Children {
		children = append(children, DeptInfo2SysDept(deptinfo.Children[i]))
	}

	return &v1.SysDept{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  deptinfo.CreateBy,
			CreatedAt: utils.Str2Time(deptinfo.CreateTime),
			UpdateBy:  deptinfo.UpdateBy,
			UpdatedAt: utils.Str2Time(deptinfo.UpdateTime),
			Remark:    deptinfo.Remark,
			Extend:    map[string]interface{}{"Params": deptinfo.Params},
		},
		DeptId:     deptinfo.DeptId,
		ParentId:   deptinfo.GetParentId(),
		Ancestors:  deptinfo.GetAncestors(),
		DeptName:   deptinfo.GetDeptName(),
		OrderNum:   deptinfo.GetOrderNum(),
		Leader:     deptinfo.GetLeader(),
		Phone:      deptinfo.GetPhone(),
		Email:      deptinfo.GetEmail(),
		Status:     deptinfo.GetStatus(),
		DelFlag:    deptinfo.GetDelFlag(),
		ParentName: deptinfo.GetParentName(),
		Children:   children,
	}
}

func MDeptInfo2SysDept(deptInfos []*DeptInfo) []*v1.SysDept {
	sDepts := make([]*v1.SysDept, 0, len(deptInfos))
	for i := range deptInfos {
		sDepts = append(sDepts, DeptInfo2SysDept(deptInfos[i]))
	}
	return sDepts
}

func RoleInfo2SysRole(roleinfo *RoleInfo) *v1.SysRole {
	if roleinfo == nil {
		return nil
	}
	return &v1.SysRole{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  roleinfo.CreateBy,
			CreatedAt: utils.Str2Time(roleinfo.CreateTime),
			UpdateBy:  roleinfo.UpdateBy,
			UpdatedAt: utils.Str2Time(roleinfo.UpdateTime),
			Remark:    roleinfo.Remark,
			Extend:    map[string]interface{}{"Params": roleinfo.Params},
		},
		RoleId:            roleinfo.RoleId,
		RoleName:          roleinfo.GetRoleName(),
		RoleKey:           roleinfo.GetRoleKey(),
		RoleSort:          roleinfo.GetRoleSort(),
		DataScope:         roleinfo.GetDataScope(),
		MenuCheckStrictly: roleinfo.GetMenuCheckStrictly(),
		DeptCheckStrictly: roleinfo.GetDeptCheckStrictly(),
		Status:            roleinfo.GetStatus(),
		DelFlag:           roleinfo.GetDelFlag(),
		Flag:              roleinfo.GetFlag(),
		MenuIds:           roleinfo.MenuIds,
		DeptIds:           roleinfo.DeptIds,
		Permissions:       roleinfo.Permissions,
	}
}

func MRoleInfo2SysRole(roleInfos []*RoleInfo) []*v1.SysRole {
	sRoles := make([]*v1.SysRole, 0, len(roleInfos))
	for i := range roleInfos {
		sRoles = append(sRoles, RoleInfo2SysRole(roleInfos[i]))
	}
	return sRoles
}

func LoginInfo2SysLogininfo(loginfo *Logininfo) *v1.SysLogininfor {
	if loginfo == nil {
		return nil
	}
	return &v1.SysLogininfor{
		InfoId:     loginfo.InfoId,
		UserName:   loginfo.GetUserName(),
		Ipaddr:     loginfo.GetIpaddr(),
		Status:     loginfo.GetStatus(),
		Msg:        loginfo.GetMsg(),
		AccessTime: utils.Str2Time(loginfo.GetAccessTime()),
	}
}

func MLoginInfo2SysLogininfo(loginfos []*Logininfo) []*v1.SysLogininfor {
	sLogins := make([]*v1.SysLogininfor, 0, len(loginfos))
	for i := range loginfos {
		sLogins = append(sLogins, LoginInfo2SysLogininfo(loginfos[i]))
	}
	return sLogins
}

func PostInfo2SysPost(postinfo *PostInfo) *v1.SysPost {
	if postinfo == nil {
		return nil
	}

	return &v1.SysPost{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  postinfo.CreateBy,
			CreatedAt: utils.Str2Time(postinfo.CreateTime),
			UpdateBy:  postinfo.UpdateBy,
			UpdatedAt: utils.Str2Time(postinfo.UpdateTime),
			Remark:    postinfo.Remark,
			Extend:    map[string]interface{}{"Params": postinfo.Params},
		},
		PostId:   postinfo.PostId,
		PostCode: postinfo.GetPostCode(),
		PostName: postinfo.GetPostName(),
		PostSort: postinfo.GetPostSort(),
		Status:   postinfo.GetStatus(),
		Flag:     postinfo.GetFlag(),
	}
}

func MPostInfo2SysPost(postinfos []*PostInfo) []*v1.SysPost {
	res := make([]*v1.SysPost, 0, len(postinfos))
	for i := range postinfos {
		res = append(res, PostInfo2SysPost(postinfos[i]))
	}
	return res
}

func ConfigInfo2SysConfig(configInfo *ConfigInfo) *v1.SysConfig {
	if configInfo == nil {
		return nil
	}

	return &v1.SysConfig{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  configInfo.CreateBy,
			CreatedAt: utils.Str2Time(configInfo.CreateTime),
			UpdateBy:  configInfo.UpdateBy,
			UpdatedAt: utils.Str2Time(configInfo.UpdateTime),
			Remark:    configInfo.Remark,
			Extend:    map[string]interface{}{"Params": configInfo.Params},
		},
		ConfigId:    configInfo.ConfigId,
		ConfigName:  configInfo.GetConfigName(),
		ConfigKey:   configInfo.GetConfigKey(),
		ConfigValue: configInfo.GetConfigValue(),
		ConfigType:  configInfo.GetConfigType(),
	}
}

func MConfigInfo2SysConfig(configs []*ConfigInfo) []*v1.SysConfig {
	res := make([]*v1.SysConfig, 0, len(configs))
	for i := range configs {
		res = append(res, ConfigInfo2SysConfig(configs[i]))
	}
	return res
}

func DictData2SysDictData(dictData *DictData) *v1.SysDictData {
	if dictData == nil {
		return nil
	}
	return &v1.SysDictData{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  dictData.CreateBy,
			CreatedAt: utils.Str2Time(dictData.CreateTime),
			UpdateBy:  dictData.UpdateBy,
			UpdatedAt: utils.Str2Time(dictData.UpdateTime),
			Remark:    dictData.Remark,
			Extend:    map[string]interface{}{"Params": dictData.Params},
		},
		DictCode:  dictData.GetDictCode(),
		DictSort:  dictData.GetDictSort(),
		DictLabel: dictData.GetDictLabel(),
		DictType:  dictData.GetDictType(),
		CssClass:  dictData.GetCssClass(),
		ListClass: dictData.GetListClass(),
		IsDefault: dictData.GetIsDefault(),
		Status:    dictData.GetStatus(),
	}
}

func MDictData2SysDictData(datas []*DictData) []*v1.SysDictData {
	res := make([]*v1.SysDictData, 0, len(datas))
	for i := range datas {
		res = append(res, DictData2SysDictData(datas[i]))
	}
	return res
}

func DictType2SysDictType(dictType *DictType) *v1.SysDictType {
	if dictType == nil {
		return nil
	}
	return &v1.SysDictType{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  dictType.CreateBy,
			CreatedAt: utils.Str2Time(dictType.CreateTime),
			UpdateBy:  dictType.UpdateBy,
			UpdatedAt: utils.Str2Time(dictType.UpdateTime),
			Remark:    dictType.Remark,
			Extend:    map[string]interface{}{"Params": dictType.Params},
		},
		DictId:   dictType.DictId,
		DictName: dictType.GetDictName(),
		DictType: dictType.GetDictType(),
		Status:   dictType.GetStatus(),
	}
}

func MDictType2SysDictType(dictTypes []*DictType) []*v1.SysDictType {
	res := make([]*v1.SysDictType, 0, len(dictTypes))
	for i := range dictTypes {
		res = append(res, DictType2SysDictType(dictTypes[i]))
	}
	return res
}

func MenuInfo2SysMenu(menuinfo *MenuInfo) *v1.SysMenu {
	if menuinfo == nil {
		return nil
	}

	children := make([]*v1.SysMenu, 0, len(menuinfo.Children))
	for i := range menuinfo.Children {
		children = append(children, MenuInfo2SysMenu(menuinfo.Children[i]))
	}
	return &v1.SysMenu{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  menuinfo.CreateBy,
			CreatedAt: utils.Str2Time(menuinfo.CreateTime),
			UpdateBy:  menuinfo.UpdateBy,
			UpdatedAt: utils.Str2Time(menuinfo.UpdateTime),
			Remark:    menuinfo.Remark,
			Extend:    map[string]interface{}{"Params": menuinfo.Params},
		},
		MenuId:     menuinfo.MenuId,
		MenuName:   menuinfo.GetMenuName(),
		ParentName: menuinfo.GetParentName(),
		ParentId:   menuinfo.GetParentId(),
		OrderNum:   menuinfo.GetOrderNum(),
		Path:       menuinfo.GetPath(),
		Component:  menuinfo.GetComponent(),
		Query:      menuinfo.GetQuery(),
		IsFrame:    menuinfo.GetIsFrame(),
		IsCache:    menuinfo.GetIsCache(),
		MenuType:   menuinfo.GetMenuType(),
		Visible:    menuinfo.GetVisible(),
		Status:     menuinfo.GetStatus(),
		Perms:      menuinfo.GetPerms(),
		Icon:       menuinfo.GetIcon(),
		Children:   children,
	}
}

func MMenuInfo2SysMenu(menuinfos []*MenuInfo) []*v1.SysMenu {
	res := make([]*v1.SysMenu, 0, len(menuinfos))
	for i := range menuinfos {
		res = append(res, MenuInfo2SysMenu(menuinfos[i]))
	}
	return res
}

func NoticeInfo2SysNotice(notice *NoticeInfo) *v1.SysNotice {
	if notice == nil {
		return nil
	}

	return &v1.SysNotice{
		ObjectMeta: api.ObjectMeta{
			CreateBy:  notice.CreateBy,
			CreatedAt: utils.Str2Time(notice.CreateTime),
			UpdateBy:  notice.UpdateBy,
			UpdatedAt: utils.Str2Time(notice.UpdateTime),
			Remark:    notice.Remark,
			Extend:    map[string]interface{}{"Params": notice.Params},
		},
		NoticeId:      notice.NoticeId,
		NoticeTitle:   notice.GetNoticeTitle(),
		NoticeType:    notice.GetNoticeType(),
		NoticeContent: notice.GetNoticeContent(),
		Status:        notice.GetStatus(),
	}
}

func MNoticeInfo2SysNotice(notices []*NoticeInfo) []*v1.SysNotice {
	res := make([]*v1.SysNotice, 0, len(notices))
	for i := range notices {
		res = append(res, NoticeInfo2SysNotice(notices[i]))
	}
	return res
}

func OperLog2SysOperLog(operlog *OperLog) *v1.SysOperLog {
	if operlog == nil {
		return nil
	}

	return &v1.SysOperLog{
		OperId:        operlog.GetOperId(),
		Title:         operlog.GetTitle(),
		BusinessType:  &operlog.BusinessType,
		Method:        operlog.GetMethod(),
		RequestMethod: operlog.GetRequestMethod(),
		OperatorType:  &operlog.OperatorType,
		OperName:      operlog.GetOperName(),
		DeptName:      operlog.GetDeptName(),
		OperUrl:       operlog.GetOperUrl(),
		OperIp:        operlog.GetOperIp(),
		OperParam:     operlog.GetOperParam(),
		JsonResult:    operlog.GetJsonResult_(),
		Status:        operlog.GetStatus(),
		ErrorMsg:      operlog.GetErrorMsg(),
		OperTime:      utils.Str2Time(operlog.GetOperTime()),
		CostTime:      operlog.GetCostTime(),
	}
}

func MOperLog2SysOperLog(logs []*OperLog) []*v1.SysOperLog {
	res := make([]*v1.SysOperLog, 0, len(logs))
	for i := range logs {
		res = append(res, OperLog2SysOperLog(logs[i]))
	}
	return res
}

func MTreeSelectTrans(vs []vo2.TreeSelect) []*TreeSelect {
	res := make([]*TreeSelect, 0, len(vs))
	for i := range vs {
		res = append(res, TreeSelectTrans(vs[i]))
	}
	return res
}

func TreeSelectTrans(v vo2.TreeSelect) *TreeSelect {
	return &TreeSelect{
		Id:       v.Id,
		Label:    v.Label,
		Children: MTreeSelectTrans(v.Children),
	}
}

func MRouterVoTrans(rs []vo2.RouterVo) []*RouterInfo {
	res := make([]*RouterInfo, 0, len(rs))
	for i := range rs {
		res = append(res, RouterVoTrans(rs[i]))
	}
	return res
}

func RouterVoTrans(r vo2.RouterVo) *RouterInfo {
	return &RouterInfo{
		Name:       r.Name,
		Path:       r.Path,
		Hidden:     r.Hidden,
		Redirect:   r.Redirect,
		Component:  r.Component,
		Query:      r.Query,
		AlwaysShow: r.AlwaysShow,
		Meta: &Meta{
			Title:   r.Meta.Title,
			Icon:    r.Meta.Icon,
			NoCache: r.Meta.NoCache,
			Link:    r.Meta.Link,
		},
		Children: MRouterVoTrans(r.Children),
	}
}

// logininfo 适配redis存储
func (l *Logininfo) Marshal() map[string]any {
	return map[string]any{
		"UserName":   l.UserName,
		"Status":     l.Status,
		"Ipaddr":     l.Ipaddr,
		"Msg":        l.Msg,
		"AccessTime": l.AccessTime,
		"TokenId":    l.TokenId,
		"DeptName":   l.DeptName,
		"Browser":    l.Browser,
		"Os":         l.Os,
		"LoginTime":  l.LoginTime,
	}
}

func (l *Logininfo) Unmarshal(mp map[string]string) {
	mapstructure.Decode(mp, l)
}

func LoginInfo2UserOnlineInfo(info *Logininfo) *UserOnlineInfo {
	return &UserOnlineInfo{
		UserName:  info.UserName,
		TokenId:   info.TokenId,
		Ipaddr:    info.Ipaddr,
		Browser:   info.Browser,
		Os:        info.Os,
		LoginTime: info.LoginTime,
	}
}
