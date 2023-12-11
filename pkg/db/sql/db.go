package sql

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrUnsupportedDB     = errors.New("DB not supported, only MySQL, PostgreSQL, SQLite, SQL Server or TiDB is valid")
	ErrConfigTypeInvalid = errors.New("cannot correctly convert the config type")
)

func NewDB(name string, config any) (*gorm.DB, error) {
	switch name {
	case "mysql":
		return NewMysqlDB(config)
	default:
		return nil, ErrUnsupportedDB
	}
}
