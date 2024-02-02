package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/system/v1"
	v12 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
)

type roles struct{}

func newRoles() *roles {
	return &roles{}
}

func (r *roles) List(ctx context.Context, sysUser *v1.SysUser, opts api.ListOptions) (*v1.RoleList, error) {
	resp, err := rpc.Remoting.GetSysRoleByUser(ctx, sysUser.UserId)
	if err = parseRpcErr(resp.BaseResp, err); err != nil {
		return nil, err
	}
	var res v1.RoleList
	for i := 0; i < int(resp.Total); i++ {
		res.Items = append(res.Items, *v12.RoleInfo2SysRole(resp.Rows[i]))
	}
	res.TotalCount = resp.Total
	return &res, nil
}
