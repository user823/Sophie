package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/errors"
	"strings"
)

type BasicStrategy struct {
	authenticate func(c *app.RequestContext, username, password string) (interface{}, bool)
}

var _ AuthStrategy = (*BasicStrategy)(nil)

func NewBasicStrategy(authenticate func(c *app.RequestContext, username, password string) (interface{}, bool)) *BasicStrategy {
	return &BasicStrategy{
		authenticate: authenticate,
	}
}

func getUsernameAndPassword(authorization string) (string, string, error) {
	auth := strings.SplitN(authorization, " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", "", fmt.Errorf("Authorization header format is wrong. ")
	}
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", fmt.Errorf("Authorization header format is wrong. ")
	}
	return pair[0], pair[1], nil
}

func (b *BasicStrategy) AuthFunc() app.HandlerFunc {
	realm := "Basic realm=Authorization Required"

	return func(ctx context.Context, c *app.RequestContext) {
		// 格式校验
		username, password, err := getUsernameAndPassword(c.Request.Header.Get(AUTHENTICATION))
		if err != nil {
			c.Header("WWW-Authenticate", realm)
			core.WriteResponseE(c, errors.WithCodeMessage(nil, code.UNAUTHRIZED, err.Error()), nil)
			c.Abort()
			return
		}

		data, ok := b.authenticate(c, username, password)
		if !ok {
			core.WriteResponseE(c, errors.WithCodeMessage(nil, code.UNAUTHRIZED, "Username or password wrong. "), nil)
			c.Abort()
			return
		}

		c.Set(api.LOGIN_INFO_KEY, data)
		c.Next(ctx)
	}
}
