package kv

// 缓存常量信息
const (
	// 缓存有效期， 默认180(分钟）
	EXPIRATION int64 = 180 * 60
	// 缓存刷新时间，默认60（分钟）
	REFRESH_TIME = 60 * 60
	// 密码最大错误次数
	PASSWORD_MAX_RETRY_COUNT = 5
	// 密码锁定时间，默认10（分钟）
	PASSWORD_LOCK_TIME = 10
	// 权限缓存前缀
	LOGIN_TOKEN_KEY = "sophie-login_tokens-"
	// 验证码 redis key
	CAPTHA_CODE_KEY = "sohie-captcha_codes-"
	// 验证码有效期
	CAPTHA_CODE_KEY_VALID = 5 * 60
	// 参数管理 cache key
	SYS_CONFIG_KEY = "sophie-config-"
	// 字典管理 cache key
	SYS_DICT_TYPE = "sophie-dict-"
	// 登陆账户密码错误次数 redis key
	PWD_ERR_CNT_KEY = "sophie-pwd_err_cnt-"
	// 登陆IP 黑名单cache key
	SYS_LOGIN_BLACKIPLIST = SYS_CONFIG_KEY + "loginIP_blacklist-"
	// 登录用户cache key
	SYS_LOGIN_USER = "sophie-loginuser-"
	// 登录用户cache_id_key
	SYS_LOGIN_USER_IDS = "ids"
)
