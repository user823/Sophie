package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"github.com/hertz-contrib/requestid"
)

const (
	XRequestIDKey = "X-Request-ID"
)

func RequestID() app.HandlerFunc {
	return requestid.New(
		requestid.WithGenerator(func(ctx context.Context, c *app.RequestContext) string {
			return uuid.NewString()
		}),
		requestid.WithCustomHeaderStrKey(XRequestIDKey),
	)
}
