namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/system/v1 -service SystemService -gen-path api/thrift/system  system.thrift

service SystemService {
    // config service
    ListConfigsResponse ListConfigs(1:ListConfigsRequest req)
    ExportConfigResponse ExportConfig(1:ExportConfigRequest req)
    ConfigResponse GetConfigById(1:i64 id)
    BaseResp GetConfigByKey(1:string key)
    BaseResp CreateConfig(1:CreateConfigRequest req)
    BaseResp UpdateConfig(1:UpdateConfigReqeust req)
    BaseResp DeleteConfig(1:DeleteConfigReqeust req)
    BaseResp RefreshConfig()

    // dept service
    ListDeptsResponse ListDepts(1:ListDeptsRequest req)
    ListDeptsResponse ListDeptsExcludeChild(1:i64 id)
    DeptResponse GetDeptById(1:GetDeptByIdReq req)
    BaseResp CreateDept(1:CreateDeptRequest req)
    BaseResp UpdateDept(1:UpdateDeptRequest req)
    BaseResp DeleteDept(1:DeleteDeptRequest req)

    // dict data service
    ListDictDatasResponse ListDictDatas(1:ListDictDatasRequest req)
    ExportDictDataResponse ExportDictData(1:ExportDictDataRequest req)
    DictDataResponse GetDictDataByCode(1:i64 code)
    ListDictDatasResponse ListDictDataByType(1:string dictType)
    BaseResp CreateDictData(1:CreateDictDataRequest req)
    BaseResp UpdateDictData(1:UpdateDictDataRequest req)
    BaseResp DeleteDictData(1:DeleteDictDataRequest req)

    // dict type service 
    ListDictTypesResponse ListDictTypes(1:ListDictTypesRequest req)
    ExportDictTypeResponse ExportDictType(1:ExportDictTypeRequest req)
    DictTypeResponse GetDictTypeById(1:i64 id)
    BaseResp CreateDictType(1:CreateDictTypeRequest req)
    BaseResp UpdateDictType(1:UpdateDictTypeRequest req)
    BaseResp DeleteDictType(1:DeleteDictTypeRequest req)
    BaseResp RefreshDictType()
    DictTypeOptionSelectResponse DictTypeOptionSelect()

    // sys logininfo service
    ListSysLogininfosResponse ListSysLogininfos(1:ListSysLogininfosRequest req)
    ExportLogininfoResponse ExportLogininfo(1:ExportLogininfoRequest req)
    BaseResp RemoveSysLogininfosById(1:RemoveSysLogininfosByIdRequest req)
    BaseResp LogininfoClean()
    BaseResp UnlockByUserName(1:string username)
    BaseResp CreateSysLogininfo(1:CreateSysLogininfoRequest req)

    // sys menu service
    ListSysMenusResponse ListSysMenus(1:ListSysMenusRequest req)
    SysMenuResponse GetSysMenuById(1:i64 id)
    ListSysMenusResponse ListTreeMenu(1:ListTreeMenuRequest req)
    RoleMenuResponse ListTreeMenuByRoleid(1:ListTreeMenuByRoleidRequest req)
    BaseResp CreateMenu(1:CreateMenuRequest req)
    BaseResp UpdateMenu(1:UpdateMenuRequest req)
    BaseResp DeleteMenu(1:DeleteMenuRequest req)
    RoutersResonse GetRouters(1:GetRoutersRequest req)

    // notice service 
    ListSysNoticesResponse ListSysNotices(1:ListSysNoticesRequest req)
    SysNoticeResponse GetSysNoticeById(1:i64 id)
    BaseResp CreateSysNotice(1:CreateSysNoticeRequest req)
    BaseResp DeleteSysNotice(1:DeleteSysNoticeRequest req)
    BaseResp UpdateSysNotice(1:UpdateSysNoticeRequest req)

    // operlog service
    ListSysOperLogsResponse ListSysOperLogs(1:ListSysOperLogsRequest req)
    ExportSysOperLogResponse ExportSysOperLog(1:ExportSysOperLogRequest req)
    BaseResp DeleteSysOperLog(1:DeleteSysOperLogRequest req)
    BaseResp OperLogClean()
    BaseResp CreateSysOperLog(1:CreateSysOperLogRequest req)

    // syspost service
    ListSysPostsResponse ListSysPosts(1:ListSysPostsRequest req)
    ExportSysPostResponse ExportSysPost(1:ExportSysPostRequest req)
    SysPostResponse GetSysPostById(1:i64 id)
    BaseResp CreateSysPost(1:CreateSysPostRequest req)
    BaseResp UpdateSysPost(1:UpdateSysPostRequest req)
    BaseResp DeleteSysPost(1:DeleteSysPostRequest req)
    PostOptionSelectResponse PostOptionSelect()

    // profile service
    ProfileResponse Profile(1:ProfileRequest req)
    BaseResp UpdateProfile(1:UpdateProfileRequest req)
    BaseResp UpdatePassword(1:UpdatePasswordRequest req)

    // role service
    ListSysRolesResponse ListSysRole(1:ListSysRolesRequest req)
    ExportSysRoleResponse ExportSysRole(1:ExportSysRoleRequest req)
    SysRoleResponse GetSysRoleByid(1:i64 id)
    BaseResp CreateSysRole(1:CreateSysRoleRequest req)
    BaseResp UpdateSysRole(1:UpdateSysRoleRequest req)
    BaseResp DataScope(1:DataScopeRequest req)
    BaseResp ChangeSysRoleStatus(1:ChangeSysRoleStatusRequest req)
    BaseResp DeleteSysRole(1:DeleteSysRoleRequest req)
    ListSysRolesResponse ListRoleOption()
    ListSysUsersResponse AllocatedList(1:AllocatedListRequest req)
    ListSysUsersResponse UnallocatedList(1:UnallocatedListRequest req)
    BaseResp CancelAuthUser(1:CancelAuthUserRequest req)
    BaseResp CancelAuthUserAll(1:CancelAuthUserAllRequest req)
    BaseResp SelectAuthUserAll(1:SelectAuthUserAllRequest req)
    DeptTreeByRoleIdResponse DeptTreeByRoleId(1:i64 id)

    // user service
    ListSysUsersResponse ListSysUsers(1:ListSysUsersRequest req)
    ExportSysUserResponse ExportSysUser(1:ExportSysUserRequest req)
    BaseResp ImportUserData(1:ImportUserDataRequest req)
    // 使用用户名获取当前登录用户权限信息
    UserInfoResponse GetUserInfoByName(1:string name)
    // 使用id获取当前登录用户权限信息
    UserInfoResponse GetUserInfo(1:i64 id)
    RegisterSysUserResponse RegisterSysUser(1:RegisterSysUserRequest req)
    // 设置id = -1 表示不需要额外信息
    UserInfoByIdResponse GetUserInfoById(1:i64 id)
    BaseResp CreateSysUser(1:CreateSysUserRequest req)
    BaseResp UpdateSysUser(1:UpdateSysUserRequest req)
    BaseResp DeleteSysUser(1:DeleteSysUserRequest req)
    BaseResp ResetPassword(1:ResetPasswordRequest req)
    BaseResp ChangeSysUserStatus(1:ChangeSysUserStatus req)
    AuthRoleInfoResponse GetAuthRoleById(1:i64 id)
    BaseResp AuthRole(1:AuthRoleRequest req) // 给其他人授权
    ListDeptsTreeResponse ListDeptsTree(1:ListDeptsTreeRequest req)

    // online user service
    ListSysUserOnline ListSysUserOnlines(1:ListSysUserOnlinesRequest req)
    BaseResp ForceLogout(1:ForceLogoutRequest req)
}

