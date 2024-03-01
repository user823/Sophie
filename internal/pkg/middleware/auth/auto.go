package auth

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/errors"
	"strings"
)

type AutoStrategy struct {
	basic AuthStrategy
	jwt   AuthStrategy
}

func NewAutoStrategy(basic, jwt AuthStrategy) *AutoStrategy {
	return &AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}

func (a *AutoStrategy) AuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		operator := AuthOperator{}
		auth := strings.SplitN(c.Request.Header.Get(AUTHENTICATION), " ", 2)
		if len(auth) != 2 {
			core.WriteResponseE(c, errors.WithCodeMessage(nil, code.UNAUTHRIZED, "Authorization header format is wrong. "), nil)
			c.Abort()
			return
		}

		switch auth[0] {
		case BASICPREFIX:
			operator.SetStrategy(a.basic)
		case JWTPREFIX:
			operator.SetStrategy(a.jwt)
		default:
			core.WriteResponseE(c, errors.WithCodeMessage(nil, code.UNAUTHRIZED, "Unrecognized Authorization header. "), nil)
			c.Abort()
			return
		}
		operator.AuthFunc()(ctx, c)
		c.Next(ctx)
	}
}
