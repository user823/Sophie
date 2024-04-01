package auth

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
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
	TokenIdKey  = "tokenId"
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
func GetLogininfo(c *app.RequestContext, userinfo *v12.UserInfo, tokenId string) *v12.Logininfo {
	accessTime := time.Now()
	uaStr := utils.B2s(c.UserAgent())

	res := &v12.Logininfo{
		UserName:   userinfo.UserName,
		Status:     api.LOGIN_SUCCESS,
		Ipaddr:     utils.GetClientIP(c),
		Msg:        "登录成功",
		AccessTime: utils.Time2Str(accessTime),
		TokenId:    tokenId,
		LoginTime:  utils.Time2Str(accessTime),
	}

	if userinfo.Dept != nil {
		res.DeptName = userinfo.Dept.DeptName
	}

	if uaStr != "" {
		ua := useragent.Parse(uaStr)
		res.Browser = ua.Name + ":" + ua.Version
		res.Os = ua.OS
	}

	if tokenId == "" {
		res.TokenId = uuid.NewString()
	}
	return res
}