struct BaseResp {
    1:i64 code
    2:string msg
}

struct PageInfo {
    1:i64 pageNum
    2:i64 pageSize
    3:string orderByColumn
    4:string isAsc
}

struct DateRange {
    1:i64 beginTime // 使用毫秒
    2:i64 endTime
}

struct ConfigInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 configId
    8:string configName
    9:string configKey
    10:string configValue
    11:string configType
}

struct ListConfigsRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:ConfigInfo configInfo
    4:LoginUser loginUser
}

struct ListConfigsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<ConfigInfo> rows
}

struct ExportConfigRequest {
    1:PageInfo pageInfo
    2:ConfigInfo configInfo
    3:LoginUser loginUser
}

struct ExportConfigResponse {
    1:BaseResp baseResp
    2:list<ConfigInfo> list
    3:string sheetName
    4:string title
}

struct ConfigResponse {
    1:BaseResp baseResp
    2:ConfigInfo data
}

struct CreateConfigRequest {
    1:ConfigInfo configInfo
    2:LoginUser user
}

struct UpdateConfigReqeust {
    1:ConfigInfo configInfo
    2:LoginUser user
}

struct DeleteConfigReqeust {
    1:list<i64> configIds
    2:LoginUser user
}

struct DeptInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 deptId
    8:i64 parentId
    9:string ancestors
    10:string deptName
    11:i64 orderNum
    12:string leader
    13:string phone
    14:string email
    15:string status
    16:string delFlag
    17:string parentName
    18:list<DeptInfo> children
}

