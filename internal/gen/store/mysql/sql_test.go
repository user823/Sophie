package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/db/sql"
	"testing"
)

var (
	ctx context.Context
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
		Debug:                 true,
	}
	ctx = context.Background()
	GetMySQLFactoryOr(mysqlConfig)
}

func TestSelectGenTableList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, total, err := sqlCli.GenTables().SelectGenTableList(ctx, &v1.GenTable{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("总记录数: %d", total)
	for i := range result {
		t.Logf("%v", result[i])
		for j := range result[i].Columns {
			t.Logf("No.%d Columnt: %v", i, result[i].Columns[j])
		}
	}
}

func TestSelectDbTableList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, total, err := sqlCli.GenTables().SelectDbTableList(ctx, &v1.GenTable{Tablename: "sys_user"}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("总记录数: %d", total)
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectDbTableListByNames(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, total, err := sqlCli.GenTables().SelectDbTableListByNames(ctx, []string{"sys_user_role", "sys_notice"}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("总记录数: %d", total)
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectGenTableAll(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, total, err := sqlCli.GenTables().SelectGenTableAll(ctx, &api.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("总记录数: %d", total)
	for i := range result {
		t.Logf("%v", result[i])
		for j := range result[i].Columns {
			t.Logf("No.%d Columnt: %v", i, result[i].Columns[j])
		}
	}
}

func TestSelectDbTableColumnsByName(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, total, err := sqlCli.GenTableColumns().SelectDbTableColumnsByName(ctx, "sys_user_role", &api.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("总记录数: %d", total)
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestColumnSelectDbTableColumnsByName(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, total, err := sqlCli.GenTableColumns().SelectDbTableColumnsByName(ctx, "sys_menu", &api.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("总记录数: %d", total)
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectGenTableColumnListByTableId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, total, err := sqlCli.GenTableColumns().SelectGenTableColumnListByTableId(ctx, 12, &api.GetOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("总记录数: %d", total)
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestInsertGenTable(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	sqlCli.GenTables().InsertGenTable(ctx, &v1.GenTable{Tablename: "test1"}, &api.CreateOptions{})
	sqlCli.GenTables().InsertGenTable(ctx, &v1.GenTable{Tablename: "test2"}, &api.CreateOptions{})
}

func TestSub(t *testing.T) {
	InitSQL()

	// table
	t.Run("test-SelectGenTableList", TestSelectGenTableList)
	t.Run("test-SelectDbTableList", TestSelectDbTableList)
	t.Run("test-SelectDbTableListByNames", TestSelectDbTableListByNames)
	t.Run("test-SelectGenTableAll", TestSelectGenTableAll)
	t.Run("test-SelectDbTableColumnsByName", TestSelectDbTableColumnsByName)
	t.Run("test-InsertGenTable", TestInsertGenTable)

	// column
	t.Run("test-ColumnSelectDbTableColumnsByName", TestColumnSelectDbTableColumnsByName)
	t.Run("test-SelectGenTableColumnListByTableId", TestSelectGenTableColumnListByTableId)
}
