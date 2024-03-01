package auth

import (
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
