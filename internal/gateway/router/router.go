package gateway

import "github.com/user823/Sophie/pkg/log"

func (s *GateWayServer) initRouter() {
	// 安装通用中间件
	if s.insecureServer == nil {
		log.Panicf("Insecure engine has not been prepared already: %s", ServiceName)
	}

	s.insecureServer.Use(s.config.Middlewares...)

	// 安装路由
	jwtStrategy, _ :=
}

