package router

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/jwt"
	"github.com/user823/Sophie/api"
	v13 "github.com/user823/Sophie/api/domain/gateway/v1"
	"github.com/user823/Sophie/api/domain/system/v1"
	v12 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	code2 "github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/internal/pkg/middleware/auth"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"gorm.io/gorm"
	"time"
)

var (
	fieldTranslations = map[string]string{
		"loginInfo.Username": "用户名",
		"loginInfo.Password": "用户密码",
	}
)

// 登陆参数绑定与校验
type loginInfo struct {
	Username string `json:"username" form:"username" vd:"required,min=2,max=20"`
	Password string `json:"password" form:"password" vd:"required,min=5,max=20"`
}

func newBasicAuth() auth.AuthStrategy {
	return auth.NewBasicStrategy(func(c *app.RequestContext, username, password string) (data interface{}, ok bool) {
		var msg string
		sysLoginInfo := auth.GetLogininfo(c, username)
		defer func() {
			if msg != "" {
				sysLoginInfo.Status = api.LOGIN_FAIL
				sysLoginInfo.Msg = msg
			}
			appendLogininfo(sysLoginInfo)
		}()

		resp, err := rpc.Remoting.GetUserInfoByName(context.Background(), username)
		if err != nil || resp.BaseResp.Code != code2.SUCCESS {
			msg = fmt.Sprintf("登陆用户：%s 出现未知错误，请重试", username)
			return nil, false
		}

		user := v12.UserInfo2SysUser(resp.Data)

		if err = auth.Compare(user.Password, password); err != nil {
			msg = fmt.Sprint("用户名或者密码错误")
			return nil, false
		}

		// 添加环境信息
		return v13.LoginUser{
			Roles:       resp.GetRoles(),
			Permissions: resp.GetPermissions(),
			User:        *user,
		}, true
	})
}

func newJWTAuth(info *JwtInfo) auth.AuthStrategy {
	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:                 info.Realm,
		Key:                   utils.S2b(info.Key),
		Timeout:               info.Timeout,
		MaxRefresh:            info.Timeout,
		IdentityKey:           auth.UsernameKey,
		PayloadFunc:           payloadFunc(),
		Authenticator:         authenticator(),
		Authorizator:          authentizator(),
		Unauthorized:          unauthorized(),
		SigningAlgorithm:      "HS256",
		LoginResponse:         loginResponse(),
		LogoutResponse:        logoutResponse(),
		TokenLookup:           "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:         "Bearer",
		SendCookie:            false,
		HTTPStatusMessageFunc: httpMessageFunc(),
	})
	if err != nil {
		log.Fatalf("Service auth initial error: %s, %s", v13.ServiceName, err.Error())
	}
	return auth.NewJWTStrategy(*authMiddleware)
}

// LoginHandler 中调用
// 授权成功则返回user, role
func authenticator() func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
	return func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
		var login loginInfo
		var err error
		var msgErr error

		// 参数验证/校验失败
		if err = c.BindAndValidate(&login); err != nil {
			return "", errors.New(validLoginErrMsg(err))
		}

		// 设置基本登陆信息
		sysLoginInfo := auth.GetLogininfo(c, login.Username)
		defer func() {
			if msgErr != nil {
				sysLoginInfo.Status = api.LOGIN_FAIL
				sysLoginInfo.Msg = msgErr.Error()
			}
			appendLogininfo(sysLoginInfo)
		}()

		resp, err := rpc.Remoting.GetUserInfoByName(context.Background(), login.Username)
		user := v12.UserInfo2SysUser(resp.Data)
		if err != nil {
			log.Errorf("get user information failed: %s", err.Error())

			// 未找到用户
			if err == gorm.ErrRecordNotFound {
				log.Errorf("get user information failed: %s", err.Error())
				msgErr = fmt.Errorf("登陆用户：%s 不存在", login.Username)
				return "", msgErr
			}

			// 其他错误
			msgErr = fmt.Errorf("登陆用户：%s 出现未知错误，请重试", login.Username)
			return "", msgErr
		}

		if v1.UserStatus["DELETED"].Code == user.DelFlag {
			msgErr = fmt.Errorf("对不起，您的账号：%s 已被删除")
			return "", msgErr
		}

		if v1.UserStatus["DISABLE"].Code == user.Status {
			msgErr = fmt.Errorf("对不起，您的账号：%s 已被停用")
			return "", msgErr
		}

		// 验证密码
		if err = auth.Compare(user.Password, login.Password); err != nil {
			msgErr = fmt.Errorf("用户名或者密码错误")
			return "", msgErr
		}

		// 登陆成功，设置授权上下文信息
		loginUser := v13.LoginUser{
			Roles:       resp.GetRoles(),
			Permissions: resp.GetPermissions(),
			User:        *user,
		}
		c.Set(api.LOGIN_INFO_KEY, loginUser)
		return *user, nil
	}
}

// 登陆成功时设置jwt 载荷信息
func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(v1.SysUser); ok {
			return jwt.MapClaims{
				auth.UsernameKey: v.String(),
			}
		}
		return jwt.MapClaims{}
	}
}

