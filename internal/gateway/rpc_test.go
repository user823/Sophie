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

func TestRPCLogininfor(t *testing.T) {

	resp, err := rpc.Remoting.CreateSysLogininfo(ctx, &v12.CreateSysLogininfoRequest{
		LoginInfo: &v12.Logininfo{UserName: "test"},
	})
	if err = rpc.ParseRpcErr(resp, err); err != nil {
		t.Error(err)
	}
	t.Logf("%s", resp.Msg)
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

	t.Run("test-RPCLogininfor", TestRPCLogininfor)
	t.Run("test-RPCRole", TesetRPCRole)
}
