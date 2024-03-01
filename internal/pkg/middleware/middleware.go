package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// 注册所有支持的通用中间件
// 还差权限中间件
var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]app.HandlerFunc {
	return map[string]app.HandlerFunc{
		"cache":     Cache(),
		"recovery":  Recovery(),
		"cors":      Cors(),
		"requestid": RequestID(),
		"accesslog": AccessLog(),
	}
}

func Get(mws ...string) (res []app.HandlerFunc) {
	res = make([]app.HandlerFunc, 0, len(mws))
	for _, mw := range mws {
		if a, ok := Middlewares[mw]; ok {
			res = append(res, a)
		}
	}
	return
}
