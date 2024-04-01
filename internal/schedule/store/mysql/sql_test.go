package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/db/sql"
	"testing"
)

func InitSQL() {
	opts := options.NewMySQLOptions()
	mysqlConfig := &sql.MysqlConfig{
		Host:                  opts.Host,
		Username:              opts.Username,
		Password:              opts.Password,
		Database:              opts.Database,
		MaxIdleConnections:    opts.MaxIdleConnections,
		MaxOpenConnections:    opts.MaxOpenConnections,
		MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
		LogLevel:              opts.LogLevel,
		Logger:                nil,
	}
	GetMySQLFactoryOr(mysqlConfig)
}

func TestJobList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	list, _, err := sqlCli.Jobs().SelectJobList(context.Background(), &v1.SysJob{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range list {
		t.Logf("%v", list[i])
	}
}

func TestJobInsert(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Jobs().InsertJob(context.Background(), &v1.SysJob{
		JobId:   2,
		JobName: "test",
	}, &api.CreateOptions{})
	if err != nil {
		t.Logf("%v", err.Error())
	}
}

func TestSqlSub(t *testing.T) {
	InitSQL()

	t.Run("test-jobList", TestJobList)
	t.Run("test-jobCreate", TestJobInsert)
}
