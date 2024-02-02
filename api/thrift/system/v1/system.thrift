namespace go v1

// code-gen
// kitex -module github.com/user823/Sophie -I api/thrift/system/v1 -service SystemService -gen-path api/thrift/system  system.thrift

service SystemService {
    // config service
    ListConfigsResponse ListConfigs(1:ListConfigsRequest req)
    ExportConfigResponse ExportConfig(1:ExportConfigRequest req)
    ConfigResponse GetConfigById(1:i64 id)
    ConfigResponse GetConfigByKey(1:string key)
    BaseResp CreateConfig(1:CreateConfigRequest req)
    BaseResp UpdateConfig(1:UpdateConfigReqeust req)
    BaseResp DeleteConfig(1:DeleteConfigReqeust req)
    BaseResp RefreshConfig()

    // dept service
    ListDeptsResponse ListDepts(1:ListDeptsRequest req)
    ListDeptsResponse ListDeptsExcludeChild(1:i64 id)
    DeptResponse GetDeptById(1:i64 id)
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
    RoleMenuResponse ListTreeMenuByRoleid(1:i64 id)
    BaseResp CreateMenu(1:CreateMenuRequest req)
    BaseResp UpdateMenu(1:UpdateMenuRequest req)
    BaseResp DeleteMenu(1:DeleteMenuRequest req)
    RoutersResonse GetRouters()

    SysMenuPermsResponse GetSysMenuPermsByRoleIds(1:GetSysMenuPermsByRoleIdsRequest req)

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
    ProfileResponse Profile()
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
    ListSysRolesResponse AllocatedList(1:AllocatedListRequest req)
    ListSysRolesResponse UnallocatedList(1:UnallocatedListRequest req)
    BaseResp CancelAuthUser(1:CancelAuthUserRequest req)
    BaseResp CancelAuthUserAll(1:CancelAuthUserAllRequest req)
    BaseResp SelectAuthUserAll(1:SelectAuthUserAllRequest req)
    DeptTreeByRoleIdResponse DeptTreeByRoleId(1:i64 id)

    ListSysRolesResponse GetSysRoleByUser(1:i64 id)


    // user service
    ListSysUsersResponse ListSysUsers(1:ListSysUsersRequest req)
    ExportSysUserResponse ExportSysUser(1:ExportSysUserRequest req)
    BaseResp ImportUserData(1:ImportUserDataRequest req)
    UserInfoResponse GetUserInfoByName(1:string name)
    RegisterSysUserResponse RegisterSysUser(1:RegisterSysUserRequest req)
    UserInfoResponse GetUserInfoById(1:i64 id)
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
}

struct DateRange {
    1:i64 beginTime // 使用毫秒
    2:i64 endTime
}

struct BaseInfo {
    1:string createBy
    2:i64 createTime
    3:string updateBy
    4:i64 updateTime
    5:string remark
    6:map<string,string> params
}

struct ConfigInfo {
    1:BaseInfo baseInfo
    2:i64 configId
    3:optional string configName
    4:optional string configKey
    5:optional string configValue
    6:optional string configType
}

struct ListConfigsRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:ConfigInfo configInfo
}

struct ListConfigsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<ConfigInfo> rows
}

struct ExportConfigRequest {
    1:PageInfo pageInfo
    2:ConfigInfo configInfo
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
}

struct UpdateConfigReqeust {
    1:ConfigInfo configInfo
}

struct DeleteConfigReqeust {
    1:list<i64> configIds
}

struct DeptInfo {
    1:BaseInfo baseInfo
    2:i64 deptId
    3:optional i64 parentId
    4:optional string ancestors
    5:optional string deptName
    6:optional i64 orderNum
    7:optional string leader
    8:optional string phone
    9:optional string email
    10:optional string status
    11:optional string delFlag
    12:optional string parentName
    13:optional list<DeptInfo> children
}

struct ListDeptsRequest {
    1:string deptName
    2:string status
}

struct ListDeptsResponse {
    1:BaseResp baseResp
    2:list<DeptInfo> data
}

struct DeptResponse {
    1:BaseResp baseResp
    2:DeptInfo data
}

struct CreateDeptRequest {
    1:DeptInfo dept
}

struct UpdateDeptRequest {
    1:DeptInfo dept
}

struct DeleteDeptRequest {
    1:i64 deptId
}

struct DictData {
    1:BaseInfo baseInfo
    2:i64 dictCode
    3:optional i64 dictSort
    4:optional string dictLabel
    5:optional string dictValue
    6:optional string dictType
    7:optional string cssClass
    8:optional string listClass
    9:optional string isDefault
    10:optional string status
}

struct ListDictDatasRequest {
    1:PageInfo pageInfo
    2:DictData dictData
}

struct ListDictDatasResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<DictData> rows
}

struct ExportDictDataRequest {
    1:PageInfo pageInfo
    2:DictData dictData
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
}

