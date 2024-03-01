package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/log"
)

func Recovery() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(RecoveryHandler))
}

func RecoveryHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	log.Errorf("[Recovery] err=%v\nstack=%s", err, stack)
	log.Infof("Client: %s", c.Request.Header.UserAgent())
	c.Abort()
	core.Fail(c, "系统内部错误，请重试", nil)
}
