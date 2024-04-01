package v1

// 系统服务名称
const (
	ServiceName = "Sophie System"
)

// 用户常量信息
const (
	// 管理员ID
	ROOT_ID int64 = 1
	// 系统用户标识
	ROOT = "SYS_USER"
	// 正常状态
	USERNORMAL = "0"
	// 停用状态
	USERDISABLE = "1"
	// 删除状态
	USERDELETED = "2"
	// 部门正常状态
	DEPTNORMAL = "0"
	// 部门停用状态
	DEPTDISABLE = "1"
	// 字典正常状态
	DICTNORMAL = "0"
	// 是否为系统默认（是）
	YES = "Y"
	// 是否为菜单外链（是）
	YES_FRAME = "0"
	// 是否为菜单外链（否）
	NO_FRAME = "1"
	// 菜单类型（目录）
	TYPE_DIR = "M"
	// 菜单类型（菜单）
	TYPE_MENU = "C"
	// 菜单类型（按钮）
	TYPE_BUTTON = "F"
	// Layout 组件标识
	LAYOUT = "Layout"
	// ParentView 组件标识
	PARENT_VIEW = "ParentView"
	// InnerLink 组件标识
	INNER_LINK = "InnerLink"
	// 校验是否唯一的返回标识
	UNIQUE     = true
	NOT_UNIQUE = false
	// 用户名长度限制
	USERNAME_MIN_LENGTH = 2
	USERNAME_MAX_LENGTH = 20
	// 密码长度限制
	PASSWORDMINLENGTH = 5
	PASSWORDMAXLENGTH = 20
	// 默认头像
	AVATAR_URL = "https://picsum.photos/200"
)

// 操作相关信息
const (
	// 动作类型
	BUSINESSTYPE_NULL int64 = iota
	BUSINESSTYPE_INSERT
	BUSINESSTYPE_UPDATE
	BUSINESSTYPE_DELETE
	BUSINESSTYPE_GRANT
	BUSINESSTYPE_EXPORT
	BUSINESSTYPE_IMPORT
	BUSINESSTYPE_FORCE   // 强退
	BUSINESSTYPE_GENCODE // 生成代码
	BUSINESSTYPE_CLEAN   // 清空数据
	BUSINESSTYPE_OTHER   // 其他类型

	// 操作类型
	OPERATORTYPE_OTHER  int64 = iota - 10 // 其他用户
	OPERATORTYPE_MANAGE                   // 后台用户
	OPERATORTYPR_MOBILE                   // 手机端用户

	// 操作状态
	BUSINESS_SUCCESS = "0"
	BUSINESS_FAIL    = "1"
)
