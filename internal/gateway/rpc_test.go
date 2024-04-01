package gateway

import (
	"context"
	"fmt"
	v12 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/pkg/db/kv"
	"testing"
	"time"
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

	redisConfig := &kv.RedisConfig{
		Addrs:                 opts.RedisOptions.Addrs,
		MasterName:            opts.RedisOptions.MasterName,
		Username:              opts.RedisOptions.Username,
		Password:              opts.RedisOptions.Password,
		Database:              opts.RedisOptions.Database,
		MaxIdle:               opts.RedisOptions.MaxIdle,
		MaxActive:             opts.RedisOptions.MaxActive,
		Timeout:               opts.RedisOptions.Timeout,
		EnableCluster:         opts.RedisOptions.EnableCluster,
		UseSSL:                opts.RedisOptions.UseSSL,
		SSLInsecureSkipVerify: opts.RedisOptions.SSLInsecureSkipVerify,
	}

	go kv.KeepConnection(ctx, redisConfig)
	time.Sleep(2 * time.Second)
	if !kv.Connected() {
		fmt.Printf("redis 未连接成功")
	}
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

func TestRPCUser(t *testing.T) {
	resp, err := rpc.Remoting.GetUserInfoById(ctx, &v12.GetUserInfoByIdRequest{
		Id: 1,
		User: &v12.LoginUser{
			User: &v12.UserInfo{UserId: 1},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", resp.Data)
}

func TestRPCDept(t *testing.T) {
	resp, err := rpc.Remoting.GetDeptById(ctx, &v12.GetDeptByIdReq{
		Id: 103,
		User: &v12.LoginUser{
			User: &v12.UserInfo{UserId: 1},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", resp.Data)
}

func TestGetRouters(t *testing.T) {
	resp, err := rpc.Remoting.GetRouters(ctx, &v12.GetRoutersRequest{
		User: &v12.LoginUser{
			User: &v12.UserInfo{
				UserId:   1,
				UserName: "admin",
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", resp.Data)
}

func TestGetLoginUser(t *testing.T) {
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)
	keys, _ := redisCli.GetListRange(ctx, kv.SYS_LOGIN_USER_IDS, 0, -1)
	list, _ := redisCli.MGetFromHash(ctx, keys)
	var result v12.Logininfo
	result.Unmarshal(list[0])
	t.Log(result)
}

func TestForceLogout(t *testing.T) {
	tokenId := "3e06d4e6-8a93-4ec5-84eb-b863ddfc0f99"
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)
	redisCli.DeleteKey(context.Background(), tokenId)
	redisCli.RemoveFromList(context.Background(), kv.SYS_LOGIN_USER_IDS, tokenId)
}

func TestSub(t *testing.T) {
	RPCInit()

	t.Run("test-RPCLogininfor", TestRPCGetLogininfo)
	t.Run("test-RPCRole", TesetRPCRole)
	t.Run("test-RPCUser", TestRPCUser)
	t.Run("test-RPCDept", TestRPCDept)
	t.Run("test-GetRouters", TestGetRouters)
	t.Run("test-SetLoginUser", TestGetLoginUser)
	t.Run("test-ForceLogout", TestForceLogout)
}
