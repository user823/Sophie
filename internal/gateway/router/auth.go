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
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/internal/pkg/middleware/auth"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"time"
)

var (
	fieldTranslations = map[string]string{
		"loginInfo.Username": "用户名",
		"loginInfo.Password": "用户密码",
	}
)

// 未授权时返回消息
var (
	ErrExpire       = errors.New("登录状态已过期，请重新登录")
	ErrUnauthorized = errors.New("账号或密码错误")
	ErrInternal     = errors.New("系统内部错误")
	ErrDeleted      = errors.New("对不起，您的账号已删除")
	ErrDisabled     = errors.New("对不起，您的账号已停用")
)

// 登陆参数绑定与校验
type loginInfo struct {
	Username string `json:"username" form:"username" vd:"required,min=2,max=20"`
	Password string `json:"password" form:"password" vd:"required,min=5,max=20"`
}

func newBasicAuth() auth.AuthStrategy {
	return auth.NewBasicStrategy(func(c *app.RequestContext, username, password string) (data interface{}, ok bool) {
		var msg string
		sysLoginInfo := auth.GetLogininfo(c, &v12.UserInfo{UserName: username}, "")
		defer func() {
			if msg != "" {
				sysLoginInfo.Status = api.LOGIN_FAIL
				sysLoginInfo.Msg = msg
			}
			appendLogininfo(sysLoginInfo)
		}()

		resp, err := rpc.Remoting.GetUserInfoByName(context.Background(), username)
		if err != nil || resp.BaseResp.Code != code.SUCCESS {
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
		Realm:            info.Realm,
		Key:              utils.S2b(info.Key),
		Timeout:          info.Timeout,
		MaxRefresh:       info.Timeout,
		IdentityKey:      auth.UsernameKey,
		PayloadFunc:      payloadFunc(),
		Authenticator:    authenticator(),
		Authorizator:     authentizator(),
		Unauthorized:     unauthorized(),
		SigningAlgorithm: "HS256",
		LoginResponse:    loginResponse(),
		LogoutResponse:   logoutResponse(),
		TokenLookup:      "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:    "Bearer",
		SendCookie:       false,
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

		// 参数验证/校验失败
		if err = c.BindAndValidate(&login); err != nil {
			return "", errors.New(validLoginErrMsg(err))
		}

		// 记录基本登陆信息
		sysLoginInfo := auth.GetLogininfo(c, &v12.UserInfo{UserName: login.Username}, "")
		defer func() {
			if err != nil {
				sysLoginInfo.Status = api.LOGIN_FAIL
				sysLoginInfo.Msg = err.Error()
			}
			appendLogininfo(sysLoginInfo)
		}()

		resp, err := rpc.Remoting.GetUserInfoByName(context.Background(), login.Username)
		if err != nil || resp.BaseResp.Code != code.SUCCESS {
			return nil, ErrInternal
		}

		user := resp.Data

		if v1.UserStatus["DELETED"].Code == user.DelFlag {
			return "", ErrDeleted
		}

		if v1.UserStatus["DISABLE"].Code == user.Status {
			return "", ErrDisabled
		}

		// 验证密码
		if err = auth.Compare(user.Password, login.Password); err != nil {
			return "", ErrUnauthorized
		}

		// 登陆成功，设置授权上下文信息
		loginUser := v12.LoginUser{
			Roles:       resp.GetRoles(),
			Permissions: resp.GetPermissions(),
			User:        resp.Data,
		}
		c.Set(api.LOGIN_INFO_KEY, loginUser)
		c.Set(auth.TokenIdKey, sysLoginInfo.TokenId)
		return map[string]any{
			"user":    user.UserName,
			"tokenId": sysLoginInfo.TokenId,
		}, nil
	}
}

// 登陆成功时设置jwt 载荷信息
func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(map[string]any); ok {
			return jwt.MapClaims{
				auth.UsernameKey: v["user"],
				// 每个token 绑定一个tokenId
				auth.TokenIdKey: v["tokenId"],
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

		// 检查用户登录状态
		claims := jwt.ExtractClaims(ctx, c)
		if tokenId, ok := claims[auth.TokenIdKey].(string); ok && checkLogin(tokenId) {
			if username, ok := data.(string); ok {

				// 用于后续步骤的权限校验
				resp, err := rpc.Remoting.GetUserInfoByName(ctx, username)
				if err != nil || resp.BaseResp.Code != code.SUCCESS {
					log.Errorf("Get user info error, code: %d, error: %v", resp.BaseResp.Code, err)
					return false
				}

				loginUser := v12.LoginUser{
					Roles:       resp.GetRoles(),
					Permissions: resp.GetPermissions(),
					User:        resp.Data,
				}
				c.Set(api.LOGIN_INFO_KEY, loginUser)
				return true
			}
		}

		return false
	}
}

// message 是根据authenticator的err构造的信息
func unauthorized() func(ctx context.Context, c *app.RequestContext, code int, message string) {
	return func(ctx context.Context, c *app.RequestContext, code int, message string) {
		if strutil.ContainsAny(message, ErrUnauthorized.Error(), ErrDisabled.Error(), ErrDeleted.Error(), ErrInternal.Error()) {
			core.WriteResponse(c, core.ErrResponse{
				Code:    500,
				Message: message,
			})
			return
		}

		core.WriteResponse(c, core.ErrResponse{
			Code:    401,
			Message: ErrExpire.Error(),
		})
	}
}

// 设置token 和 过期时间
func loginResponse() func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
	return func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
		data := map[string]interface{}{
			"access_token": message,
			"expires_in":   time,
		}

		if token, ok := c.Get(auth.TokenIdKey); ok {
			tokenId := token.(string)
			carry, _ := c.Get(api.LOGIN_INFO_KEY)
			v := carry.(v12.LoginUser)
			sysLoginInfo := auth.GetLogininfo(c, v.User, tokenId)
			redisLoginStatus(sysLoginInfo)
		} else {
			log.Warn(errors.New("Set user login status failed"))
		}

		core.OK(c, "", data)
	}
}

