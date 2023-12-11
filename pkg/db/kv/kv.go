package db

import (
	"context"
	"errors"
	"github.com/user823/Sophie/pkg/db/redis"
	"github.com/user823/Sophie/pkg/log"
)

var ErrKeyNotFound = errors.New("key not found")

// KeyValueStore 是所有k/v数据库存储后端的标准接口
type KeyValueStore interface {
	Connect(ctx context.Context, config any)
	Connected() bool
	Disconnect() error
	GetKey(context.Context, string) (string, error)
	GetMultiKey(context.Context, []string) ([]string, error)
	SetKey(context.Context, string, string, int64) error
	SetExp(context.Context, string, int64) error
	GetExp(context.Context, string) (int64, error)
	GetKeys(context.Context, string) []string
	DeleteKey(context.Context, string) bool
	DeleteAllKeys(context.Context) bool
	GetKeysAndValues(context.Context) map[string]string
	GetKeysAndValuesWithFilter(context.Context, string) map[string]string
	DeleteKeys(context.Context, []string) bool
	Decrement(context.Context, string)
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

type RedisStore interface {
	KeyValueStore
	SetKeyPrefix(string)
	SetHashKey(bool)
	GetKeyPrefix() string
	GetKeyTTL(context.Context, string) (int64, error)
	GetRawKey(context.Context, string) (string, error)
	SetRawKey(context.Context, string, string, int64) error
	DeleteRawKey(context.Context, string) bool
}

func NewKVStore(name string) KeyValueStore {
	switch name {
	case "redis":
		return redis.NewRedisClient()
	default:
		log.Warn("Cannot match any key-value instance. return redisclient")
		return redis.NewRedisClient()
	}
}
