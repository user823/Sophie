package main

import (
	"context"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
)

// SystemServiceImpl implements the last service interface defined in the IDL.
type SystemServiceImpl struct{}

// ListConfigs implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListConfigs(ctx context.Context, req *v1.ListConfigsRequest) (resp *v1.ListConfigsResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportConfig(ctx context.Context, req *v1.ExportConfigRequest) (resp *v1.ExportConfigResponse, err error) {
	// TODO: Your code here...
	return
}

// GetConfigById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetConfigById(ctx context.Context, id int64) (resp *v1.ConfigResponse, err error) {
	// TODO: Your code here...
	return
}

// GetConfigByKey implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetConfigByKey(ctx context.Context, key string) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// CreateConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateConfig(ctx context.Context, req *v1.CreateConfigRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateConfig(ctx context.Context, req *v1.UpdateConfigReqeust) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteConfig(ctx context.Context, req *v1.DeleteConfigReqeust) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// RefreshConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RefreshConfig(ctx context.Context) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListDepts implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDepts(ctx context.Context, req *v1.ListDeptsRequest) (resp *v1.ListDeptsResponse, err error) {
	// TODO: Your code here...
	return
}

// ListDeptsExcludeChild implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDeptsExcludeChild(ctx context.Context, id int64) (resp *v1.ListDeptsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetDeptById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetDeptById(ctx context.Context, req *v1.GetDeptByIdReq) (resp *v1.DeptResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateDept implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateDept(ctx context.Context, req *v1.CreateDeptRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateDept implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateDept(ctx context.Context, req *v1.UpdateDeptRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteDept implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteDept(ctx context.Context, req *v1.DeleteDeptRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListDictDatas implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDictDatas(ctx context.Context, req *v1.ListDictDatasRequest) (resp *v1.ListDictDatasResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportDictData(ctx context.Context, req *v1.ExportDictDataRequest) (resp *v1.ExportDictDataResponse, err error) {
	// TODO: Your code here...
	return
}

// GetDictDataByCode implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetDictDataByCode(ctx context.Context, code int64) (resp *v1.DictDataResponse, err error) {
	// TODO: Your code here...
	return
}

// ListDictDataByType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDictDataByType(ctx context.Context, dictType string) (resp *v1.ListDictDatasResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateDictData(ctx context.Context, req *v1.CreateDictDataRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateDictData(ctx context.Context, req *v1.UpdateDictDataRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteDictData(ctx context.Context, req *v1.DeleteDictDataRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListDictTypes implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDictTypes(ctx context.Context, req *v1.ListDictTypesRequest) (resp *v1.ListDictTypesResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportDictType(ctx context.Context, req *v1.ExportDictTypeRequest) (resp *v1.ExportDictTypeResponse, err error) {
	// TODO: Your code here...
	return
}

// GetDictTypeById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetDictTypeById(ctx context.Context, id int64) (resp *v1.DictTypeResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateDictType(ctx context.Context, req *v1.CreateDictTypeRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateDictType(ctx context.Context, req *v1.UpdateDictTypeRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteDictType(ctx context.Context, req *v1.DeleteDictTypeRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// RefreshDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RefreshDictType(ctx context.Context) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DictTypeOptionSelect implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DictTypeOptionSelect(ctx context.Context) (resp *v1.DictTypeOptionSelectResponse, err error) {
	// TODO: Your code here...
	return
}

// ListSysLogininfos implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysLogininfos(ctx context.Context, req *v1.ListSysLogininfosRequest) (resp *v1.ListSysLogininfosResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportLogininfo implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportLogininfo(ctx context.Context, req *v1.ExportLogininfoRequest) (resp *v1.ExportLogininfoResponse, err error) {
	// TODO: Your code here...
	return
}

// RemoveSysLogininfosById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RemoveSysLogininfosById(ctx context.Context, req *v1.RemoveSysLogininfosByIdRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// LogininfoClean implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) LogininfoClean(ctx context.Context) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UnlockByUserName implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UnlockByUserName(ctx context.Context, username string) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// CreateSysLogininfo implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysLogininfo(ctx context.Context, req *v1.CreateSysLogininfoRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListSysMenus implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysMenus(ctx context.Context, req *v1.ListSysMenusRequest) (resp *v1.ListSysMenusResponse, err error) {
	// TODO: Your code here...
	return
}

// GetSysMenuById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysMenuById(ctx context.Context, id int64) (resp *v1.SysMenuResponse, err error) {
	// TODO: Your code here...
	return
}

// ListTreeMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListTreeMenu(ctx context.Context, req *v1.ListTreeMenuRequest) (resp *v1.ListSysMenusResponse, err error) {
	// TODO: Your code here...
	return
}

// ListTreeMenuByRoleid implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListTreeMenuByRoleid(ctx context.Context, req *v1.ListTreeMenuByRoleidRequest) (resp *v1.RoleMenuResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateMenu(ctx context.Context, req *v1.CreateMenuRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateMenu(ctx context.Context, req *v1.UpdateMenuRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteMenu(ctx context.Context, req *v1.DeleteMenuRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// GetRouters implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetRouters(ctx context.Context, req *v1.GetRoutersRequest) (resp *v1.RoutersResonse, err error) {
	// TODO: Your code here...
	return
}

// ListSysNotices implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysNotices(ctx context.Context, req *v1.ListSysNoticesRequest) (resp *v1.ListSysNoticesResponse, err error) {
	// TODO: Your code here...
	return
}

