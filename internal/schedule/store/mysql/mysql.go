package mysql

import (
	"fmt"
	store2 "github.com/user823/Sophie/internal/schedule/store"
	"github.com/user823/Sophie/pkg/db/sql"
	"gorm.io/gorm"
	"sync"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Jobs() store2.JobStore {
	return &mysqlJobStore{ds.db}
}

func (ds *datastore) JobLogs() store2.JobLogStore {
	return &mysqlJobLogStore{ds.db}
}

var (
	mysqlFactory store2.Factory
	once         sync.Once
)

func GetMySQLFactoryOr(cfg *sql.MysqlConfig) (store2.Factory, error) {
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