struct UpdateDictDataRequest {
    1:DictData dictData
}

struct DeleteDictDataRequest {
    1:DictData dictData
}

struct DictType {
    1:BaseInfo baseInfo
    2:i64 dictId
    3:optional string dictName
    4:optional string dictType
    5:optional string status
}

struct ListDictTypesRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:DictType dictType
}

struct ListDictTypesResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<DictType> rows
}

struct ExportDictTypeRequest {
    1:PageInfo pageInfo
    2:DictType dictType
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
}

struct UpdateDictTypeRequest {
    1:DictType dictType
}

struct DeleteDictTypeRequest {
    1:list<i64> dictIds
}

struct DictTypeOptionSelectResponse {
    1:BaseResp baseResp
    2:list<DictType> data
}

struct Logininfo {
    1:BaseInfo baseInfo
    2:i64 infoId
    3:optional string userName
    4:optional string status
    5:optional string ipaddr
    6:optional string msg
    7:optional i64 accessTime
}

struct ListSysLogininfosRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:Logininfo loginInfo
}

struct ListSysLogininfosResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<Logininfo> rows
}

struct ExportLogininfoRequest {
    1:PageInfo pageInfo
    2:Logininfo loginInfo
}

struct ExportLogininfoResponse {
    1:BaseResp baseResp
    2:list<Logininfo> list
    3:string sheetName
    4:string title
}

struct RemoveSysLogininfosByIdRequest {
    1:list<i64> infoIds
}

struct CreateSysLogininfoRequest {
    1:Logininfo loginInfo
}

struct MenuInfo {
    1:BaseInfo baseInfo
    2:i64 menuId
    3:optional string menuName
    4:optional string parentName
    5:optional i64 parentId
    6:optional i64 orderNum
    7:optional string path
    8:optional string component
    9:optional string query
    10:optional string isFrame
    11:optional string isCache
    12:optional string menuType
    13:optional string visible
    14:optional string status
    15:optional string perms
    16:optional string icon
    17:optional list<MenuInfo> children
}

struct ListSysMenusRequest {
    1:MenuInfo menuInfo
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

struct CreateMenuRequest {
    1:MenuInfo menuInfo
}

struct UpdateMenuRequest {
    1:MenuInfo menuInfo
}

struct DeleteMenuRequest {
    1:i64 menuId
}

struct GetSysMenuPermsByRoleIdsRequest {
    1:list<i64> roleIds
}

struct SysMenuPermsResponse {
    1:BaseResp baseResp
    2:list<string> Perms
}

struct RouterInfo {
    1:string name
    2:string path
    3:bool hidden
    4:string redirect
    5:string component
    6:string query
    7:bool alwaysShow
    8:string title
    9:string icon
    10:bool noCache
    11:string link
    12:list<RouterInfo> children
}

struct RoutersResonse {
    1:BaseResp baseResp
    2:list<RouterInfo> data
}

struct NoticeInfo {
    1:BaseInfo baseInfo
    2:i64 noticeId
    3:optional string noticeTitle
    4:optional string noticeType
    5:optional string noticeContent
    6:optional string status
}

struct ListSysNoticesRequest {
    1:PageInfo pageInfo
    2:NoticeInfo noticeInfo
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
}

struct DeleteSysNoticeRequest {
    1:list<i64> noticeIds
}

struct UpdateSysNoticeRequest {
    1:NoticeInfo noticeInfo
}

struct OperLog {
    1:BaseInfo baseInfo
    2:i64 operId
    3:optional string title
    4:optional i64 businessType
    5:optional list<i64> businessTypes
    6:optional string method
    7:optional string requestMethod
    8:optional i64 operatorType
    9:optional string operName
    10:optional string deptName
    11:optional string operUrl
    12:optional string operIp
    13:optional string operParam
    14:optional string jsonResult
    15:optional i64 status
    16:optional string errorMsg
    17:optional i64 operTime
    18:optional i64 costTime
}

struct ListSysOperLogsRequest {
    1:PageInfo pageInfo
    2:OperLog operLog
}

struct ListSysOperLogsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<OperLog> rows
}

struct ExportSysOperLogRequest {
    1:PageInfo pageInfo
    2:OperLog operLog
}

struct ExportSysOperLogResponse {
    1:BaseResp baseResp
    2:list<OperLog> operLogs
    3:string sheetName
    4:string title
}

struct DeleteSysOperLogRequest {
    1:list<i64> operIds
}

struct CreateSysOperLogRequest {
    1:OperLog operLog
}

struct PostInfo {
    1:BaseInfo baseInfo
    2:i64 poseId
    3:optional string poseCode
    4:optional string poseName
    5:optional i64 postSort
    6:optional string status
    7:optional bool flag
}

struct ListSysPostsRequest {
    1:PageInfo pageInfo
    2:PostInfo postInfo
}

struct ListSysPostsResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<PostInfo> rows
}

