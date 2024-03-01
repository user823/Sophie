package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/user823/Sophie/pkg/log"
)

func AccessLog() app.HandlerFunc {
	logFn := func(ctx context.Context, format string, v ...interface{}) {
		log.Infof(format, v...)
	}
	format := "[${time}] ${status} - ${latency} ${method} ${path} ${body}"
	return accesslog.New(accesslog.WithAccessLogFunc(logFn), accesslog.WithFormat(format))
}
