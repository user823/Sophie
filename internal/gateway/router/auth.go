package gateway

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/system/v1"
	"github.com/user823/Sophie/internal/gateway/store"
	"github.com/user823/Sophie/internal/pkg/middleware"
	"github.com/user823/Sophie/internal/pkg/middleware/auth"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

type loginInfo struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,password"`
}

func newBasicAuth() auth.AuthStrategy {
	// 存储层写好了之后修改
	return auth.NewBasicStrategy(func(username string, password string) bool {
		// 拉取用户信息
		user, err := store.Client().Users().Get(context.Background(), username, api.GetOptions{})
		if err != nil {
			return false
		}

		if err := auth.Compare(user.Password, password); err != nil {
			return false
		}

		// 登陆后处理逻辑
		return true
	})
}

func newJWTAuth() auth.AuthStrategy {
	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       viper.GetString("jwt.realm"),
		Key:         utils.S2b("jwt.key"),
		Timeout:     viper.GetDuration("jwt.timeout"),
		MaxRefresh:  viper.GetDuration("jwt.timeout"),
		IdentityKey: middleware.UsernameKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*v1.SysUser); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return claims[jwt.IdentityKey]
		},
		Authenticator:    authenticator(),
		Authorizator:     authentizator(),
		Unauthorized:     unauthorized(),
		SigningAlgorithm: "HS256",
		LoginResponse:    loginResponse(),
		LogoutResponse: func(ctx context.Context, c *app.RequestContext, code int) {
			c.JSON(200, nil)
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
	})
	if err != nil {
		log.Fatalf("Service auth initial error: %s, %s", ServiceName, err.Error())
	}
	return auth.NewJWTStrategy(*authMiddleware)
}

func authenticator() func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	var login loginInfo
	var err error
}

func authentizator() func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
	return nil
}

func unauthorized() func(ctx context.Context, c *app.RequestContext, code int, message string) {
	return nil
}

func loginResponse() func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
	return nil
}

func newAutoAuth() auth.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth().(*auth.BasicStrategy), newJWTAuth().(*auth.JWTStrategy))
}
