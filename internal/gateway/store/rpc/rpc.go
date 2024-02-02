package mysql

import (
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/store"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/errors"
)

type datastore struct {
}

func (ds *datastore) Users() store.UserStore {
	return newUsers()
}

func (ds *datastore) Logininfo() store.LogininfoStore {
	return newLogininfos()
}

func (db *datastore) Roles() store.RoleStore {
	return newRoles()
}

func (db *datastore) Menus() store.MenuStore {
	return newMenus()
}

func GetRPCFactory() store.Factory {
	return &datastore{}
}

//func (ds *datastore) Close() error {
//	db, err := ds.db.DB()
//	if err != nil {
//		return errors.WithMessage(err, "get gorm db instance failed")
//	}
//	return db.Close()
//}

//var (
//	mysqkFactory store.Factory
//	once         sync.Once
//)

// 懒加载
//func GetMySQLFactoryOr(cfg *sql.MysqlConfig) (store.Factory, error) {
//	if cfg == nil && mysqkFactory == nil {
//		return nil, fmt.Errorf("failed to get mysql store fatory")
//	}
//
//	var err error
//	var dbIns *gorm.DB
//	once.Do(func() {
//		dbIns, err = sql.NewMysqlDB(cfg)
//		mysqkFactory = &datastore{dbIns}
//	})
//
//	if mysqkFactory == nil || err != nil {
//		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqkFactory, err)
//	}
//	return mysqkFactory, nil
//}

func parseRpcErr(b *v1.BaseResp, err error) error {
	if err != nil {
		return err
	}

	if b.Code != code.SUCCESS {
		return errors.WithMessage(err, b.Msg)
	}
	return nil
}
