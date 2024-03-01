package cache

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// mysql缓存层
// 使用db执行的动作
type (
	ExecFn func(ctx context.Context, db *gorm.DB) error
	// 基于唯一索引进行查询, v 需要是指针
	// 唯一索引未找到的情况下返回 主键
	IndexQueryFn func(ctx context.Context, db *gorm.DB, v any) (any, error)
	// 基于主键进行查询
	PrimaryQueryFn func(ctx context.Context, db *gorm.DB, primary any, v any) error
	// 普通查询
	QueryFn func(ctx context.Context, db *gorm.DB, v any) error
)

type CachedDB struct {
	db    *gorm.DB
	cache Cache
}

func NewCachedDB(db *gorm.DB, cache Cache) *CachedDB {
	return &CachedDB{
		db:    db,
		cache: cache,
	}
}

func (c *CachedDB) DB() *gorm.DB {
	return c.db
}

// 删除缓存层中的key
func (c *CachedDB) DelCache(ctx context.Context, keys ...string) error {
	return c.cache.Del(ctx, keys...)
}

// 从缓存中获取key
func (c *CachedDB) GetCache(ctx context.Context, key string, dst any) error {
	return c.cache.Get(ctx, key, dst)
}

// 清空缓存
func (s *CachedDB) CleanCache(ctx context.Context) error {
	return s.cache.Clean(ctx)
}

// 执行Exec动作， 并且删除相关的keys
func (c *CachedDB) Exec(ctx context.Context, execFn ExecFn, keys ...string) error {
	if err := execFn(ctx, c.db); err != nil {
		return err
	}

	return c.cache.Del(ctx, keys...)
}

// 使用读缓存策略执行给定查询
func (c *CachedDB) QueryRow(ctx context.Context, key string, v any, query QueryFn) error {
	return c.cache.Take(ctx, key, expireTime*time.Second, v, func(v any) error {
		return query(ctx, c.db, v)
	})
}

// 使用读缓存策略执行唯一索引上的查询
func (c *CachedDB) QueryRowIndex(ctx context.Context, key string, v any, keyer func(primaryKey any) string, indexQuery IndexQueryFn, primaryQuery PrimaryQueryFn) error {
	var primaryKey any
	var found bool

	// 首先查询唯一索引
	if err := c.cache.Take(ctx, key, expireTime*time.Second, v, func(v any) (err error) {
		primaryKey, err = indexQuery(ctx, c.db, v)
		if err != nil {
			return err
		}
		found = true
		return c.cache.Set(ctx, keyer(primaryKey), v)
	}); err != nil {
		return err
	}
	// v 值已经设置
	if found {
		return nil
	}

	// 唯一索引未找到值， 从主键索引查询
	return c.cache.Take(ctx, keyer(primaryKey), expireTime*time.Second, v, func(val any) error {
		return primaryQuery(ctx, c.db, primaryKey, v)
	})
}
