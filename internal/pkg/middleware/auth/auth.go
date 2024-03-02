package auth

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
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
