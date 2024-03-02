package v1

import (
	v12 "github.com/user823/Sophie/api/domain/gateway/v1"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/pkg/utils"
)

func SysUser2UserInfo(user *v1.SysUser) *UserInfo {
	if user == nil {
		return nil
	}
	var params map[string]string
	if data, ok := user.Extend["Params"].(map[string]string); ok {
		params = data
	}

	deptInfo := SysDept2DeptInfo(&user.Dept)
	var roles []*RoleInfo
	for i := range user.Roles {
		roles = append(roles, SysRole2RoleInfo(&user.Roles[i]))
	}

	userInfo := &UserInfo{
		CreateBy:    user.CreateBy,
		CreateTime:  utils.Time2Second(user.CreatedAt),
		UpdateBy:    user.UpdateBy,
		UpdateTime:  utils.Time2Second(user.UpdatedAt),
		Remark:      user.Remark,
		Params:      params,
		UserId:      user.UserId,
		DeptId:      user.DeptId,
		UserName:    user.Username,
		NickName:    user.Nickname,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      user.Avatar,
		Password:    user.Password,
		Status:      user.Status,
		DelFlag:     user.DelFlag,
		LoginIp:     user.LoginIp,
		Dept:        deptInfo,
		Roles:       roles,
		RoleIds:     user.RoleIds,
		PostIds:     user.PostIds,
		RoleId:      user.RoleId,
	}

	if user.LoginDate != nil {
		userInfo.LoginDate = utils.Time2Second(*user.LoginDate)
	}
	return userInfo
}

func MSysUser2UserInfo(users []*v1.SysUser) []*UserInfo {
	res := make([]*UserInfo, 0, len(users))
	for i := range users {
		res = append(res, SysUser2UserInfo(users[i]))
	}
	return res
}

func SysDept2DeptInfo(dept *v1.SysDept) *DeptInfo {
	if dept == nil {
		return nil
	}

	var params map[string]string
	if data, ok := dept.Extend["Params"].(map[string]string); ok {
		params = data
	}

	var children []*DeptInfo
	for i := range dept.Children {
		children = append(children, SysDept2DeptInfo(dept.Children[i]))
	}

	return &DeptInfo{
		CreateBy:   dept.CreateBy,
		CreateTime: utils.Time2Second(dept.CreatedAt),
		UpdateBy:   dept.UpdateBy,
		UpdateTime: utils.Time2Second(dept.UpdatedAt),
		Remark:     dept.Remark,
		Params:     params,
		DeptId:     dept.DeptId,
		ParentId:   dept.ParentId,
		Ancestors:  dept.Ancestors,
		DeptName:   dept.DeptName,
		OrderNum:   dept.OrderNum,
		Leader:     dept.Leader,
		Phone:      dept.Phone,
		Email:      dept.Email,
		Status:     dept.Status,
		DelFlag:    dept.DelFlag,
		ParentName: dept.ParentName,
		Children:   children,
	}
}

func MSysDept2DeptInfo(depts []*v1.SysDept) []*DeptInfo {
	res := make([]*DeptInfo, 0, len(depts))
	for i := range depts {
		res = append(res, SysDept2DeptInfo(depts[i]))
	}
	return res
}

func SysRole2RoleInfo(role *v1.SysRole) *RoleInfo {
	if role == nil {
		return nil
	}

	var params map[string]string
	if data, ok := role.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &RoleInfo{
		CreateBy:          role.CreateBy,
		CreateTime:        utils.Time2Second(role.CreatedAt),
		UpdateBy:          role.UpdateBy,
		UpdateTime:        utils.Time2Second(role.UpdatedAt),
		Remark:            role.Remark,
		Params:            params,
		RoleId:            role.RoleId,
		RoleName:          role.RoleName,
		RoleKey:           role.RoleKey,
		RoleSort:          role.RoleSort,
		DataScope:         role.DataScope,
		MenuCheckStrictly: role.MenuCheckStrictly,
		DeptCheckStrictly: role.DeptCheckStrictly,
		Status:            role.Status,
		DelFlag:           role.DelFlag,
		Flag:              role.Flag,
		MenuIds:           role.MenuIds,
		DeptIds:           role.DeptIds,
		Permissions:       role.Permissions,
	}
}

func MSysRole2RoleInfo(roles []*v1.SysRole) []*RoleInfo {
	res := make([]*RoleInfo, 0, len(roles))
	for i := range roles {
		res = append(res, SysRole2RoleInfo(roles[i]))
	}
	return res
}

