package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/user823/Sophie/internal/pkg/cache/evict"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/errors"
	"sync"
	"time"
)

type redisCache struct {
	rds            kv.RedisStore
	barrier        SingleFlight
	expiry         time.Duration
	notFoundExpiry time.Duration
	lru            evict.EvictPolicy
	mu             sync.RWMutex
	tmu            sync.RWMutex
	timers         map[string]*time.Timer
}

type StringMarshaler interface {
	Marshal() string
}

type StringUnmarshaler interface {
	Unmarshal(string)
}

type Option func(*Options)
type Options struct {
	expire         time.Duration
	notFoundExpiry time.Duration
	lru            evict.EvictPolicy
}

func NewRedisCache(rds kv.RedisStore, barrier SingleFlight, opts ...Option) Cache {
	options := NewOptions()
	for i := range opts {
		opts[i](options)
	}
	return &redisCache{
		rds:            rds,
		barrier:        barrier,
		expiry:         options.expire,
		notFoundExpiry: options.notFoundExpiry,
		lru:            options.lru,
		timers:         map[string]*time.Timer{},
	}
}

func NewOptions() *Options {
	return &Options{
		expire:         expireTime * time.Second,
		notFoundExpiry: notFoundExpireTime * time.Second,
		// 使用redis自己的缓存更新策略
		lru: evict.EmptyLRU{},
	}
}

func WithExpire(expire time.Duration) Option {
	return func(opt *Options) {
		opt.expire = expire
	}
}

func WithNotFoundExpiry(expire time.Duration) Option {
	return func(opt *Options) {
		opt.notFoundExpiry = expire
	}
}

func WithEvictStrategy(e evict.EvictPolicy) Option {
	return func(opt *Options) {
		opt.lru = e
	}
}

func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		r.mu.Lock()
		r.lru.Remove(key)
		r.mu.Unlock()
		r.rds.DeleteKey(ctx, key)
	}
	return nil
}

func (r *redisCache) IsNotFound(err error) bool {
	switch err {
	case redis.Nil:
		return true
	default:
		return false
	}
}

func (r *redisCache) Get(ctx context.Context, key string, dst any) error {
	if um, ok := dst.(StringUnmarshaler); !ok {
		return fmt.Errorf("dst must implement func [Unmarshal(string)] ")
	} else {
		value, err := r.rds.GetKey(ctx, key)
		if err != nil {
			return err
		}

		um.Unmarshal(value)
		r.mu.Lock()
		defer func() { r.mu.Unlock() }()
		r.lru.Add(key)
		return nil
	}
}

func (r *redisCache) setTimer(key string, expire time.Duration) {
	r.tmu.RLock()
	if t, ok := r.timers[key]; ok {
		r.tmu.RUnlock()
		t.Stop()
		t.Reset(expire)
		return
	}
	r.tmu.RUnlock()

	t := time.AfterFunc(expire, func() {
		r.mu.Lock()
		r.lru.Remove(key)
		r.mu.Unlock()

		r.tmu.Lock()
		delete(r.timers, key)
		r.tmu.Unlock()
	})

	r.tmu.Lock()
	r.timers[key] = t
	r.tmu.Unlock()
}

func (r *redisCache) Set(ctx context.Context, key string, val any) error {
	return r.SetWithExp(ctx, key, val, r.expiry)
}

func (r *redisCache) SetWithExp(ctx context.Context, key string, val any, expire time.Duration) error {
	if m, ok := val.(StringMarshaler); !ok {
		return fmt.Errorf("val must implement func [String() string]")
	} else {
		err := r.rds.SetKey(ctx, key, m.Marshal(), expire.Nanoseconds())
		if err != nil {
			return err
		}
		r.mu.Lock()
		r.lru.Add(key)
		r.mu.Unlock()

		// 过期删除lru key
		if expire > 0 {
			r.setTimer(key, expire)
		}
	}
	return nil
}

func (r *redisCache) Take(ctx context.Context, key string, expire time.Duration, dst any, query func(val any) error) error {
	defer func() {
		r.mu.Lock()
		r.lru.Add(key)
		r.mu.Unlock()
	}()

	// 首先尝试从缓存中获取
	if err := r.Get(ctx, key, dst); err == nil {
		return err
	}

	if expire <= 0 {
		expire = r.expiry
	}

	// 从底层数据库获取
	val, err := r.barrier.Do(key, func() (any, error) {
		// 再次检查缓存
		if err := r.Get(ctx, key, dst); err == nil {
			return dst, nil
		}

		if err := query(dst); err != nil {
			return nil, err
		}
		return dst, nil
	})

	// 从数据库中获取值失败，设置占位符
	if err != nil {
		r.rds.SetKey(ctx, key, placeHolder, r.notFoundExpiry.Nanoseconds())
		r.setTimer(key, r.notFoundExpiry)
		return redis.Nil
	}

	//	从数据库中获取值成功，则设置缓存
	if m, ok := val.(StringMarshaler); ok {
		r.rds.SetKey(ctx, key, m.Marshal(), expire.Nanoseconds())
		r.setTimer(key, expire)
	}
	return nil
}

func (r *redisCache) Clean(ctx context.Context) error {
	r.mu.Lock()
	r.lru.Clean()
	r.mu.Unlock()

	r.tmu.Lock()
	r.timers = map[string]*time.Timer{}
	r.tmu.Unlock()

	// 不能简单调用cleanAll
	keys := r.rds.GetKeys(ctx, "")
	if !r.rds.DeleteKeys(ctx, keys) {
		return errors.New("删除缓存失败")
	}
	return nil
}
