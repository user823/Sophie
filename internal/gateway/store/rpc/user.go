package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v12 "github.com/user823/Sophie/api/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
)

type users struct{}

func newUsers() *users {
	return &users{}
}

// 用户注册
func (u *users) Create(ctx context.Context, userinfo *v1.UserInfo, opts api.CreateOptions) error {
	resp, err := rpc.Remoting.RegisterSysUser(ctx, &v1.RegisterSysUserRequest{userinfo})
	return parseRpcErr(resp.BaseResp, err)
}

func (u *users) Get(ctx context.Context, username string, opts api.GetOptions) (*v12.SysUser, error) {
	var err error
	resp, err := rpc.Remoting.GetUserInfoByName(ctx, username)
	if err = parseRpcErr(resp.BaseResp, err); err != nil {
		return nil, err
	}
	return v1.UserInfo2SysUser(resp.Data), nil
}