func logoutResponse() func(ctx context.Context, c *app.RequestContext, code int) {
	return func(ctx context.Context, c *app.RequestContext, cod int) {
		claims := jwt.ExtractClaims(ctx, c)
		user := claims[auth.UsernameKey]
		tokenid := claims[auth.TokenIdKey]
		if v, ok := user.(v1.SysUser); ok {
			// 设置登录日志
			logininfo := auth.GetLogininfo(c, &v12.UserInfo{UserName: v.Username}, "xxx")
			logininfo.Status = api.LOGOUT
			logininfo.Msg = "登出成功"
			appendLogininfo(logininfo)
		}

		if v, ok := tokenid.(string); ok {
			// 设置登出状态
			log.Infof("设置登出状态， 登出成功")
			redisLogoutStatus(v)
		}

		c.JSON(code.SUCCESS, map[string]interface{}{
			"code": cod,
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

	if err != nil || resp.Code != code.SUCCESS {
		log.Errorf("login info append error: %s", err.Error())
	}
}

// 检查用户登录状态
func checkLogin(tokenId string) bool {
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)
	ok, err := redisCli.Exists(context.Background(), tokenId)
	if err != nil {
		return false
	}
	return ok
}

// 设置用户登录状态
func redisLoginStatus(logininfo *v12.Logininfo) {
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)
	expireTime := api.LOGIN_TIMEOUT * time.Second
	// 设置登录状态
	_, result := redisCli.SetRollingWindow(context.Background(), kv.SYS_LOGIN_USER_IDS, expireTime.Milliseconds(), logininfo.TokenId, true)
	// 清理过期token
	for i := range result {
		redisCli.DeleteKey(context.Background(), result[i])
	}
	// 设置登录信息
	redisCli.AddToHash(context.Background(), logininfo.TokenId, logininfo.Marshal())
}

// 设置用户登出状态
func redisLogoutStatus(tokenId string) {
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
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
