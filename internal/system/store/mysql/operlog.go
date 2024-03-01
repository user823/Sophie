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

type mysqlOperLogStore struct {
	db *gorm.DB
}

var _ store.OperLogStore = &mysqlOperLogStore{}

func selectOperLogVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_oper_log")
}

func (s *mysqlOperLogStore) InsertOperLog(ctx context.Context, operlog *v1.SysOperLog, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(operlog)
	return create.Error
}

func (s *mysqlOperLogStore) SelectOperLogList(ctx context.Context, operlog *v1.SysOperLog, opts *api.GetOptions) ([]*v1.SysOperLog, error) {
	query := selectOperLogVo(s.db)
	if operlog.OperIp != "" {
		query = query.Where("oper_ip like %?%", operlog.OperId)
	}
	if operlog.Title != "" {
		query = query.Where("title like %?%", operlog.Title)
	}
	if operlog.BusinessType != nil {
		query = query.Where("business_type = ?", operlog.BusinessType)
	}
	if len(operlog.BusinessTypes) > 0 {
		query = query.Where("business_type in ?", operlog.BusinessTypes)
	}
	if operlog.Status != "" {
		query = query.Where("status = ?", operlog.Status)
	}
	if operlog.OperName != "" {
		query = query.Where("oper_name like %?%", operlog.OperName)
	}
	query = opts.SQLCondition(query, "")
	query = query.Order("oper_id DESC")

	var result []*v1.SysOperLog
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlOperLogStore) DeleteOperLogByIds(ctx context.Context, operids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("oper_id in ?", operids).Delete(&v1.SysOperLog{})

	cacheKeys := make([]string, 0, len(operids))
	for i := range operids {
		cacheKeys = append(cacheKeys, s.CacheKey(operids[i]))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return del.Error
}

func (s *mysqlOperLogStore) SelectOperLogById(ctx context.Context, operid int64, opts *api.GetOptions) (*v1.SysOperLog, error) {
	query := selectOperLogVo(s.db).Where("oper_id = ?", operid)
	query = opts.SQLCondition(query, "")

	var result v1.SysOperLog
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(operid), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlOperLogStore) CleanOperLog(ctx context.Context, opts *api.DeleteOptions) error {
	del := s.db.Table("sys_oper_log").Where("true").Delete(nil)

	s.CachedDB().CleanCache(ctx)
	return del.Error
}

var operLogCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlOperLogStore) CachedDB() *cache.CachedDB {
	operLogCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis").(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-operlogstore-")
		rdsCli.SetRandomExp(true)

		operLogCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, operLogCache.rdsCache)
}

func (s *mysqlOperLogStore) CacheKey(operLogId int64) string {
	// id
	return fmt.Sprintf("%d", operLogId)
}