struct ListDeptsRequest {
    1:DeptInfo deptInfo
    2:LoginUser loginUser
}

struct ListDeptsResponse {
    1:BaseResp baseResp
    2:list<DeptInfo> data
}

struct DeptResponse {
    1:BaseResp baseResp
    2:DeptInfo data
}

struct GetDeptByIdReq {
    1:i64 id
    2:LoginUser user
}

struct CreateDeptRequest {
    1:DeptInfo dept
    2:LoginUser user
}

struct UpdateDeptRequest {
    1:DeptInfo dept
    2:LoginUser user
}

struct DeleteDeptRequest {
    1:i64 deptId
    2:LoginUser user
}

struct DictData {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 dictCode
    8:i64 dictSort
    9:string dictLabel
    10:string dictValue
    11:string dictType
    12:string cssClass
    13:string listClass
    14:string isDefault
    15:string status
}

struct ListDictDatasRequest {
    1:PageInfo pageInfo
    2:DictData dictData
    3:LoginUser user
}

struct ListDictDatasResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<DictData> rows
}

struct ExportDictDataRequest {
    1:PageInfo pageInfo
    2:DictData dictData
    3:LoginUser user
}

struct ExportDictDataResponse {
    1:BaseResp baseResp
    2:list<DictData> list
    3:string sheetName
    4:string title
}

struct DictDataResponse {
    1:BaseResp baseResp
    2:DictData dictData
}

struct CreateDictDataRequest {
    1:DictData dictData
    2:LoginUser user
}

struct UpdateDictDataRequest {
    1:DictData dictData
    2:LoginUser user
}

struct DeleteDictDataRequest {
    1:list<i64> dictCodes
    2:LoginUser user
}

struct DictType {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 dictId
    8:string dictName
    9:string dictType
    10:string status
}

struct ListDictTypesRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:DictType dictType
    4:LoginUser user
}

struct ListDictTypesResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<DictType> rows
}

struct ExportDictTypeRequest {
    1:PageInfo pageInfo
    2:DictType dictType
    3:LoginUser user
}

struct ExportDictTypeResponse {
    1:BaseResp baseResp
    2:list<DictType> list
    3:string sheetName
    4:string title
}

struct DictTypeResponse {
    1:BaseResp baseResp
    2:DictType data
}

struct CreateDictTypeRequest {
    1:DictType dictType
    2:LoginUser user
}

struct UpdateDictTypeRequest {
    1:DictType dictType
    2:LoginUser user
}

struct DeleteDictTypeRequest {
    1:list<i64> dictIds
    2:LoginUser user
}

struct DictTypeOptionSelectResponse {
    1:BaseResp baseResp
    2:list<DictType> data
}

