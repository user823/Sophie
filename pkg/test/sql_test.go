package test

import (
	"fmt"
	"github.com/user823/Sophie/pkg/db/sql"
	"testing"
	"time"
)

func TestSQLConnection(t *testing.T) {
	db, err := sql.NewDB("mysql", &sql.MysqlConfig{
		Host:                  "127.0.0.1:3306",
		Username:              "sophie",
		Password:              "123456",
		Database:              "sophie",
		MaxIdleConnections:    10,
		MaxOpenConnections:    10,
		MaxConnectionLifeTime: 3600 * time.Second,
		LogLevel:              0,
		Logger:                nil,
	})
	if err != nil {
		t.Fatalf("Connection error: %s", err.Error())
	}

	var user User
	db.First(&user, &User{RoleId: 1})
	fmt.Println(user)
}

type User struct {
	RoleId   int    `gorm:"role_id"`
	RoleName string `gorm:"role_name"`
}

func (u *User) TableName() string {
	return "sys_role"
}
