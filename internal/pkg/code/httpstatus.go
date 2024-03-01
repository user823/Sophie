package code

// 使用Http状态码定义系统业务码
const (
	SUCCESS          = 200
	CREATED          = 201
	ACCEPTED         = 202
	NO_CONTENT       = 204
	MOVED_PERM       = 301
	SEE_OTHER        = 303
	NOT_MODIFIED     = 304
	BAD_REQUEST      = 400
	UNAUTHRIZED      = 401
	FORBIDDEN        = 403
	NOT_FOUND        = 404
	BAD_METHOD       = 405
	CONFLICT         = 409
	UNSUPPORTED_TYPE = 415
	ERROR            = 500
	NOT_IMPLEMENTED  = 501
	WARN             = 601
)

// 注册业务码和状态码及消息
func init() {
	register(SUCCESS, SUCCESS, "操作成功")
	register(CREATED, CREATED, "对象创建成功")
	register(ACCEPTED, ACCEPTED, "请求已经被接受")
	register(NO_CONTENT, NO_CONTENT, "操作已经执行成功，但是没有返回数据")
	register(MOVED_PERM, MOVED_PERM, "资源已被移除")
	register(SEE_OTHER, SEE_OTHER, "重定向")
	register(NOT_MODIFIED, NOT_MODIFIED, "资源没有被修改")
	register(BAD_REQUEST, BAD_REQUEST, "参数列表错误（缺少，格式不匹配）")
	register(UNAUTHRIZED, UNAUTHRIZED, "未授权")
	register(FORBIDDEN, FORBIDDEN, "访问受限，授权过期")
	register(NOT_FOUND, NOT_FOUND, "资源，服务未找到")
	register(BAD_METHOD, BAD_METHOD, "不允许的http方法")
	register(CONFLICT, CONFLICT, "资源冲突，或者资源被锁")
	register(UNSUPPORTED_TYPE, UNSUPPORTED_TYPE, "不支持的数据，媒体类型")
	register(ERROR, ERROR, "系统内部错误")
	register(NOT_IMPLEMENTED, NOT_IMPLEMENTED, "接口未实现")
	register(WARN, WARN, "系统警告消息")
}
