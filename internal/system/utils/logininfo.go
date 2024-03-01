package utils

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/gateway/v1"
)

var (
	ErrPermsNotGet = fmt.Errorf("获取用户权限信息异常")
)

func GetLogininfoFromCtx(ctx context.Context) (v1.LoginUser, error) {
	data := ctx.Value(api.LOGIN_INFO_KEY)
	if data == nil {
		return v1.LoginUser{}, ErrPermsNotGet
	}
	logininfo, ok := data.(v1.LoginUser)
	if !ok {
		return v1.LoginUser{}, ErrPermsNotGet
	}
	return logininfo, nil
}
