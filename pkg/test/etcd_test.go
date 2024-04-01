package test

import (
	"context"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/utils"
	"testing"
)

var (
	etcdCli kv.KeyValueStore
)

func InitEtcd() {
	cfg := &kv.EtcdConfig{
		Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		//Username:  "sophie",
		//Password:  "123456",
		UseSSL:  false,
		Timeout: 3,
	}
	etcdCli = kv.NewKVStore("etcd", cfg)
}

func TestEtcdGetKeys(t *testing.T) {
	allKeys := etcdCli.GetKeys(context.TODO(), "")
	for i := range allKeys {
		t.Logf("%v", allKeys[i])
	}
}

func TestEtcdGet(t *testing.T) {
	value, err := etcdCli.GetKey(context.TODO(), "/config/test")
	if err != nil {
		t.Log(err)
	}
	t.Logf(value)
}

func TestEtcdPut(t *testing.T) {
	err := etcdCli.SetKey(context.TODO(), "/config/test", "2", utils.SecondToNano(100))
	if err != nil {
		t.Log(err)
	}
}

func TestEtcdGetKeysAndValuesWithFilter(t *testing.T) {
	res := etcdCli.GetKeysAndValuesWithFilter(context.TODO(), "sophie-schedule-nodes")
	for k, v := range res {
		t.Logf("k: %s v:%s", k, v)
	}
}

func TestEtcdSub(t *testing.T) {
	InitEtcd()

	t.Run("test-GetAllKey", TestEtcdGetKeys)
	t.Run("test-Getkey", TestEtcdGet)
	t.Run("test-Put", TestEtcdPut)
	t.Run("test-GetKeysAndValuesWithFilter", TestEtcdGetKeysAndValuesWithFilter)
}
