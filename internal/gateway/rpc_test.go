package gateway

import (
	"context"
	v12 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"testing"
)

var (
	ctx = context.Background()
)

func RPCInit() {
	opts := NewOptions()
	opts.SecureServing.Required = false
	cfg, _ := CreateConfigFromOptions(opts)
	rpcGeneralOpts := cfg.CreateRemoteClientOptions()
	rpc.Init(rpcGeneralOpts)
}

func TestRPCGetLogininfo(t *testing.T) {
	resp, err := rpc.Remoting.GetUserInfoByName(ctx, "sophie")
	if err != nil {
		t.Error(err)
	}
	user := v12.UserInfo2SysUser(resp.Data)
	t.Logf("%v", user)
	t.Logf("%v", resp.Roles)
	t.Logf("%v", resp.Permissions)
}

func TesetRPCRole(t *testing.T) {
	resp, err := rpc.Remoting.ListSysRole(ctx, &v12.ListSysRolesRequest{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", resp.Rows)

}

func TestSub(t *testing.T) {
	RPCInit()

	t.Run("test-RPCLogininfor", TestRPCGetLogininfo)
	t.Run("test-RPCRole", TesetRPCRole)
}
