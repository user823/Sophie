package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
)

type menus struct{}

func newMenus() *menus {
	return &menus{}
}

func (m *menus) GetPerms(ctx context.Context, roleIds []int64, opts api.ListOptions) ([]string, error) {
	resp, err := rpc.Remoting.GetSysMenuPermsByRoleIds(ctx, &v1.GetSysMenuPermsByRoleIdsRequest{
		RoleIds: roleIds,
	})
	if err = parseRpcErr(resp.BaseResp, err); err != nil {
		return []string{}, err
	}
	return resp.Perms, nil
}
