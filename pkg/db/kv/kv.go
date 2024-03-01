package kv

import (
	"context"
	"github.com/user823/Sophie/pkg/db/kv/redis"
)

// KeyValueStore 是所有k/v数据库存储后端的标准接口
type KeyValueStore interface {
	Connect(ctx context.Context, config any)
	Connected() bool
	Disconnect() error
	GetKey(context.Context, string) (string, error)          // 获取key 对应的值。如果不存在返回ErrkeyNotFound
	GetMultiKey(context.Context, []string) ([]string, error) // 如果所有键都没有值，则返回ErrkeyNotFound
	SetKey(context.Context, string, string, int64) error     // 设置键值，过期时间单位是纳秒
	SetExp(context.Context, string, int64) error
	GetExp(context.Context, string) (int64, error)
	GetKeys(context.Context, string) []string // 不要对匹配模式filter进行hash
	DeleteKey(context.Context, string) bool
	DeleteAllKeys(context.Context) bool
	GetKeysAndValues(context.Context) map[string]string
	GetKeysAndValuesWithFilter(context.Context, string) map[string]string
	DeleteKeys(context.Context, []string) bool
	Decrement(context.Context, string) int64
	IncrememntWithExpire(context.Context, string, int64) int64
	SetRollingWindow(ctx context.Context, key string, per int64, val string, pipeline bool) (int, []interface{})
	GetRollingWindow(ctx context.Context, key string, per int64, pipeline bool) (int, []interface{})
	GetSet(context.Context, string) (map[string]string, error)
	AddToSet(context.Context, string, string)
	GetAndDeleteSet(context.Context, string) []interface{}
	RemoveFromSet(context.Context, string, string)
	DeleteScanMatch(context.Context, string) bool
	AddToSortedSet(context.Context, string, string, float64)
	GetSortedSetRange(context.Context, string, string, string) ([]string, []float64, error)
	RemoveSortedSetRange(context.Context, string, string, string) error
	GetListRange(context.Context, string, int64, int64) ([]string, error)
	RemoveFromList(context.Context, string, string) error
	AppendToSet(context.Context, string, string)
	Exists(context.Context, string) (bool, error)
}

// redis 一般要结合keyprefix 和 hash一起使用
type RedisStore interface {
	KeyValueStore
	SetRandomExp(bool)
	SetKeyPrefix(string)
	SetHashKey(bool)
	GetAndDelete(context.Context, string) (string, error)
	GetKeyTTL(context.Context, string) (int64, error)
	GetRawKey(context.Context, string) (string, error)
	SetRawKey(context.Context, string, string, int64) error
	DeleteRawKey(context.Context, string) bool
	AppendToSetPipelined(context.Context, string, []string)
	Publish(context.Context, string, string) error
	StartPubSubHandler(context.Context, string, func(any)) error
}

func NewKVStore(name string) KeyValueStore {
	switch name {
	case "redis":
		return redis.NewRedisClient()
	default:
		return redis.NewRedisClient()
	}
}
