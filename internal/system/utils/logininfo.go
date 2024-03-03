package utils

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
)

var (
	ErrPermsNotGet = fmt.Errorf("获取用户权限信息异常")
)

func GetLogininfoFromCtx(ctx context.Context) *v1.LoginUser {
	data := ctx.Value(api.LOGIN_INFO_KEY)
	return data.(*v1.LoginUser)
}