struct Logininfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 infoId
    8:string userName
    9:string status
    10:string ipaddr
    11:string msg
    12:i64 accessTime
}

struct ListSysLogininfosRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:Logininfo loginInfo
    4:LoginUser user
}

struct ListSysLogininfosResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<Logininfo> rows
}

struct ExportLogininfoRequest {
    1:PageInfo pageInfo
    2:Logininfo loginInfo
    3:LoginUser user
}

struct ExportLogininfoResponse {
    1:BaseResp baseResp
    2:list<Logininfo> list
    3:string sheetName
    4:string title
}

struct RemoveSysLogininfosByIdRequest {
    1:list<i64> infoIds
    2:LoginUser user
}

struct CreateSysLogininfoRequest {
    1:Logininfo loginInfo
    2:LoginUser user
}

struct MenuInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 menuId
    8:string menuName
    9:string parentName
    10:i64 parentId
    11:i64 orderNum
    12:string path
    13:string component
    14:string query
    15:string isFrame
    16:string isCache
    17:string menuType
    18:string visible
    19:string status
    20:string perms
    21:string icon
    22:list<MenuInfo> children
}

struct ListSysMenusRequest {
    1:MenuInfo menuInfo
    2:LoginUser user
}

struct ListSysMenusResponse {
    1:BaseResp baseResp
    2:list<MenuInfo> data
}

struct SysMenuResponse {
    1:BaseResp baseResp
    2:MenuInfo data
}

struct ListTreeMenuRequest {
    1:MenuInfo menuInfo
    2:LoginUser user
}

struct TreeSelect {
    1:i64 id
    2:string label
    3:list<TreeSelect> children
}

struct RoleMenuResponse {
    1:BaseResp baseResp
    2:list<i64> checkedKeys
    3:list<TreeSelect> menus
}

struct ListTreeMenuByRoleidRequest {
    1:i64 id
    2:LoginUser user
}

struct CreateMenuRequest {
    1:MenuInfo menuInfo
    2:LoginUser user
}

struct UpdateMenuRequest {
    1:MenuInfo menuInfo
    2:LoginUser user
}

struct DeleteMenuRequest {
    1:i64 menuId
    2:LoginUser user
}

struct GetRoutersRequest {
    1:LoginUser user
}

struct GetSysMenuPermsByRoleIdsRequest {
    1:list<i64> roleIds
}

struct SysMenuPermsResponse {
    1:BaseResp baseResp
    2:list<string> Perms
}

struct Meta {
    1:string title
    2:string icon
    3:bool noCache
    4:string link
}

struct RouterInfo {
    1:string name
    2:string path
    3:bool hidden
    4:string redirect
    5:string component
    6:string query
    7:bool alwaysShow
    8:Meta meta
    12:list<RouterInfo> children
}

struct RoutersResonse {
    1:BaseResp baseResp
    2:list<RouterInfo> data
}

struct NoticeInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 noticeId
    8:string noticeTitle
    9:string noticeType
    10:string noticeContent
    11:string status
}

struct ListSysNoticesRequest {
    1:PageInfo pageInfo
    2:NoticeInfo noticeInfo
    3:LoginUser user
}

struct ListSysNoticesResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<NoticeInfo> rows
}

struct SysNoticeResponse {
    1:BaseResp baseResp
    2:NoticeInfo data
}

struct CreateSysNoticeRequest {
    1:NoticeInfo noticeInfo
    2:LoginUser user
}

struct DeleteSysNoticeRequest {
    1:list<i64> noticeIds
    2:LoginUser user
}

struct UpdateSysNoticeRequest {
    1:NoticeInfo noticeInfo
    2:LoginUser user
}

struct OperLog {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 operId
    8:string title
    9:i64 businessType
    10:list<i64> businessTypes
    11:string method
    12:string requestMethod
    13:i64 operatorType
    14:string operName
    15:string deptName
    16:string operUrl
    17:string operIp
    18:string operParam
    19:string jsonResult
    20:string status
    21:string errorMsg
    22:i64 operTime
    23:i64 costTime
}