func SysPost2PostInfo(post *v1.SysPost) *PostInfo {
	if post == nil {
		return nil
	}

	var params map[string]string
	if data, ok := post.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &PostInfo{
		CreateBy:   post.CreateBy,
		CreateTime: utils.Time2Second(post.CreatedAt),
		UpdateBy:   post.UpdateBy,
		UpdateTime: utils.Time2Second(post.UpdatedAt),
		Remark:     post.Remark,
		Params:     params,
		PostId:     post.PostId,
		PostCode:   post.PostCode,
		PostSort:   post.PostSort,
		Status:     post.Status,
		Flag:       post.Flag,
	}
}

func MSysPost2PostInfo(posts []*v1.SysPost) []*PostInfo {
	res := make([]*PostInfo, 0, len(posts))
	for i := range posts {
		res = append(res, SysPost2PostInfo(posts[i]))
	}
	return res
}

func SysConfig2ConfigInfo(config *v1.SysConfig) *ConfigInfo {
	if config == nil {
		return nil
	}

	var params map[string]string
	if data, ok := config.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &ConfigInfo{
		CreateBy:    config.CreateBy,
		CreateTime:  utils.Time2Second(config.CreatedAt),
		UpdateBy:    config.UpdateBy,
		UpdateTime:  utils.Time2Second(config.UpdatedAt),
		Remark:      config.Remark,
		Params:      params,
		ConfigId:    config.ConfigId,
		ConfigName:  config.ConfigName,
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		ConfigType:  config.ConfigType,
	}
}

func MSysConfig2ConfigInfo(configs []*v1.SysConfig) []*ConfigInfo {
	res := make([]*ConfigInfo, 0, len(configs))
	for i := range configs {
		res = append(res, SysConfig2ConfigInfo(configs[i]))
	}
	return res
}

func SysDictData2DictData(dictData *v1.SysDictData) *DictData {
	if dictData == nil {
		return nil
	}

	var params map[string]string
	if data, ok := dictData.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &DictData{
		CreateBy:   dictData.CreateBy,
		CreateTime: utils.Time2Second(dictData.CreatedAt),
		UpdateBy:   dictData.UpdateBy,
		UpdateTime: utils.Time2Second(dictData.UpdatedAt),
		Remark:     dictData.Remark,
		Params:     params,
		DictCode:   dictData.DictCode,
		DictSort:   dictData.DictSort,
		DictLabel:  dictData.DictLabel,
		DictValue:  dictData.DictValue,
		DictType:   dictData.DictType,
		CssClass:   dictData.CssClass,
		ListClass:  dictData.ListClass,
		IsDefault:  dictData.IsDefault,
		Status:     dictData.Status,
	}
}

func MSysDictData2DictData(dictDatas []*v1.SysDictData) []*DictData {
	res := make([]*DictData, 0, len(dictDatas))
	for i := range dictDatas {
		res = append(res, SysDictData2DictData(dictDatas[i]))
	}
	return res
}

func SysDictType2DictType(dictType *v1.SysDictType) *DictType {
	if dictType == nil {
		return nil
	}

	var params map[string]string
	if data, ok := dictType.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &DictType{
		CreateBy:   dictType.CreateBy,
		CreateTime: utils.Time2Second(dictType.CreatedAt),
		UpdateBy:   dictType.UpdateBy,
		UpdateTime: utils.Time2Second(dictType.UpdatedAt),
		Remark:     dictType.Remark,
		Params:     params,
		DictId:     dictType.DictId,
		DictName:   dictType.DictName,
		DictType:   dictType.DictType,
		Status:     dictType.Status,
	}
}

func MSysDictType2DictType(dictTypes []*v1.SysDictType) []*DictType {
	res := make([]*DictType, 0, len(dictTypes))
	for i := range dictTypes {
		res = append(res, SysDictType2DictType(dictTypes[i]))
	}
	return res
}

func SysLogininfor2Logininfor(logininfor *v1.SysLogininfor) *Logininfo {
	if logininfor == nil {
		return nil
	}

	accessTime := utils.Time2Second(logininfor.AccessTime)
	return &Logininfo{
		InfoId:     logininfor.InfoId,
		UserName:   logininfor.UserName,
		Status:     logininfor.Status,
		Ipaddr:     logininfor.Ipaddr,
		Msg:        logininfor.Msg,
		AccessTime: accessTime,
	}
}

