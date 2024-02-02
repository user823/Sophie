package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/pkg/utils"
	"golang.org/x/crypto/bcrypt"
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
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return utils.B2s(hashedBytes), err
}

func Compare(hashedPassword, plain string) error {
	return bcrypt.CompareHashAndPassword(utils.S2b(hashedPassword), utils.S2b(plain))
}
