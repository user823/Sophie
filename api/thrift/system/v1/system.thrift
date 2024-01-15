namespace go github.com/user823/Sophie/api/thrift/system/v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/system/v1 -service SystemService -gen-path api/thrift/system  system.thrift

service SystemService {
    // config service
    ListConfigsResponse ListConfigs(1:ListConfigRequest req)
    ExportConfigResponse ExportConfig(1:ExportConfigRequest req)
    ConfigResponse GetConfigById(1:i64 id)
    ConfigResponse GetConfigByKey(1:string key)
    CreateConfigResponse CreateConfig(1:CreateConfigRequest req)
    UpdateConfigResponse UpdateConfig(1:UpdateConfigReqeust req)
    DeleteConfigResponse DeleteConfig(1:DeleteConfigReqeust req)
    RefreshConfigResponse RefreshConfig()

    // dept service
    ListDeptsResponse ListDepts(1:ListDeptsRequest req)
    ListDeptsResponse ListDeptsExcludeChild(1:i64 id)
    DeptResponse GetDeptById(1:i64 id)
    CreateDeptResponse CreateDept(1:CreateDeptRequest req)
    UpdateDeptResponse UpdateDept(1:UpdateDeptRequest req)
    DeleteDeptResponse DeleteDept(1:DeleteDeptRequest req)

    // dict data service
    ListDictDatasResponse ListDictDatas(1:ListDictDatasRequest req)
    ExportDictDataResponse ExportDictData(1:ExportDictDataRequest req)
    DictDataResponse GetDictDataByCode(1:i64 code)
    ListDictDataResponse ListDictDataByType(1:string type)
    CreateDictDataResponse CreateDictData(1:CreateDictDataRequest req)
    UpdateDictDataReponse UpdateDictData(1:UpdateDictDataRequest req)
    DeleteDictDataReponse DeleteDictData(1:DeleteDictDataRequest req)

    // dict type service 
    ListDictTypesResponse ListDictTypes(1:ListDictTypesRequest req)
    ExportDictTypeResponse ExportDictType(1:ExportDictTypeRequest req)
    DictTypeResponse GetDictTypeById(1:i64 id)
    CreateDictTypeResponse CreateDictType(1:CreateDictTypeRequest req)
    UpdateDictTypeResponse UpdateDictType(1:UpdateDictTypeRequest req)
    DeleteDictTypeResponse DeleteDictType(1:DeleteDictTypeRequest req)
    RefreshDictTypeResponse RefreshDictType()

    // sys logininfo service
    ListSysLogininfosResponse ListSysLogininfos(1:ListSysLogininfosRequest req)
    ExportLogininfoResponse ExportLogininfo(1:ExportLogininfoRequest req)
    ListSysLogininfosResponse ListSysLogininfosById(1:ListSysLogininfosByIdRequest req)
    CleanSysLogininfoResponse Clean()
    UnlockResponse UnlockByUserName(1:string username)
    CreateSysLogininfoResponse CreateSysLogininfo(1:CreateSysLogininfoRequest req)

    // sys menu service
    ListSysMenusResponse ListSysMenus(1:ListSysMenusRequest req)
    SysMenuResponse GetSysMenuById(1:i64 id)
    ListSysMenusResponse ListTreeMenu(1:ListTreeMenuRequest req)
    RoleMenuResponse ListTreeMenuByRoleid(i64: id)
    CreateMenuResponse CreateMenu(1:CreateMenuRequest req)
    UpdateMemuResponse UpdateMenu(1:UpdateMenuRequest req)
    DeleteMenuResponse DeleteMenu(1:DeleteMenuRequest req)
    RoutersResonse GetRouters()

    // notice service 
    ListSysNoticesResponse ListSysNotices(1:ListSysNoticesRequest req)
    SysNoticeResponse GetSysNoticeById(1:i64 id)
    CreateSysNoticeResponse CreateSysNotice(1:CreateSysNoticeRequest req)
    DeleteSysNoticeRepsonse DeleteSysNotice(1:DeleteSysNoticeRequest req)
    UpdateSysNoticeResponse UpdateSysNotice(1:UpdateSysNoticeRequest req)

    // opelog service
    ListSysOperLogsResponse ListSysOperLogs(1:ListSysOperLogsRequest req)
    ExportSysOperLogResponse ExportSysOperLog(1:ExportSysOperLogRequest req)
    DeleteSysOperLogResponse DeleteSysOperLog(1:DeleteSysOperLogRequest req)
    CleanSysOperLogResponse Clean()
    CreateSysOperLogResponse CreateSysOperLog(1:CreateSysOperLogRequest req)

    // syspost service
    ListSysPostsResponse ListSysPosts(1:ListSysPostsRequest req)
    ExportSysPostResponse ExportSysPost(1:ExportSysPostRequest req)
    SysPostResponse GetSysPostById(1:i64:id)
    CreateSysPostResponse CreateSysPost(1:CreateSysPostRequest req)
    UpdateSysPostResponse UpdateSysPost(1:UpdateSysPostRequest req)
    DeleteSysPostResponse DeleteSysPost(1:DeleteSysPostRequest req)
    ListSysPostResponse ListPostOption()

    // profile service
    ProfileResponse Profile()
    UpdateProfileResponse UpdateProfile(1:UpdateProfileRequest req)  
    UpdatePasswordResponse UpdatePassword(1:UpdatePasswordRequest req)
    UpdateUserAvatarResponse UpdateUserAvatar(1:UpdateUserAvatarRequest req)

    // role service
    ListSysRolesResponse ListSysRole(1:ListSysRolesRequest req)
    ExportSysRoleResponse ExportSysRole(1:ExportSysRoleRequest req)
    SysRoleResponse GetSysRoleByid(1:i64 id)
    CreateSysRoleResponse CreateSysRole(1:CreateSysRoleRequest req)
    UpdateSysRoleResponse UpdateSysRole(1:UpdateSysRoleRequest req)
    DataScopeResponse DataScope(1:DataScopeRequest req)
    ChangeSysRoleStatusResponse ChangeSysRoleStatus(1:ChangeSysRoleStatusRequest req)
    DeleteSysRoleResponse DeleteSysRole(1:DeleteSysRoleRequest req)
    ListSysRolesResponse ListRoleOption()
    ListSysRolesResponse AllocatedList(1:AllocatedListRequest req)
    ListSysRolesResponse UnallocatedList(1:UnallocatedListRequest req)
    CancelAuthUserResponse CancelAuthUser(1:CancelAuthUserRequest req)
    CancelAuthUserAllResponse CancelAuthUserAll(1:CancelAuthUserRequest req)
    SelectAuthUserAllResponse SelectAuthUserAll(1:SelectAuthUserAllRequest req)
    ListSysRolesResponse GetRoleDeptsById(1:i64 id)

    // user service
    ListSysUsersResponse ListSysUsers(1:ListSysUsersRequest req)
    ExportSysUserResponse ExportSysUser(1:ExportSysUserRequest req)
    ImportUserDataResponse ImportUserData(1:ImportUserDataRequest req)
    ImportTemplateResponse ImportTemplate()
    UserInfoResponse GetUserInfoByName(1:string name)
    RegisterSysUserResponse RegisterSysUser(1:RegisterSysUserRequest req)
    CurrentUserInfoResponse CurrentUserInfo()
    UserInfoResponse GetUserInfoById(1:i64 id)
    CreateSysUserResponse CreateSysUser(1:CreateSysUserRequest req)
    UpdateSysUserResponse UpdateSysUser(1:UpdateSysUserRequest req)
    DeleteSysUserResposne DeleteSysUser(1:DeleteSysUserRequest req)
    ResetPasswordResponse ResetSysUser(1:ResetPasswordRequest req)
    ChangeSysUserStatusResponse ChangeSysUserStatus(1:ChangeSysUserStatus req)
    AuthRoleInfoResponse GetAuthRoleById(1:i64 id)
    AuthRoleResponse AuthRole(1:AuthRoleRequest req) // 给其他人授权
    ListDeptsTreeResponse ListDeptsTree(1:ListDeptsTreeRequest req)

    // online user service
    ListSysUserOnline ListSysUserOnlines(1:ListSysUserOnlinesRequest req)
    ForceLogoutResponse ForceLogout(1:ForceLogoutRequest req)
}

struct BaseResp {
    1:i64 code 
    2:string message
}