func MSysLogininfor2Logininfor(logininfors []*v1.SysLogininfor) []*Logininfo {
	res := make([]*Logininfo, 0, len(logininfors))
	for i := range logininfors {
		res = append(res, SysLogininfor2Logininfor(logininfors[i]))
	}
	return res
}

func SysMenu2MenuInfo(menu *v1.SysMenu) *MenuInfo {
	if menu == nil {
		return nil
	}

	var params map[string]string
	if data, ok := menu.Extend["Params"].(map[string]string); ok {
		params = data
	}

	children := make([]*MenuInfo, 0, len(menu.Children))
	for i := range menu.Children {
		children = append(children, SysMenu2MenuInfo(menu.Children[i]))
	}
	return &MenuInfo{
		CreateBy:   menu.CreateBy,
		CreateTime: utils.Time2Second(menu.CreatedAt),
		UpdateBy:   menu.UpdateBy,
		UpdateTime: utils.Time2Second(menu.UpdatedAt),
		Remark:     menu.Remark,
		Params:     params,
		MenuId:     menu.MenuId,
		MenuName:   menu.MenuName,
		ParentName: menu.ParentName,
		ParentId:   menu.ParentId,
		OrderNum:   menu.OrderNum,
		Path:       menu.Path,
		Component:  menu.Component,
		Query:      menu.Query,
		IsFrame:    menu.IsFrame,
		IsCache:    menu.IsCache,
		MenuType:   menu.MenuType,
		Visible:    menu.Visible,
		Status:     menu.Status,
		Perms:      menu.Perms,
		Icon:       menu.Icon,
		Children:   children,
	}
}

func MSysMenu2MenuInfo(menus []*v1.SysMenu) []*MenuInfo {
	res := make([]*MenuInfo, 0, len(menus))
	for i := range menus {
		res = append(res, SysMenu2MenuInfo(menus[i]))
	}
	return res
}

func SysNotice2NoticeInfo(notice *v1.SysNotice) *NoticeInfo {
	if notice == nil {
		return nil
	}

	var params map[string]string
	if data, ok := notice.Extend["Params"].(map[string]string); ok {
		params = data
	}

	return &NoticeInfo{
		CreateBy:      notice.CreateBy,
		CreateTime:    utils.Time2Second(notice.CreatedAt),
		UpdateBy:      notice.UpdateBy,
		UpdateTime:    utils.Time2Second(notice.UpdatedAt),
		Remark:        notice.Remark,
		Params:        params,
		NoticeId:      notice.NoticeId,
		NoticeTitle:   notice.NoticeTitle,
		NoticeType:    notice.NoticeType,
		NoticeContent: notice.NoticeContent,
		Status:        notice.Status,
	}
}

func MSysNotice2NoticeInfo(notices []*v1.SysNotice) []*NoticeInfo {
	res := make([]*NoticeInfo, 0, len(notices))
	for i := range notices {
		res = append(res, SysNotice2NoticeInfo(notices[i]))
	}
	return res
}

func SysOperLog2OperLog(log *v1.SysOperLog) *OperLog {
	if log == nil {
		return nil
	}

	var businessType, operatorType int64
	if log.BusinessType != nil {
		businessType = *log.BusinessType
	}
	if log.OperatorType != nil {
		operatorType = *log.OperatorType
	}

	operTime := utils.Time2Second(log.OperTime)
	return &OperLog{
		OperId:        log.OperId,
		Title:         log.Title,
		BusinessType:  businessType,
		Method:        log.Method,
		RequestMethod: log.RequestMethod,
		OperatorType:  operatorType,
		OperName:      log.OperName,
		DeptName:      log.DeptName,
		OperUrl:       log.OperUrl,
		OperIp:        log.OperIp,
		OperParam:     log.OperParam,
		JsonResult_:   log.JsonResult,
		Status:        log.Status,
		ErrorMsg:      log.ErrorMsg,
		OperTime:      operTime,
		CostTime:      log.CostTime,
	}
}

func MSysOperLog2Operlog(logs []*v1.SysOperLog) []*OperLog {
	res := make([]*OperLog, 0, len(logs))
	for i := range logs {
		res = append(res, SysOperLog2OperLog(logs[i]))
	}
	return res
}

func LoginUserTrans(loginUser *v12.LoginUser) *LoginUser {
	return &LoginUser{
		Permissions: loginUser.Permissions,
		Roles:       loginUser.Roles,
		User:        SysUser2UserInfo(&loginUser.User),
	}
}
