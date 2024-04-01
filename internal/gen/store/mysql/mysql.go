package mysql

import (
	"fmt"
	"github.com/user823/Sophie/internal/gen/store"
	"github.com/user823/Sophie/pkg/db/sql"
	"github.com/user823/Sophie/pkg/errors"
	"gorm.io/gorm"
	"sync"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) GenTables() store.GenTableStore {
	return &mysqlGenTableStore{ds.db}
}

func (ds *datastore) GenTableColumns() store.GenTableColumnStore {
	return &mysqlGenTableColumnStore{ds.db}
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}
	return db.Close()
}

func (ds *datastore) Begin() store.Factory {
	return &datastore{ds.db.Begin()}
}

// Commit 后不能再使用该对象创建的xxxStore对象
func (ds *datastore) Commit() error {
	return ds.db.Commit().Error
}

// Rollback 后不能再使用该对象创建的xxxStore对象
func (ds *datastore) Rollback() error {
	return ds.db.Rollback().Error
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

func GetMySQLFactoryOr(cfg *sql.MysqlConfig) (store.Factory, error) {
	if cfg == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}

	// 已经初始化过
	if mysqlFactory != nil {
		return mysqlFactory, nil
	}

	// 尝试初始化
	if cfg != nil {
		dbIns, err := sql.NewMysqlDB(cfg)
		if err != nil {
			return nil, err
		}
		once.Do(func() {
			mysqlFactory = &datastore{dbIns}
		})
	}

	return mysqlFactory, nil
}
