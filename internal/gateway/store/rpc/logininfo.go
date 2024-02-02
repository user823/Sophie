package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v12 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
)

type logininfos struct{}

func newLogininfos() *logininfos {
	return &logininfos{}
}

func (l *logininfos) Create(ctx context.Context, logininfo *v12.Logininfo, opts api.CreateOptions) error {
	resp, err := rpc.Remoting.CreateSysLogininfo(ctx, &v12.CreateSysLogininfoRequest{
		LoginInfo: logininfo,
		Source:    api.INNER,
	})
	if err = parseRpcErr(resp, err); err != nil {
		return err
	}
	return nil
}
