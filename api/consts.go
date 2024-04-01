package api

// 通用常量配置
const (
	// utf8字符集
	UTF8 = "UTF-8"
	// GBK 字符集
	GBK = "GBK"
	// www主域
	WWW = "www."
	// http请求
	HTTP = "http://"
	// https请求
	HTTPS = "https://"
	// 成功标记
	SUCCESS = 200
	// 失败标记
	FAIL = 500
	// 登陆成功状态
	LOGIN_SUCCESS_STATUS = "0"
	// 登陆失败状态
	LOGIN_FAIL_STATUS = "1"
	// 登陆成功
	LOGIN_SUCCESS = "Success"
	// 注销
	LOGOUT = "Logout"
	// 注册
	REGISTER = "Register"
	// 登陆失败
	LOGIN_FAIL = "Error"
	// 当前记录起始索引
	PAGE_NUM = "pageNum"
	// 每页显示记录数
	PAGE_SIZE = "pageSize"
	// 排序列
	ORDER_BY_COLUMN = "orderByColumn"
	// 排序方向
	IS_ASC = "ascending"
	// 资源映射路径 前缀
	RESOURCE_PREFIX = "/profile"
	// 内部请求
	INNER = "inner"
	// 登陆用户信息设置
	LOGIN_INFO_KEY = "sysloginInfo"
	// 数据范围 key（用来设置 sql查询时模版字符串）
	DATA_SCOPE = "dataScope"
	// 管理员权限字符
	ALL_PERMISSIONS = "*:*:*"
	// 默认登录有效期（3600s)
	LOGIN_TIMEOUT = 3600
)

var (
	// 自动识别json对象白名单配置（仅允许解析的包名，范围越小越安全）
	JSON_WHITELIST_STR = []string{"Sophie"}
	// 定时任务白名单配置（仅允许访问的包名，其他需要可以自行添加）
	JOB_WHITELIST_STR = []string{"Sophie"}
	// 定时任务违规字符串（根据安全需求配置）
	JOB_ERROR_STR = []string{}
)

// 日志打印相关的环境信息
const (
	// 服务模块
	LOG_SERVICE = "service"
)
