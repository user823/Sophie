package cacheutils

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/pkg/cache"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/store/mysql"
	"github.com/user823/Sophie/pkg/db/kv"
	"sync"
)

var dictCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

// 加载字典缓存数据
func LoadingDictCache(s store.Factory) {
	// 初始化
	dictCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis").(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-dictcache-")
		rdsCli.SetRandomExp(true)

		dictCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	if s == nil {
		s, _ = mysql.GetMySQLFactoryOr(nil)
	}
	ctx := context.Background()
	sysDictDatas, _ := s.DictData().SelectDictDataList(ctx, &v1.SysDictData{Status: "0"}, &api.GetOptions{})

	// 聚合
	mp := make(map[string]*v1.DictDataList)
	for i := range sysDictDatas {
		key := sysDictDatas[i].DictType
		if _, ok := mp[key]; !ok {
			mp[key] = &v1.DictDataList{}
		}
		mp[key].TotalCount++
		mp[key].Items = append(mp[key].Items, sysDictDatas[i])
	}

	// 缓存
	for k, v := range mp {
		dictCache.rdsCache.Set(ctx, k, v)
	}
}

// 清空字段缓存数据
func CleanDictCache() {
	dictCache.rdsCache.Clean(context.Background())
}

// 删除指定字典缓存
func RemoveDictCache(key string) {
	dictCache.rdsCache.Del(context.Background(), key)
}

// 设置字典缓存
func SetDictCache(key string, dictDatas []*v1.SysDictData) {
	dictDataList := &v1.DictDataList{
		ListMeta: api.ListMeta{int64(len(dictDatas))},
		Items:    dictDatas,
	}
	dictCache.rdsCache.Set(context.Background(), key, dictDataList)
}

// 获取字典缓存
func GetDictCache(key string) []*v1.SysDictData {
	var result v1.DictDataList
	err := dictCache.rdsCache.Get(context.Background(), key, &result)
	if err != nil {
		return []*v1.SysDictData{}
	}
	return result.Items
}
