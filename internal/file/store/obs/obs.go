package obs

import (
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/goss"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/internal/file/store"
	"github.com/user823/Sophie/pkg/log"
	"sync"
)

type datastore struct {
	storage core.Storage
}

var _ store.Factory = &datastore{}

func (ds *datastore) Files() store.FileStore {
	return &obsFileStore{ds.storage}
}

var (
	storage core.Storage
	once    sync.Once
)

func GetOBSFactoryOr() store.Factory {
	once.Do(func() {
		// 如果没有加载，则手动加载一次
		if !viper.IsSet("driver") {
			viper.SetConfigFile("../../../configs/goss.yml")
			if err := viper.ReadInConfig(); err != nil {
				log.Errorf("viper 读取obs配置失败: %s", err.Error())
				return
			}
		}
		gs, err := goss.NewWithViper(viper.GetViper())
		if err != nil {
			log.Fatal("OBS 服务加载失败，请检查配置: %s", err.Error())
			return
		}
		storage = gs.Storage
	})

	return &datastore{storage: storage}
}
