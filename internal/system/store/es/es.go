package es

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/store/mysql"
	"github.com/user823/Sophie/pkg/db/doc"
	"sync"
)

// 作为mysql的缓存层
type datastore struct {
	es *elasticsearch.TypedClient
}

func (ds *datastore) Users() store.UserStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Users()
}

func (ds *datastore) UserPosts() store.UserPostStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.UserPosts()
}

func (ds *datastore) UserRoles() store.UserRoleStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.UserRoles()
}

func (ds *datastore) RoleMenus() store.RoleMenuStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.RoleMenus()
}

func (ds *datastore) Roles() store.RoleStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Roles()
}

func (ds *datastore) RoleDepts() store.RoleDeptStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.RoleDepts()
}

func (ds *datastore) Posts() store.PostStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Posts()
}

func (ds *datastore) OperLogs() store.OperLogStore {
	return &esOperLogStore{ds.es}
}

func (ds *datastore) Notices() store.NoticeStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Notices()
}

func (ds *datastore) Menus() store.MenuStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Menus()
}

func (ds *datastore) Logininfors() store.LogininforStore {
	return &esLogininforStore{ds.es}
}

func (ds *datastore) DictTypes() store.DictTypeStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.DictTypes()
}

func (ds *datastore) DictData() store.DictDataStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.DictData()
}

func (ds *datastore) Depts() store.DeptStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Depts()
}

func (ds *datastore) Configs() store.ConfigStore {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Configs()
}

func (ds *datastore) Close() error {
	return nil
}

func (ds *datastore) Begin() store.Factory {
	sqlCli, _ := mysql.GetMySQLFactoryOr(nil)
	return sqlCli.Begin()
}

func (ds *datastore) Commit() error {
	return nil
}

func (ds *datastore) Rollback() error {
	return nil
}

var (
	esFactory store.Factory
	once      sync.Once
)

// es作为数据库读缓存
// 对于写请求直接使用数据库
// 如果es客户端获取失败或者无法建立连接则只使用 数据库进行读写
func GetESFactoryOr(cfg *doc.ESConfig) (store.Factory, error) {
	if cfg == nil && esFactory == nil {
		return mysql.GetMySQLFactoryOr(nil)
	}

	var escli *elasticsearch.TypedClient
	var err error
	once.Do(func() {
		escli, err = doc.NewES(cfg)
		// 测试es连接
		if err != nil {
			return
		}

		ok, e := escli.Ping().Do(context.Background())
		if e != nil || !ok {
			err = fmt.Errorf("elasticsearch connection failed")
			return
		}
		esFactory = &datastore{escli}
	})

	return esFactory, err
}