struct ListSysOperLogsRequest {
    1:PageInfo pageInfo
    2:OperLog operLog
    3:LoginUser user
}

struct ListSysOperLogsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<OperLog> rows
}

struct ExportSysOperLogRequest {
    1:PageInfo pageInfo
    2:OperLog operLog
    3:LoginUser user
}

struct ExportSysOperLogResponse {
    1:BaseResp baseResp
    2:list<OperLog> operLogs
    3:string sheetName
    4:string title
}

struct DeleteSysOperLogRequest {
    1:list<i64> operIds
    2:LoginUser user
}

struct CreateSysOperLogRequest {
    1:OperLog operLog
    2:LoginUser user
}

struct PostInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 postId
    8:string postCode
    9:string postName
    10:i64 postSort
    11:string status
    12:bool flag
}

struct ListSysPostsRequest {
    1:PageInfo pageInfo
    2:PostInfo postInfo
    3:LoginUser user
}

struct ListSysPostsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<PostInfo> rows
}

struct ExportSysPostRequest {
    1:PageInfo pageInfo
    2:PostInfo postInfo
    3:LoginUser user
}

struct ExportSysPostResponse {
    1:BaseResp baseResp
    2:list<PostInfo> list
    3:string sheetName
    4:string title
}

struct SysPostResponse {
    1:BaseResp baseResp
    2:PostInfo postInfo
}

struct CreateSysPostRequest {
    1:PostInfo postInfo
    2:LoginUser user
}

struct UpdateSysPostRequest {
    1:PostInfo postInfo
    2:LoginUser user
}

struct DeleteSysPostRequest {
    1:list<i64> postIds
    2:LoginUser user
}

struct PostOptionSelectResponse {
    1:BaseResp baseResp
    2:list<PostInfo> data
}

struct RoleInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 roleId
    8:string roleName
    9:string roleKey
    10:i64 roleSort
    11:string dataScope
    12:bool menuCheckStrictly
    13:bool deptCheckStrictly
    14:string status
    15:string delFlag
    16:bool flag
    17:list<i64> menuIds
    18:list<i64> deptIds
    19:list<string> permissions
}

struct UserInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:i64 userId
    8:i64 deptId
    9:string userName
    10:string nickName
    11:string email
    12:string phonenumber
    13:string sex
    14:string avatar
    15:string password
    16:string status
    17:string delFlag
    18:string loginIp
    19:i64 loginDate
    20:DeptInfo dept
    21:list<RoleInfo> roles
    22:list<i64> roleIds
    23:list<i64> postIds
    24:i64 roleId
}

struct LoginUser {
    1:UserInfo user
    2:list<string> permissions
    3:list<string> roles
}

struct ProfileResponse {
    1:BaseResp baseResp
    2:UserInfo userInfo
    3:string roleGroup
    4:string postGroup
}

struct ProfileRequest {
    1:LoginUser user
}

struct UpdateProfileRequest {
    1:UserInfo userInfo
    2:LoginUser user
}

struct UpdatePasswordRequest {
    1:string oldPassword
    2:string newPassword
    3:LoginUser user
}

struct ListSysRolesRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:RoleInfo roleInfo
    4:LoginUser user
}

struct ListSysRolesResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<RoleInfo> rows
}

struct ExportSysRoleRequest {
    1:PageInfo pageInfo
    2:RoleInfo roleInfo
    3:LoginUser user
}

struct ExportSysRoleResponse {
    1:BaseResp baseResp
    2:list<RoleInfo> list
    3:string sheetName
    4:string title
}

struct SysRoleResponse {
    1:BaseResp baseResp
    2:RoleInfo data
}

struct CreateSysRoleRequest {
    1:RoleInfo roleInfo
    2:LoginUser user
}

struct UpdateSysRoleRequest {
    1:RoleInfo roleInfo
    2:LoginUser user
}

struct DataScopeRequest {
    1:RoleInfo roleInfo
    2:LoginUser user
}

