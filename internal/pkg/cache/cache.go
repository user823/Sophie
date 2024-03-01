package cache

import (
	"context"
	"time"
)

// 缓存需要实现的接口
type Cache interface {
	// 删除多个缓存
	Del(ctx context.Context, keys ...string) error
	// 从缓存中取数据
	Get(ctx context.Context, key string, dst any) error
	Set(ctx context.Context, key string, val any) error
	SetWithExp(ctx context.Context, key string, val any, expire time.Duration) error
	// 从缓存中取数据，如果未获取成功，则调用query从DB获取数据，并使用指定过期时间(s)设置缓存，（如果指定时间小于等于0则使用默认缓存时间）
	// 如果未成功获取数据则采用palceHolder占位
	Take(ctx context.Context, key string, expire time.Duration, dst any, query func(val any) error) error
	// 判断给定的error是否为errNotFound
	IsNotFound(err error) bool
	// 清空缓存
	Clean(ctx context.Context) error
}

const (
	// 默认过期时间为 1h
	expireTime = 3600
	// 默认未找到时过期时间 10 分钟
	notFoundExpireTime = 600
	// 未找到占位符
	placeHolder = "not found"
)
