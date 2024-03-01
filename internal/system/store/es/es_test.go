package es

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/gateway/v1"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/internal/system/store/mysql"
	"github.com/user823/Sophie/pkg/db/sql"
	"testing"
	"time"
)

var (
	ctx context.Context
)

func Init() {
	esOptions := &options.ESOptions{
		Addrs:    []string{"https://localhost:9200"},
		Username: "sophie",
		Password: "123456",
		MaxIdle:  10,
		UseSSL:   false,
		Timeout:  5 * time.Second,
	}
	_, err := GetESFactoryOr(esOptions)
	if err != nil {
		fmt.Printf("出错了 %s", err.Error())
		panic(err)
	}

	testLogininfo := v1.LoginUser{
		User: v12.SysUser{
			UserId:      2,
			DeptId:      105,
			Username:    "sophie",
			Nickname:    "sophie",
			Email:       "sophie@qq.com",
			Phonenumber: "15666666666",
			Sex:         "1",
			Password:    "$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2",
			DelFlag:     "0",
			Roles: []v12.SysRole{
				{RoleId: 2, RoleName: "普通角色", RoleKey: "common", RoleSort: 2, DataScope: "2", MenuCheckStrictly: true, DeptCheckStrictly: true},
			},
		},
		Roles:       []string{"common"},
		Permissions: []string{"*.*.*"},
	}
	ctx = context.WithValue(context.Background(), api.LOGIN_INFO_KEY, testLogininfo)

	cfg := &sql.MysqlConfig{
		Host:                  "127.0.0.1:3306",
		Username:              "sophie",
		Password:              "123456",
		Database:              "sophie",
		MaxIdleConnections:    10,
		MaxOpenConnections:    10,
		MaxConnectionLifeTime: 3600 * time.Second,
		LogLevel:              2,
	}
	_, err = mysql.GetMySQLFactoryOr(cfg)
	if err != nil {
		fmt.Printf("出错了 %s", err.Error())
		panic(err)
	}
}

func TestESSelectLogininforList(t *testing.T) {
	esCli, _ := GetESFactoryOr(nil)
	result, err := esCli.Logininfors().SelectLogininforList(ctx, &v12.SysLogininfor{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestESCleanLogininfor(t *testing.T) {
	esCli, _ := GetESFactoryOr(nil)
	err := esCli.Logininfors().CleanLogininfor(ctx, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestESSelectOperLogList(t *testing.T) {
	esCli, _ := GetESFactoryOr(nil)
	result, err := esCli.OperLogs().SelectOperLogList(ctx, &v12.SysOperLog{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestESSelectOperLogById(t *testing.T) {
	esCli, _ := GetESFactoryOr(nil)
	result, err := esCli.OperLogs().SelectOperLogById(ctx, 1, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestESCleanOperLog(t *testing.T) {
	esCli, _ := GetESFactoryOr(nil)
	err := esCli.OperLogs().CleanOperLog(ctx, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestESSub(t *testing.T) {
	Init()

	t.Run("test-SelectLogininforList", TestESSelectLogininforList)
	t.Run("test-CleanLogininfor", TestESCleanLogininfor)

	t.Run("test-SelectOperLogList", TestESSelectOperLogList)
	t.Run("test-SelectOperLogById", TestESSelectOperLogById)
	t.Run("test-TestESCleanOperLog", TestESCleanOperLog)
}