// GetSysNoticeById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysNoticeById(ctx context.Context, id int64) (resp *v1.SysNoticeResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateSysNotice implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysNotice(ctx context.Context, req *v1.CreateSysNoticeRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteSysNotice implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysNotice(ctx context.Context, req *v1.DeleteSysNoticeRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateSysNotice implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysNotice(ctx context.Context, req *v1.UpdateSysNoticeRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListSysOperLogs implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysOperLogs(ctx context.Context, req *v1.ListSysOperLogsRequest) (resp *v1.ListSysOperLogsResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportSysOperLog implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysOperLog(ctx context.Context, req *v1.ExportSysOperLogRequest) (resp *v1.ExportSysOperLogResponse, err error) {
	// TODO: Your code here...
	return
}

// DeleteSysOperLog implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysOperLog(ctx context.Context, req *v1.DeleteSysOperLogRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// OperLogClean implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) OperLogClean(ctx context.Context) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// CreateSysOperLog implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysOperLog(ctx context.Context, req *v1.CreateSysOperLogRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListSysPosts implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysPosts(ctx context.Context, req *v1.ListSysPostsRequest) (resp *v1.ListSysPostsResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysPost(ctx context.Context, req *v1.ExportSysPostRequest) (resp *v1.ExportSysPostResponse, err error) {
	// TODO: Your code here...
	return
}

// GetSysPostById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysPostById(ctx context.Context, id int64) (resp *v1.SysPostResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysPost(ctx context.Context, req *v1.CreateSysPostRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysPost(ctx context.Context, req *v1.UpdateSysPostRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysPost(ctx context.Context, req *v1.DeleteSysPostRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// PostOptionSelect implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) PostOptionSelect(ctx context.Context) (resp *v1.PostOptionSelectResponse, err error) {
	// TODO: Your code here...
	return
}

// Profile implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) Profile(ctx context.Context, req *v1.ProfileRequest) (resp *v1.ProfileResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateProfile implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdatePassword implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdatePassword(ctx context.Context, req *v1.UpdatePasswordRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysRole(ctx context.Context, req *v1.ListSysRolesRequest) (resp *v1.ListSysRolesResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysRole(ctx context.Context, req *v1.ExportSysRoleRequest) (resp *v1.ExportSysRoleResponse, err error) {
	// TODO: Your code here...
	return
}

// GetSysRoleByid implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysRoleByid(ctx context.Context, id int64) (resp *v1.SysRoleResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysRole(ctx context.Context, req *v1.CreateSysRoleRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysRole(ctx context.Context, req *v1.UpdateSysRoleRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DataScope implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DataScope(ctx context.Context, req *v1.DataScopeRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ChangeSysRoleStatus implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ChangeSysRoleStatus(ctx context.Context, req *v1.ChangeSysRoleStatusRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysRole(ctx context.Context, req *v1.DeleteSysRoleRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListRoleOption implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListRoleOption(ctx context.Context) (resp *v1.ListSysRolesResponse, err error) {
	// TODO: Your code here...
	return
}

// AllocatedList implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) AllocatedList(ctx context.Context, req *v1.AllocatedListRequest) (resp *v1.ListSysUsersResponse, err error) {
	// TODO: Your code here...
	return
}

// UnallocatedList implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UnallocatedList(ctx context.Context, req *v1.UnallocatedListRequest) (resp *v1.ListSysUsersResponse, err error) {
	// TODO: Your code here...
	return
}

// CancelAuthUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CancelAuthUser(ctx context.Context, req *v1.CancelAuthUserRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// CancelAuthUserAll implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CancelAuthUserAll(ctx context.Context, req *v1.CancelAuthUserAllRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// SelectAuthUserAll implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) SelectAuthUserAll(ctx context.Context, req *v1.SelectAuthUserAllRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeptTreeByRoleId implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeptTreeByRoleId(ctx context.Context, id int64) (resp *v1.DeptTreeByRoleIdResponse, err error) {
	// TODO: Your code here...
	return
}

// ListSysUsers implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysUsers(ctx context.Context, req *v1.ListSysUsersRequest) (resp *v1.ListSysUsersResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysUser(ctx context.Context, req *v1.ExportSysUserRequest) (resp *v1.ExportSysUserResponse, err error) {
	// TODO: Your code here...
	return
}

// ImportUserData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ImportUserData(ctx context.Context, req *v1.ImportUserDataRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfoByName implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetUserInfoByName(ctx context.Context, name string) (resp *v1.UserInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetUserInfo(ctx context.Context, id int64) (resp *v1.UserInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// RegisterSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RegisterSysUser(ctx context.Context, req *v1.RegisterSysUserRequest) (resp *v1.RegisterSysUserResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfoById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetUserInfoById(ctx context.Context, req *v1.GetUserInfoByIdRequest) (resp *v1.UserInfoByIdResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysUser(ctx context.Context, req *v1.CreateSysUserRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// UpdateSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysUser(ctx context.Context, req *v1.UpdateSysUserRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// DeleteSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysUser(ctx context.Context, req *v1.DeleteSysUserRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ResetPassword implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ChangeSysUserStatus implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ChangeSysUserStatus(ctx context.Context, req *v1.ChangeSysUserStatus) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// GetAuthRoleById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetAuthRoleById(ctx context.Context, id int64) (resp *v1.AuthRoleInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// AuthRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) AuthRole(ctx context.Context, req *v1.AuthRoleRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}

// ListDeptsTree implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDeptsTree(ctx context.Context, req *v1.ListDeptsTreeRequest) (resp *v1.ListDeptsTreeResponse, err error) {
	// TODO: Your code here...
	return
}

// ListSysUserOnlines implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysUserOnlines(ctx context.Context, req *v1.ListSysUserOnlinesRequest) (resp *v1.ListSysUserOnline, err error) {
	// TODO: Your code here...
	return
}

// ForceLogout implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ForceLogout(ctx context.Context, req *v1.ForceLogoutRequest) (resp *v1.BaseResp, err error) {
	// TODO: Your code here...
	return
}
