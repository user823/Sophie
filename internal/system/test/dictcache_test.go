package test

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/internal/system/store/mysql"
	"github.com/user823/Sophie/internal/system/utils/cacheutils"
	"github.com/user823/Sophie/pkg/db/kv/redis"
	"testing"
	"time"
)

func Init() {
	mysqlOptions := &options.MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "sophie",
		Password:              "123456",
		Database:              "sophie",
		MaxIdleConnections:    10,
		MaxOpenConnections:    10,
		MaxConnectionLifeTime: 3600 * time.Second,
		LogLevel:              2,
	}
	mysql.GetMySQLFactoryOr(mysqlOptions)

	connectionConfig := &redis.RedisConfig{
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "123456",
		Database: 0,
	}
	go redis.KeepConnection(context.Background(), connectionConfig)
	time.Sleep(2 * time.Second)
	if !redis.Connected() {
		fmt.Printf("redis 未连接成功")
	}

	cacheutils.LoadingDictCache(nil)
}

func TestCleanDictCache(t *testing.T) {
	cacheutils.CleanDictCache()
}

func TestGetDictCache(t *testing.T) {
	result := cacheutils.GetDictCache("sys_user_sex")
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestDictCacheSub(t *testing.T) {
	Init()

	t.Run("test-CleanDictCache", TestCleanDictCache)
	t.Run("test-GetDictCache", TestGetDictCache)
}
