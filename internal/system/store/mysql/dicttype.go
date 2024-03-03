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

type mysqlDictType struct {
	db *gorm.DB
}

var _ store.DictTypeStore = &mysqlDictType{}

func selectDictTypeVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_dict_type")
}

func (s *mysqlDictType) SelectDictTypeList(ctx context.Context, dictType *v1.SysDictType, opts *api.GetOptions) ([]*v1.SysDictType, error) {
	query := selectDictTypeVo(s.db)
	if dictType.DictName != "" {
		query = query.Where("dict_name like ?", "%"+dictType.DictName+"%")
	}
	if dictType.Status != "" {
		query = query.Where("status = ?", dictType.Status)
	}
	if dictType.DictType != "" {
		query = query.Where("dict_type like ?", "%"+dictType.DictType+"%")
	}
	query = opts.SQLCondition(query, "sys_dict_type.create_time")

	var result []*v1.SysDictType
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlDictType) SelectDictTypeAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysDictType, error) {
	return s.SelectDictTypeList(ctx, &v1.SysDictType{}, opts)
}

func (s *mysqlDictType) SelectDictTypeById(ctx context.Context, dictid int64, opts *api.GetOptions) (*v1.SysDictType, error) {
	query := selectDictTypeVo(s.db).Where("dict_id = ?", dictid)
	query = opts.SQLCondition(query, "")

	var result v1.SysDictType
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(dictid, ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlDictType) SelectDictTypeByType(ctx context.Context, dictType string, opts *api.GetOptions) (*v1.SysDictType, error) {
	query := selectDictTypeVo(s.db).Where("dict_type = ?", dictType)
	query = opts.SQLCondition(query, "")

	var result v1.SysDictType
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, dictType), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlDictType) DeleteDictTypeById(ctx context.Context, dictid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("dict_id = ?", dictid).Delete(&v1.SysDictType{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(dictid, ""))
}

func (s *mysqlDictType) DeleteDictTypeByIds(ctx context.Context, dictids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("dict_id in ?", dictids).Delete(&v1.SysDictType{})

	cacheKeys := make([]string, 0, len(dictids))
	for i := range dictids {
		cacheKeys = append(cacheKeys, s.CacheKey(dictids[i], ""))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return del.Error
}

func (s *mysqlDictType) InsertDictType(ctx context.Context, dictType *v1.SysDictType, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(dictType)
	return create.Error
}

func (s *mysqlDictType) UpdateDictType(ctx context.Context, dictType *v1.SysDictType, opts *api.UpdateOptions) error {
	if dictType.DictId == 0 {
		return fmt.Errorf("更新字典类型必须指明id")
	}

	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("dict_id = ?", dictType.DictId).Updates(dictType).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(dictType.DictId, ""))
}

func (s *mysqlDictType) CheckDictTypeUnique(ctx context.Context, dictType string, opts *api.GetOptions) *v1.SysDictType {
	query := selectDictTypeVo(s.db).Where("dict_type = ?", dictType)
	query = opts.SQLCondition(query, "")

	var result v1.SysDictType
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, dictType), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

var dictTypeCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlDictType) CachedDB() *cache.CachedDB {
	dictTypeCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis").(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-dicttypestore-")
		rdsCli.SetRandomExp(true)

		dictTypeCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, dictTypeCache.rdsCache)
}

func (s *mysqlDictType) CacheKey(dictTypeId int64, dictType string) string {
	// id:name
	return fmt.Sprintf("%d:%s", dictTypeId, dictType)
}