// 设置已授权角色资源访问权限，中间件中调用
// data是IdentityHandler 返回的值
func authentizator() func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
	return func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
		// 授权后要从token中获取载荷信息，结构体被解析成map[string]interface{}
		// 因此在放入载荷之前就把结构体解析成string，这样就可以使用自己的方式解析数据
		if data == nil {
			return false
		}

		if str, ok := data.(string); ok {
			var v v1.SysUser
			v.Unmarshal(str)

			// 用于后续步骤的权限校验
			resp, err := rpc.Remoting.GetUserInfoByName(ctx, v.Username)
			if err != nil || resp.BaseResp.Code != code2.SUCCESS {
				log.Errorf("Get user info error: %s", err.Error())
				return false
			}

			loginUser := v13.LoginUser{
				Roles:       resp.GetRoles(),
				Permissions: resp.GetPermissions(),
				User:        v,
			}
			c.Set(api.LOGIN_INFO_KEY, loginUser)
			return true
		}

		return false
	}
}

// message 是根据authenticator的err构造的信息
func unauthorized() func(ctx context.Context, c *app.RequestContext, code int, message string) {
	return func(ctx context.Context, c *app.RequestContext, code int, message string) {
		core.WriteResponse(c, core.ErrResponse{
			Code:    code,
			Message: message,
		})
	}
}

func httpMessageFunc() func(e error, ctx context.Context, c *app.RequestContext) string {
	return func(e error, ctx context.Context, c *app.RequestContext) string {
		switch e {
		case jwt.ErrExpiredToken:
			return "登录状态已过期"
		default:
			return "身份验证失败"
		}
	}
}

// 设置token 和 过期时间
func loginResponse() func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
	return func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
		data := map[string]interface{}{
			"access_token": message,
			"expires_in":   time,
		}

		// 设置登录状态
		carry, _ := c.Get(api.LOGIN_INFO_KEY)
		v := carry.(v13.LoginUser)
		sysLoginInfo := auth.GetLogininfo(c, v.User.Username)
		record := fmt.Sprintf("%s:%s:%s:%d", v.User.Username, message, sysLoginInfo.Ipaddr, sysLoginInfo.AccessTime)
		redisLoginStatus(message, record, utils.Time2Second(time)*1e9)

		core.WriteResponse(c, core.ErrResponse{
			Code: code,
			Data: data,
		})
	}
}

func logoutResponse() func(ctx context.Context, c *app.RequestContext, code int) {
	return func(ctx context.Context, c *app.RequestContext, code int) {
		claims := jwt.ExtractClaims(ctx, c)
		user := claims[auth.UsernameKey]
		if v, ok := user.(v1.SysUser); ok {
			// 设置登出状态
			tokenId := jwt.GetToken(ctx, c)
			redisLogoutStatus(tokenId)

			logininfo := auth.GetLogininfo(c, v.Username)
			logininfo.Status = api.LOGOUT
			logininfo.Msg = "登出成功"
			appendLogininfo(logininfo)
		}

		c.JSON(code2.SUCCESS, map[string]interface{}{
			"code": code,
		})
	}
}

func newAutoAuth(info *JwtInfo) auth.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth().(*auth.BasicStrategy), newJWTAuth(info).(*auth.JWTStrategy))
}

// 添加登陆信息
func appendLogininfo(logininfo *v12.Logininfo) {
	if strutil.ContainsAny(logininfo.Status, api.LOGIN_SUCCESS, api.LOGOUT, api.REGISTER) {
		logininfo.Status = api.LOGIN_SUCCESS_STATUS
	} else if logininfo.Status == api.LOGIN_FAIL {
		logininfo.Status = api.LOGIN_FAIL_STATUS
	}

	resp, err := rpc.Remoting.CreateSysLogininfo(context.Background(), &v12.CreateSysLogininfoRequest{
		LoginInfo: logininfo,
	})

	if err != nil || resp.Code != code2.SUCCESS {
		log.Errorf("login info append error: %s", err.Error())
	}
}

// 设置用户登录状态
func redisLoginStatus(tokenid string, token_ip_accesstime string, expireTime int64) {
	redisCli := kv.NewKVStore("redis").(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)
	redisCli.SetKey(context.Background(), tokenid, token_ip_accesstime, expireTime)
}

// 设置用户登出状态
func redisLogoutStatus(tokenId string) {
	redisCli := kv.NewKVStore("redis").(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)
	redisCli.DeleteKey(context.Background(), tokenId)
}

// 自定义参数校验
func validLoginErrMsg(err error) string {
	for _, e := range err.(validator.ValidationErrors) {
		switch e.ActualTag() {
		case "required":
			return fmt.Sprintf("%s 必须填写\n", fieldTranslations[e.StructNamespace()])
		case "min", "max":
			return fmt.Sprintf("%s 不在指定范围内\n", fieldTranslations[e.StructNamespace()])
		}
	}
	return err.Error()
}
