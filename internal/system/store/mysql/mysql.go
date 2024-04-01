package mysql

import (
	"fmt"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/pkg/db/sql"
	"github.com/user823/Sophie/pkg/errors"
	"gorm.io/gorm"
	"sync"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Users() store.UserStore {
	return &mysqlUserStore{ds.db}
}

func (ds *datastore) UserPosts() store.UserPostStore {
	return &mysqlUserPostStore{ds.db}
}

func (ds *datastore) UserRoles() store.UserRoleStore {
	return &mysqlUserRoleStore{ds.db}
}

func (ds *datastore) RoleMenus() store.RoleMenuStore {
	return &mysqlRoleMenuStore{ds.db}
}

func (ds *datastore) Roles() store.RoleStore {
	return &mysqlRoleStore{ds.db}
}

func (ds *datastore) RoleDepts() store.RoleDeptStore {
	return &mysqlRoleDeptStore{ds.db}
}

func (ds *datastore) Posts() store.PostStore {
	return &mysqlPostStore{ds.db}
}

func (ds *datastore) OperLogs() store.OperLogStore {
	return &mysqlOperLogStore{ds.db}
}

func (ds *datastore) Notices() store.NoticeStore {
	return &mysqlNoticeStore{ds.db}
}

func (ds *datastore) Menus() store.MenuStore {
	return &mysqlMenuStore{ds.db}
}

func (ds *datastore) Logininfors() store.LogininforStore {
	return &mysqlLogininforStore{ds.db}
}

func (ds *datastore) DictTypes() store.DictTypeStore {
	return &mysqlDictType{ds.db}
}

func (ds *datastore) DictData() store.DictDataStore {
	return &mysqlDictDataStore{ds.db}
}

func (ds *datastore) Depts() store.DeptStore {
	return &mysqlDeptStore{ds.db}
}

func (ds *datastore) Configs() store.ConfigStore {
	return &mysqlConfigStore{ds.db}
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
