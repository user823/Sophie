# 系统管理接口

## 用户管理
GET v1/user/list
POST v1/user/export
POST v1/user/importData
POST v1/user/importTemplate
GET v1/user/info/{username}
POST v1/user/register （用户注册）
GET v1/user/getInfo
GET v1/user/
GET v1/user/{userId}
POST v1/user/add
POST v1/user （新增用户）
PUT v1/user
DELETE v1/user/{userIds}
PUT v1/user/resetPwd
PUT v1/user/changeStatus
GET v1/user/authRole/{userId}
PUT v1/user/authRole
GET v1/user/deptTree

### 个人主页
GET v1/user/profile
PUT v1/user/profile
PUT v1/user/profile/updatedPwd
POSE v1/user/profile/avatar


## 角色管理
GET v1/role/list
POST v1/role/export
GET v1/role/{roleId}
POST v1/role
PUT v1/role
PUT v1/role/dataScope
PUT v1/role/changeStatus
DELETE v1/role/{roleIds}
GET v1/role/optionselect
GET v1/role/authUser/allocatedList
GET v1/role/authUser/unallocatedList
PUT v1/role/authUser/cancel
PUT v1/role/authUser/cancelAll
PUT v1/role/authUser/selectAll
GET v1/role/deptTree/{roleId}

## 岗位管理
GET v1/post/list
POST v1/post/export
GET v1/post/{postId}
POST v1/post
PUT v1/post
DELETE v1/post/{postIds}
GET v1/post/optionselect

## 操作日志
GET v1/operlog/list
POST v1/operlog/export
DELETE v1/operlog/{operIds}
DELETE v1/operlog/clean
POST v1/operlog

## 公告管理
GET v1/notice/list
GET v1/notice/{noticeId}
POST v1/notice
PUT v1/notice
DELETE v1/notice/{noticeIds}

## 菜单管理
GET v1/menu/list
GET v1/menu/{menuId}
GET v1/menu/treeselect
GET v1/menu/roleMenuTreeselect/{roleId}
POST v1/menu
PUT v1/menu
DELETE v1/menu/{menuId}
GET v1/menu/getRouters

## 登录日志
GET v1/logininfor/list
POST v1/logininfor/export
DELETE v1/logininfor/{inforIds}
DELETE v1/logininfor/clean
GET v1/logininfor/unlock/{userName}
POST v1/logininfor

## 字典管理
**字典类型**
GET v1/dict/type/list
POST v1/dict/type/export
GET v1/dict/type/{dictId}
POST v1/dict/type
PUT v1/dict/type
DELETE v1/dict/type/{dictIds}
DELETE v1/dict/type/refreshCache
GET v1/dict/type/optionselect

**字典数据**
GET v1/dict/data/list
POST v1/dict/data/export
GET v1/dict/data/{dictCode}
GET v1/dict/data/type/{dictType}
POST v1/dict/data
PUT v1/dict/data
DELETE v1/dict/data/{dictCodes}

## 部门管理
GET v1/dept/list
GET v1/dept/list/exclude/{deptId}
GET v1/dept/{deptId}
POST v1/dept
PUT v1/dept
DELETE v1/dept/{deptId}

## 配置管理
GET v1/config/list
POST v1/config/export
GET v1/config/{configId}
GET v1/config/configKey/{configKey}
POST v1/config
PUT v1/config
DELETE v1/config/{configIds}
DELETE v1/config/refreshCache

## 在线用户服务
GET v1/online/list
DELETE v1/online/{tokenId}