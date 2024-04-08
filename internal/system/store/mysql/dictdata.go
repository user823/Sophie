package mysql

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/pkg/cache"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/db/kv"
	"gorm.io/gorm"
	"sync"
)

type mysqlDictDataStore struct {
	db *gorm.DB
}

var _ store.DictDataStore = &mysqlDictDataStore{}

func selectDictDataVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_dict_data")
}

func (s *mysqlDictDataStore) SelectDictDataList(ctx context.Context, dictData *v1.SysDictData, opts *api.GetOptions) ([]*v1.SysDictData, int64, error) {
	query := selectDictDataVo(s.db)
	if dictData.DictType != "" {
		query = query.Where("dict_type = ?", dictData.DictType)
	}
	if dictData.DictLabel != "" {
		query = query.Where("dict_label = ?", dictData.DictLabel)
	}
	if dictData.Status != "" {
		query = query.Where("status = ?", dictData.Status)
	}
	total := utils.CountQuery(query, opts, "")
	query = opts.SQLCondition(query, "")
	query = query.Order("dict_sort ASC")

	var result []*v1.SysDictData
	err := query.Find(&result).Error
	return result, total, err
}

func (s *mysqlDictDataStore) SelectDictDataByType(ctx context.Context, dictType string, opts *api.GetOptions) ([]*v1.SysDictData, int64, error) {
	query := selectDictDataVo(s.db).Where("status = 0 and dict_type = ?", dictType)
	query = opts.SQLCondition(query, "")
	query = query.Order("dict_sort ASC")

	var result []*v1.SysDictData
	err := query.Find(&result).Error
	return result, utils.CountQuery(query, opts, ""), err
}

func (s *mysqlDictDataStore) SelectDictLabel(ctx context.Context, dictType, dictValue string, opts *api.GetOptions) (string, error) {
	query := selectDictDataVo(s.db).Where("dict_type = ? and dict_value = ?", dictType, dictValue)
	var result v1.SysDictData
	if err := query.First(&result).Error; err != nil {
		return "", err
	}
	return result.DictLabel, nil
}

func (s *mysqlDictDataStore) SelectDictDataById(ctx context.Context, dictCode int64, opts *api.GetOptions) (*v1.SysDictData, error) {
	query := selectDictDataVo(s.db).Where("dict_code = ?", dictCode)
	query = opts.SQLCondition(query, "")

	var result v1.SysDictData
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(dictCode), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlDictDataStore) CountDictDataByType(ctx context.Context, dictType string, opts *api.GetOptions) int {
	query := selectDictDataVo(s.db).Where("dict_type = ?", dictType)
	query = opts.SQLCondition(query, "")

	var result int64
	query.Count(&result)
	return int(result)
}

func (s *mysqlDictDataStore) DeleteDictDataById(ctx context.Context, dictCode int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("dict_code = ?", dictCode).Delete(&v1.SysDictData{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(dictCode))
}

func (s *mysqlDictDataStore) DeleteDictDataByIds(ctx context.Context, dictCodes []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("dict_code in ?", dictCodes).Delete(&v1.SysDictData{})

	cacheKeys := make([]string, len(dictCodes))
	for i := range dictCodes {
		cacheKeys = append(cacheKeys, s.CacheKey(dictCodes[i]))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return del.Error
}

func (s *mysqlDictDataStore) InsertDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(dictData)
	return create.Error
}

func (s *mysqlDictDataStore) UpdateDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.UpdateOptions) error {
	if dictData.DictCode == 0 {
		return fmt.Errorf("更新字典数据必须指明id")
	}
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(dictData).Where("dict_code = ?", dictData.DictCode).Updates(dictData).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(dictData.DictCode))
}

func (s *mysqlDictDataStore) UpdateDictDataType(ctx context.Context, oldType, newType string, opts *api.UpdateOptions) error {
	var ids []int64
	query := s.db.Table("sys_dict_data").Where("dict_type = ?", oldType).Select("dict_code")
	query.Find(&ids)
	cacheKeys := make([]string, 0, len(ids))
	for i := range ids {
		cacheKeys = append(cacheKeys, s.CacheKey(ids[i]))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)

	update := opts.SQLCondition(s.db).Table("sys_dict_data").Where("dict_type = ?", oldType).Update("dict_type", newType)
	return update.Error
}

var dictDataCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlDictDataStore) CachedDB() *cache.CachedDB {
	dictDataCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-dictdatastore-")
		rdsCli.SetRandomExp(true)

		dictDataCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, dictDataCache.rdsCache)
}

func (s *mysqlDictDataStore) CacheKey(dictDataId int64) string {
	// id
	return fmt.Sprintf("%d", dictDataId)
}
