package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cache"
	"github.com/hertz-contrib/cache/persist"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

var cacheWhiteList = map[string]interface{}{}

// 注册需要缓存响应的路由
func RegisterCacheRequestPath(paths ...string) {
	for _, path := range paths {
		cacheWhiteList[path] = struct{}{}
	}
}

func Cache() app.HandlerFunc {
	memoryStore := persist.NewMemoryStore(1 * time.Minute)
	return cache.NewCacheByRequestURIWithIgnoreQueryOrder(
		memoryStore,
		2*time.Second,
		cache.WithCacheStrategyByRequest(func(ctx context.Context, c *app.RequestContext) (bool, cache.Strategy) {
			if _, ok := cacheWhiteList[utils.B2s(c.Request.Path())]; ok {
				return true, cache.Strategy{
					CacheKey: c.Request.URI().String(),
				}
			}
			// 不进行缓存
			return false, cache.Strategy{}
		}),
	)
}