struct ExportSysPostRequest {
    1:PageInfo pageInfo
    2:PostInfo postInfo
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
}

struct UpdateSysPostRequest {
    1:PostInfo postInfo
}

struct DeleteSysPostRequest {
    1:list<i64> postIds
}

struct PostOptionSelectResponse {
    1:BaseResp baseResp
    2:list<PostInfo> data
}

struct RoleInfo {
    1:BaseInfo baseInfo
    2:i64 roleId
    3:optional string roleName
    4:optional string roleKey
    5:optional i64 roleSort
    6:optional string dataScope
    7:optional bool menuCheckStrictly
    8:optional bool deptCheckStrictly
    9:optional string status
    10:optional string delFlag
    11:optional bool flag
    12:optional list<i64> menuIds
    13:optional list<i64> deptIds
    14:optional list<string> permissions
}

struct UserInfo {
    1:BaseInfo baseInfo
    2:i64 userId
    3:optional i64 deptId
    4:optional string userName
    5:optional string nickName
    6:optional string email
    7:optional string phonenumber
    8:optional string sex
    9:optional string avatar
    10:optional string password
    11:optional string status
    12:optional string delFlag
    13:optional string loginIp
    14:optional i64 loginDate
    15:optional DeptInfo dept
    16:optional list<RoleInfo> roles
    17:optional list<i64> roleIds
    18:optional list<i64> postIds
    19:optional i64 roleId
}

struct ProfileResponse {
    1:BaseResp baseResp
    2:UserInfo userInfo
    3:string roleGroup
    4:string postGroup
}

struct UpdateProfileRequest {
    1:UserInfo userInfo
}

struct UpdatePasswordRequest {
    1:string oldPassword
    2:string newPassword
    3:string confirmPassword
}

struct ListSysRolesRequest {
    1:PageInfo pageInfo
    2:DateRange dateRange
    3:RoleInfo roleInfo
}

struct ListSysRolesResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<RoleInfo> rows
}

struct ExportSysRoleRequest {
    1:PageInfo pageInfo
    2:RoleInfo roleInfo
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
}

struct UpdateSysRoleRequest {
    1:RoleInfo roleInfo
}

struct DataScopeRequest {
    1:RoleInfo roleInfo
}

struct ChangeSysRoleStatusRequest {
    1:RoleInfo roleInfo
}

struct DeleteSysRoleRequest {
    1:list<i64> roleIds
}

struct RoleOptionSelectResponse {
    1:BaseResp baseResp
    2:list<RoleInfo> data
}

struct AllocatedListRequest {
    1:UserInfo userInfo
}

struct UnallocatedListRequest {
    1:UserInfo userInfo
}

struct CancelAuthUserRequest {
    1:UserInfo userInfo
}

struct CancelAuthUserAllRequest {
    1:i64 roleId
    2:list<i64> userIds
}

struct SelectAuthUserAllRequest {
    1:i64 roleId
    2:list<i64> userIds
}

struct DeptTreeByRoleIdResponse {
    1:BaseResp baseResp
    2:list<i64> checkedKeys
    3:list<TreeSelect> depts
}

struct ListSysUsersRequest {
    1:PageInfo pageInfo
    2:UserInfo userInfo
}

struct ListSysUsersResponse {
    1:BaseResp baseResp
    2:i64 total
    3:list<UserInfo> rows
}

struct ExportSysUserRequest {
    1:PageInfo pageInfo
    2:UserInfo userInfo
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
}

struct UserInfoResponse {
    1:BaseResp baseResp
    2:UserInfo data
}

struct RegisterSysUserRequest {
    1:UserInfo userInfo
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
}

struct UpdateSysUserRequest {
    1:UserInfo userInfo
}

struct DeleteSysUserRequest {
    1:list<i64> userIds
}

struct ResetPasswordRequest {
    1:UserInfo userInfo
}

struct ChangeSysUserStatus {
    1:UserInfo userInfo
}

struct AuthRoleInfoResponse {
    1:BaseResp baseResp
    2:UserInfo user
    3:list<RoleInfo> roles
}

struct AuthRoleRequest {
    1:i64 userId
    2:list<i64> roleIds
}

struct ListDeptsTreeRequest {
    1:DeptInfo deptInfo
}

struct ListDeptsTreeResponse {
    1:BaseResp baseResp
    2:list<TreeSelect> data
}

struct ListSysUserOnlinesRequest {
    1:PageInfo pageInfo
    2:string ipaddr
    3:string userName
}

struct UserOnlineInfo {
    1:BaseInfo baseInfo
    2:string tokenId
    3:optional string userName
    4:optional string ipaddr
    5:optional string loginLocation
    6:optional string browser
    7:optional string os
    8:optional i64 loginTime
}

struct ListSysUserOnline {
    1:BaseResp baseResp
    2:i64 total
    3:list<UserOnlineInfo> rows
}

struct ForceLogoutRequest {
    1:string tokenId
}

