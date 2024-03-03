package mysql

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/pkg/cache"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/pkg/db/kv"
	"gorm.io/gorm"
	"sync"
)

type mysqlConfigStore struct {
	db *gorm.DB
}

var _ store.ConfigStore = &mysqlConfigStore{}

func selectConfigVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_config")
}

func (s *mysqlConfigStore) SelectConfig(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) (*v1.SysConfig, error) {
	query := selectConfigVo(s.db)
	query = opts.SQLCondition(query, "")
	if config.ConfigId != 0 {
		query = query.Where("config_id = ?", config.ConfigId)
	}
	if config.ConfigKey != "" {
		query = query.Where("config_key = ?", config.ConfigKey)
	}

	var result v1.SysConfig
	err := query.First(&result).Error
	return &result, err
}

func (s *mysqlConfigStore) SelectConfigById(ctx context.Context, configid int64, opts *api.GetOptions) (*v1.SysConfig, error) {
	query := selectConfigVo(s.db).Where("config_id = ?", configid)
	query = opts.SQLCondition(query, "")

	var result v1.SysConfig
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(configid, ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlConfigStore) SelectConfigList(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) ([]*v1.SysConfig, error) {
	query := selectConfigVo(s.db)
	if config.ConfigName != "" {
		query = query.Where("config_name like ?", "%"+config.ConfigName+"%")
	}
	if config.ConfigType != "" {
		query = query.Where("config_type = ?", config.ConfigType)
	}
	if config.ConfigKey != "" {
		query = query.Where("config_key like ?", "%"+config.ConfigKey+"%")
	}
	query = opts.SQLCondition(query, "sys_config.create_time")

	var result []*v1.SysConfig
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlConfigStore) CheckConfigKeyUnique(ctx context.Context, configKey string, opts *api.GetOptions) *v1.SysConfig {
	query := selectConfigVo(s.db).Where("config_key = ?", configKey)
	query = opts.SQLCondition(query, "")

	var result v1.SysConfig
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, configKey), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

func (s *mysqlConfigStore) InsertConfig(ctx context.Context, config *v1.SysConfig, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(config)
	return create.Error
}

func (s *mysqlConfigStore) UpdateConfig(ctx context.Context, config *v1.SysConfig, opts *api.UpdateOptions) error {
	if config.ConfigId == 0 {
		return fmt.Errorf("更新系统配置必须指明id")
	}
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(config).Where("config_id = ?", config.ConfigId).Updates(config).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(config.ConfigId, ""))
}

func (s *mysqlConfigStore) DeleteConfigById(ctx context.Context, configid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("config_id = ?", configid).Delete(&v1.SysConfig{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(configid, ""))
}

func (s *mysqlConfigStore) DeleteConfigByIds(ctx context.Context, configids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("config_id in ?", configids).Delete(&v1.SysConfig{})

	cacheKeys := make([]string, 0, len(configids))
	for i := range configids {
		cacheKeys = append(cacheKeys, s.CacheKey(configids[i], ""))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return del.Error
}

var configCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlConfigStore) CachedDB() *cache.CachedDB {
	configCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis").(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-configstore-")
		rdsCli.SetRandomExp(true)

		configCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, configCache.rdsCache)
}

func (s *mysqlConfigStore) CacheKey(configId int64, key string) string {
	// id:name
	return fmt.Sprintf("%d:%s", configId, key)
}
