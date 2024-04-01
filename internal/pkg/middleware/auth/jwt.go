package auth

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

type JWTStrategy struct {
	jwt.HertzJWTMiddleware
}

var _ AuthStrategy = &JWTStrategy{}

func NewJWTStrategy(jwt jwt.HertzJWTMiddleware) AuthStrategy {
	return &JWTStrategy{
		jwt,
	}
}

func (j *JWTStrategy) AuthFunc() app.HandlerFunc {
	return j.MiddlewareFunc()
}

func (j *JWTStrategy) ExtractToken() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, err := j.GetClaimsFromJWT(ctx, c)
		if err == nil {
			c.Set("JWT_PAYLOAD", claims)
		}
		c.Next(ctx)
	}
}
