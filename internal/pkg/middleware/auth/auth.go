package auth

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	v12 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	// 令牌自定义标识
	AUTHENTICATION = "Authorization"
	// 令牌前缀
	BASICPREFIX = "Basic"
	JWTPREFIX   = "Bearer"
	UsernameKey = "username"
)

// 授权策略
type AuthStrategy interface {
	AuthFunc() app.HandlerFunc
}

// 策略模式支持多种不同的授权策略
type AuthOperator struct {
	strategy AuthStrategy
}

func (a *AuthOperator) SetStrategy(strategy AuthStrategy) {
	a.strategy = strategy
}

func (a *AuthOperator) AuthFunc() app.HandlerFunc {
	return a.strategy.AuthFunc()
}

func Encrypt(plain string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(utils.S2b(plain), bcrypt.DefaultCost)
	return utils.B2s(hashedBytes), err
}

func Compare(hashedPassword, plain string) error {
	return bcrypt.CompareHashAndPassword(utils.S2b(hashedPassword), utils.S2b(plain))
}

// 获取基本登陆信息
func GetLogininfo(c *app.RequestContext, username string) *v12.Logininfo {
	ip := utils.GetClientIP(c)
	accessTime := time.Now().Unix()
	status := api.LOGIN_SUCCESS
	msg := "登陆成功"
	return &v12.Logininfo{
		UserName:   username,
		InfoId:     0,
		Ipaddr:     ip,
		AccessTime: accessTime,
		Status:     status,
		Msg:        msg,
	}
}

// 授权成功后设置上下文信息
func Authorizator(data interface{}, c *app.RequestContext) bool {

	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(v1.SysUser)
		r, _ := v["role"].(v1.SysRole)
		c.Set("role", r.RoleName)
		c.Set("roleIds", r.RoleId)
		c.Set("userId", u.UserId)
		c.Set("userName", u.Username)
		c.Set("dataScope", r.DataScope)
		return true
	}
	return false
}