struct ChangeSysRoleStatusRequest {
    1:RoleInfo roleInfo
    2:LoginUser user
}

struct DeleteSysRoleRequest {
    1:list<i64> roleIds
    2:LoginUser user
}

struct RoleOptionSelectResponse {
    1:BaseResp baseResp
    2:list<RoleInfo> data
}

struct AllocatedListRequest {
    1:PageInfo pageInfo
    2:UserInfo userInfo
    3:LoginUser user
}

struct UnallocatedListRequest {
    1:PageInfo pageInfo
    2:UserInfo userInfo
    3:LoginUser user
}

struct CancelAuthUserRequest {
    1:i64 userId
    2:i64 roleId
    3:LoginUser user
}

struct CancelAuthUserAllRequest {
    1:i64 roleId
    2:list<i64> userIds
    3:LoginUser user
}

struct SelectAuthUserAllRequest {
    1:i64 roleId
    2:list<i64> userIds
    3:LoginUser user
}

struct DeptTreeByRoleIdResponse {
    1:BaseResp baseResp
    2:list<i64> checkedKeys
    3:list<TreeSelect> depts
}

struct ListSysUsersRequest {
    1:PageInfo pageInfo
    2:UserInfo userInfo
    3:LoginUser user
}

struct ListSysUsersResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<UserInfo> rows
}

struct ExportSysUserRequest {
    1:PageInfo pageInfo
    2:UserInfo userInfo
    3:LoginUser user
}

struct ExportSysUserResponse {
    1:BaseResp baseResp
    2:list<UserInfo> list
    3:string sheetName
    4:string title
}

struct ImportUserDataRequest {
    1:list<UserInfo> users
    2:bool isUpdateSupport
    3:string operName
    4:LoginUser user
}

struct UserInfoResponse {
    1:BaseResp baseResp
    2:UserInfo data
    3:list<string> roles
    4:list<string> permissions
}

struct UserInfoByIdResponse {
    1:BaseResp baseResp
    2:list<RoleInfo> roles
    3:list<PostInfo> posts
    4:UserInfo data
    5:list<i64> postIds
    6:list<i64> roleIds
}

struct RegisterSysUserRequest {
    1:UserInfo userInfo
    2:LoginUser user
}

struct RegisterSysUserResponse {
    1:BaseResp baseResp
    2:bool isOk
}

struct CurrentUserInfoResponse {
    1:BaseResp baseResp
    2:UserInfo userInfo
    3:list<string> roles
    4:list<string> permissions
}

struct CreateSysUserRequest {
    1:UserInfo userInfo
    2:LoginUser user
}

struct UpdateSysUserRequest {
    1:UserInfo userInfo
    2:LoginUser user
}

struct DeleteSysUserRequest {
    1:list<i64> userIds
    2:LoginUser user
}

struct ResetPasswordRequest {
    1:UserInfo userInfo
    2:LoginUser user
}

struct ChangeSysUserStatus {
    1:UserInfo userInfo
    2:LoginUser user
}

struct AuthRoleInfoResponse {
    1:BaseResp baseResp
    2:UserInfo user
    3:list<RoleInfo> roles
}

struct AuthRoleRequest {
    1:i64 userId
    2:list<i64> roleIds
    3:LoginUser user
}

struct ListDeptsTreeRequest {
    1:DeptInfo deptInfo
    2:LoginUser user
}

struct ListDeptsTreeResponse {
    1:BaseResp baseResp
    2:list<TreeSelect> data
}

struct ListSysUserOnlinesRequest {
    1:PageInfo pageInfo
    2:string ipaddr
    3:string userName
    4:LoginUser user
}

struct UserOnlineInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
    7:string tokenId
    8:string userName
    9:string ipaddr
    10:string loginLocation
    11:string browser
    12:string os
    13:i64 loginTime
}

struct ListSysUserOnline {
    1:BaseResp baseResp
    2:i64 total
    3:list<UserOnlineInfo> rows
}

struct ForceLogoutRequest {
    1:string tokenId
    2:LoginUser user
}

