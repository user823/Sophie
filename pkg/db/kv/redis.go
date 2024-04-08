package kv

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/hash"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type RedisConfig struct {
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	Database              int
	MaxIdle               int
	MaxActive             int
	Timeout               int
	EnableCluster         bool
	UseSSL                bool
	SSLInsecureSkipVerify bool
}

// 单例模式
type RedisClient struct {
	KeyPrefix string
	HashKey   bool
	RandomExp bool
}

var ErrRedisIsDown = errors.New("storage: Redis is either down or ws not configured")

// 单例模式
var (
	singlePool atomic.Value
	// 关闭redis
	redisUp atomic.Value
)

const (
	// 附加随机有效期占基础有效期的百分比
	RANDOM_PERCENT = 0.5
)

// true 表示禁用redis
var disableRedis atomic.Value
var (
	ErrConfigTypeInvalid = errors.New("cannot correctly convert the config type")
	ErrKeyNotFound       = errors.New("key not found")
)

func DisableRedis(disable bool) {
	if disable {
		disableRedis.Store(true)
		redisUp.Store(false)
		return
	}
	disableRedis.Store(false)
	redisUp.Store(true)
}

// 判断是否允许连接到redis
func shouldConnect() bool {
	if v := disableRedis.Load(); v != nil {
		return !v.(bool)
	}
	return true
}

func Connected() bool {
	if v := redisUp.Load(); v != nil {
		return v.(bool)
	}
	return false
}

func singleton() redis.UniversalClient {
	if v := singlePool.Load(); v != nil {
		return v.(redis.UniversalClient)
	}
	return nil
}

func connectSingleton(config *RedisConfig) bool {
	if singleton() == nil {
		singlePool.Store(ConnectToRedis(config))
	}
	return true
}

// 获取redis连接
func ConnectToRedis(config *RedisConfig) redis.UniversalClient {
	log.Debug("Creating new Redis connection pool")

	poolSize := 500
	if config.MaxActive > 0 {
		poolSize = config.MaxActive
	}

	timeout := 5 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	maxIdle := 240 * timeout
	if config.MaxIdle > 0 {
		maxIdle = time.Duration(config.MaxIdle) * time.Second
	}

	// tls 设置
	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient
	opts := &redis.UniversalOptions{
		Addrs:           getRedisAddrs(config),
		MasterName:      config.MasterName,
		Password:        config.Password,
		DB:              config.Database,
		DialTimeout:     timeout,
		ReadTimeout:     timeout,
		WriteTimeout:    timeout,
		ConnMaxIdleTime: maxIdle,
		PoolSize:        poolSize,
		TLSConfig:       tlsConfig,
	}

	if opts.MasterName != "" {
		log.Info("--> [REDIS] Creating sentinel-backed failover client")
		client = redis.NewFailoverClient(opts.Failover())
	} else if config.EnableCluster {
		log.Info("--> [REDIS] Creating cluster client")
		client = redis.NewClusterClient(opts.Cluster())
	} else {
		log.Info("--> [REDIS] Creating single-node client")
		client = redis.NewClient(opts.Simple())
	}
	return client
}

func getRedisAddrs(config *RedisConfig) []string {
	if len(config.Addrs) != 0 {
		return config.Addrs
	}

	if config.MasterName != "" {
		return []string{"127.0.0.1:26379"}
	}
	return []string{"127.0.0.1:6379"}
}

func connectionIsOpen() bool {
	c := singleton()
	if c == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	testKey := "redis-conn-test-" + uuid.New().String()
	if err := c.Set(ctx, testKey, "test", 1*time.Second).Err(); err != nil {
		log.Warnf("Error trying to set test key: %s", err.Error())
		return false
	}

	if _, err := c.Get(ctx, testKey).Result(); err != redis.Nil && err != nil {
		log.Warnf("Error trying to get test key: %s", err.Error())
		return false
	}

	return true
}

// 开启协程保持长连接状态
func KeepConnection(ctx context.Context, config *RedisConfig) {
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	connectSingleton(config)
	redisUp.Store(connectionIsOpen())
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if !shouldConnect() {
				continue
			}
			connectSingleton(config)
			redisUp.Store(connectionIsOpen())
		}
	}
}

func Disconnect() error {
	if v, ok := singlePool.Load().(redis.UniversalClient); v != nil && ok {
		if err := v.Close(); err != nil {
			return err
		}
	}
	return nil
}

func NewRedisClient() *RedisClient {
	return &RedisClient{}
}

func (r *RedisClient) Connect(ctx context.Context, config any) {
	// 如果存在连接则默认使用现有的连接
	connectSingleton(config.(*RedisConfig))
}

func (r *RedisClient) Connected() bool {
	return Connected()
}

// 不允许关闭单例的连接
func (r *RedisClient) Disconnect() error { return nil }

func (r *RedisClient) SetKeyPrefix(prefix string) {
	r.KeyPrefix = prefix
}

func (r *RedisClient) SetHashKey(ok bool) {
	r.HashKey = ok
}

func (r *RedisClient) SetRandomExp(ok bool) {
	r.RandomExp = ok
}

func (r *RedisClient) cacheKey(keyname string) string {
	return r.KeyPrefix + r.hashKey(keyname)
}

func (r *RedisClient) hashKey(keyname string) string {
	if !r.HashKey {
		return keyname
	}
	return hash.NewHasher(hash.DefaultHashAlgorithm).HashKey(utils.S2b(keyname))
}

func (r *RedisClient) cleanKey(key string) string {
	return strings.Replace(key, r.KeyPrefix, "", 1)
}

func (r *RedisClient) GetKey(ctx context.Context, keyname string) (string, error) {
	if !Connected() {
		return "", ErrRedisIsDown
	}

	c := singleton()

	value, err := c.Get(ctx, r.cacheKey(keyname)).Result()
	if err != nil {
		log.Debugf("Error trying to get value: %s", err.Error())

		return "", ErrKeyNotFound
	}
	return value, nil
}

func (r *RedisClient) GetMultiKey(ctx context.Context, keys []string) ([]string, error) {
	if !Connected() {
		return nil, ErrRedisIsDown
	}

	keynames := make([]string, len(keys))
	for i, val := range keys {
		keynames[i] = r.cacheKey(val)
	}

	c := singleton()
	values, err := c.MGet(ctx, keynames...).Result()
	if err != nil {
		log.Debugf("Error trying to get value: %s", err.Error())

		return nil, ErrKeyNotFound
	}
	result := make([]string, len(values))
	found := false
	for i, val := range values {
		strVal := fmt.Sprint(val)
		if strVal == "<nil>" {
			strVal = ""
		}
		if strVal != "" {
			found = true
		}
		result[i] = strVal
	}

	if found {
		return result, nil
	}
	return nil, ErrKeyNotFound
}

func (r *RedisClient) GetKeyTTL(ctx context.Context, keyname string) (int64, error) {
	if !Connected() {
		return 0, ErrRedisIsDown
	}
	c := singleton()
	duration, err := c.TTL(ctx, r.cacheKey(keyname)).Result()
	return int64(duration.Seconds()), err
}

func (r *RedisClient) GetRawKey(ctx context.Context, keyName string) (string, error) {
	if !Connected() {
		return "", ErrRedisIsDown
	}
	value, err := singleton().Get(ctx, r.KeyPrefix+keyName).Result()
	if err != nil {
		log.Debugf("Error trying to get value: %s", err.Error())

		return "", ErrKeyNotFound
	}

	return value, nil
}

func (r *RedisClient) GetExp(ctx context.Context, keyName string) (int64, error) {
	if !Connected() {
		return 0, ErrRedisIsDown
	}
	value, err := singleton().TTL(ctx, r.cacheKey(keyName)).Result()
	if err != nil {
		log.Errorf("Error trying to get TTL: ", err.Error())

		return 0, ErrKeyNotFound
	}

	return int64(value.Seconds()), nil
}

func (r *RedisClient) GetKeys(ctx context.Context, filter string) []string {
	if !Connected() {
		return nil
	}

	c := singleton()
	searchStr := r.KeyPrefix + filter + "*"
	log.Debugf("[STORE] Getting list by: %s", searchStr)

	fetchKeys := func(ctx context.Context, client *redis.Client) ([]string, error) {
		values := make([]string, 0)
		iter := client.Scan(ctx, 0, searchStr, 0).Iterator()
		for iter.Next(ctx) {
			values = append(values, iter.Val())
		}
		if err := iter.Err(); err != nil {
			return nil, err
		}

		return values, nil
	}

	var err error
	sessions := make([]string, 0)

	switch v := c.(type) {
	case *redis.ClusterClient:
		ch := make(chan []string)

		go func() {
			// ForEachMaster 并发在每个主节点上调用fn
			err = v.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
				results, e := fetchKeys(ctx, client)
				if e != nil {
					return nil
				}
				ch <- results
				return nil
			})
			close(ch)
		}()

		for res := range ch {
			sessions = append(sessions, res...)
		}
	case *redis.Client:
		sessions, err = fetchKeys(ctx, v)
	}

	if err != nil {
		log.Errorf("Error while fetching keys: %s", err)
		return nil
	}
	for i, v := range sessions {
		sessions[i] = r.cleanKey(v)
	}
	return sessions
}

func (r *RedisClient) DeleteKey(ctx context.Context, key string) bool {
	if !Connected() {
		return false
	}

	log.Debugf("DEL Key was: %s", key)
	log.Debugf("DEL Key became: %s", r.cacheKey(key))
	c := singleton()
	n, err := c.Del(ctx, r.cacheKey(key)).Result()
	if err != nil {
		log.Errorf("Error trying to delete key: %s", err.Error())
	}
	return n > 0
}

func (r *RedisClient) DeleteAllKeys(ctx context.Context) bool {
	if !Connected() {
		return false
	}

	n, err := singleton().FlushAll(ctx).Result()
	if err != nil {
		log.Errorf("Error try to delete keys: %s", err.Error())
	}
	return n == "OK"
}

func (r *RedisClient) DeleteRawKey(ctx context.Context, key string) bool {
	if !Connected() {
		return false
	}

	log.Debugf("DEL Key was: %s", key)
	c := singleton()
	n, err := c.Del(ctx, r.KeyPrefix+key).Result()
	if err != nil {
		log.Errorf("Error trying to delete key: %s", err.Error())
	}
	return n > 0
}

func (r *RedisClient) GetKeysAndValuesWithFilter(ctx context.Context, filter string) map[string]string {
	if !Connected() {
		return nil
	}

	searchStr := r.KeyPrefix + filter + "*"

	c := singleton()
	fetchKeysAndValues := func(ctx context.Context, client *redis.Client) (map[string]string, error) {
		result := map[string]string{}
		iter := client.Scan(ctx, 0, searchStr, 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			value, err := client.Get(ctx, key).Result()
			if err != nil && !errors.Is(err, redis.Nil) {
				return nil, err
			}
			result[r.cleanKey(key)] = value
		}
		if err := iter.Err(); err != nil {
			return nil, err
		}

		return result, nil
	}

	result := make(map[string]string)
	var err error
	switch v := c.(type) {
	case *redis.ClusterClient:
		ch := make(chan map[string]string)

		go func() {
			err = v.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
				value, e := fetchKeysAndValues(ctx, client)
				if e != nil {
					return e
				}
				ch <- value
				return nil
			})
			close(ch)
		}()
		for res := range ch {
			if len(result) < len(res) {
				result, res = res, result
			}
			for k, v := range res {
				result[k] = v
			}
		}
	case *redis.Client:
		result, err = fetchKeysAndValues(ctx, v)
	}
	if err != nil {
		log.Debugf("Error trying to fetch keys and values: %s", err.Error())
		return nil
	}
	return result
}

func (r *RedisClient) GetKeysAndValues(ctx context.Context) map[string]string {
	return r.GetKeysAndValuesWithFilter(ctx, "")
}

func (r *RedisClient) DeleteKeys(ctx context.Context, keys []string) bool {
	if !Connected() {
		return false
	}

	if len(keys) > 0 {
		for i, v := range keys {
			keys[i] = r.cacheKey(v)
		}

		log.Debugf("Deleting: %v", keys)
		c := singleton()
		switch v := c.(type) {
		case *redis.ClusterClient:
			pipe := v.Pipeline()
			for _, k := range keys {
				pipe.Del(ctx, k)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				log.Errorf("Error trying to delete keys: %s", err.Error())
			}

		case *redis.Client:
			_, err := v.Del(ctx, keys...).Result()
			if err != nil {
				log.Errorf("Error trying to delete keys: %s", err.Error())
			}
		}
	}
	return true
}

func (r *RedisClient) Decrement(ctx context.Context, key string) int64 {
	if !Connected() {
		return 0
	}

	log.Debugf("Deleting key: %s", key)
	n, err := singleton().Decr(ctx, r.cacheKey(key)).Result()
	if err != nil {
		log.Debugf("Error deleting key: %s", err.Error())
	}
	return n
}

func (r *RedisClient) IncrementWithExpire(ctx context.Context, key string, expire int64) int64 {
	log.Debugf("Incrementing raw key: %s", key)
	if !Connected() {
		return 0
	}

	cachekey := r.cacheKey(key)
	val, err := singleton().Incr(ctx, cachekey).Result()
	if err != nil {
		log.Errorf("Error trying to increment value: %s", err.Error())
	} else {
		log.Debugf("Incremented key: %s, val is: %d", cachekey, val)
	}

	if val == 1 && expire > 0 {
		log.Debug("--> Setting Expire")
		if r.RandomExp {
			expire += int64(float64(expire) * RANDOM_PERCENT)
		}
		singleton().Expire(ctx, cachekey, time.Duration(expire)*time.Second)
	}
	return val
}

// 添加一个有序集合到redis中， 并且获取过期时间窗口的值
// per 以毫秒为单位
func (r *RedisClient) SetRollingWindow(ctx context.Context, key string, per int64, val string, pipeline bool) (int, []string) {
	log.Debugf("Incrementing raw key: %s", key)
	if !Connected() {
		return 0, nil
	}
	cacheKey := r.cacheKey(key)

	log.Debugf("keyName is: %s", cacheKey)
	now := time.Now()
	log.Debugf("Now is: %v", now)
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Millisecond)
	log.Debugf("Then is: %v", onePeriodAgo)

	c := singleton()
	var zrange *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, cacheKey, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		zrange = pipe.ZRange(ctx, cacheKey, 0, -1)

		element := redis.Z{
			Score: float64(now.UnixNano()),
		}

		if val != "-1" {
			element.Member = val
		} else {
			element.Member = strconv.Itoa(int(now.UnixNano()))
		}

		pipe.ZAdd(ctx, cacheKey, element)
		pipe.Expire(ctx, cacheKey, time.Duration(per)*time.Millisecond)
		return nil
	}

	var err error
	if pipeline {
		_, err = c.Pipelined(ctx, pipeFn)
	} else {
		_, err = c.TxPipelined(ctx, pipeFn)
	}

	if err != nil {
		log.Errorf("Multi command failed: %s", err.Error())

		return 0, nil
	}
	values := zrange.Val()
	if values == nil {
		return 0, nil
	}

	intVal := len(values)
	log.Debugf("Returned: %d", intVal)
	return intVal, values
}

func (r *RedisClient) GetRollingWindow(ctx context.Context, key string, per int64, pipeline bool) (int, []string) {
	if !Connected() {
		return 0, nil
	}
	now := time.Now()
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Millisecond)
	cacheKey := r.cacheKey(key)

	c := singleton()
	var zrange *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, cacheKey, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		zrange = pipe.ZRange(ctx, cacheKey, 0, -1)
		return nil
	}

	var err error
	if pipeline {
		_, err = c.Pipelined(ctx, pipeFn)
	} else {
		_, err = c.TxPipelined(ctx, pipeFn)
	}
	if err != nil {
		log.Errorf("Multi command failed: %s", err.Error())

		return 0, nil
	}

	values := zrange.Val()
	if values == nil {
		return 0, nil
	}
	intVal := len(values)
	log.Debugf("Returned: %d", intVal)

	return intVal, values
}

func (r *RedisClient) GetSet(ctx context.Context, key string) (map[string]string, error) {
	if !Connected() {
		return nil, ErrRedisIsDown
	}

	log.Debugf("Getting from key set: %s", key)
	val, err := singleton().SMembers(ctx, r.cacheKey(key)).Result()
	if err != nil {
		log.Errorf("Error trying to get key set: %s", err.Error())
		return nil, err
	}

	result := make(map[string]string)
	for i, value := range val {
		result[strconv.Itoa(i)] = value
	}
	return result, nil
}

func (r *RedisClient) AddToSet(ctx context.Context, key string, value string) {
	if !Connected() {
		return
	}

	log.Debugf("Pushing to raw key set: %s", key)
	if err := singleton().SAdd(ctx, r.cacheKey(key), value).Err(); err != nil {
		log.Errorf("Error trying to append keys: %s", err.Error())
	}
}

func (r *RedisClient) GetAndDeleteSet(ctx context.Context, key string) []string {
	if !Connected() {
		return nil
	}

	log.Debugf("Geting and deleting form key set: %s", key)
	cacheKey := r.cacheKey(key)
	c := singleton()

	var lrange *redis.StringSliceCmd
	_, err := c.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		lrange = pipe.LRange(ctx, cacheKey, 0, -1)
		pipe.Del(ctx, cacheKey)
		return nil
	})

	if err != nil {
		log.Errorf("Multi command failed: %s", err.Error())
		return nil
	}

	vals := lrange.Val()
	log.Debugf("Analytics returned: %d", len(vals))
	if len(vals) == 0 {
		return nil
	}
	fmt.Println(len(vals))

	log.Debugf("Unpacked vals: %d", len(vals))
	return vals
}

func (r *RedisClient) RemoveFromSet(ctx context.Context, key string, value string) {
	if !Connected() {
		return
	}
	log.Debugf("Deleting key and value form set: %s %s", key, value)
	if err := singleton().SRem(ctx, r.cacheKey(key), value).Err(); err != nil {
		log.Errorf("Error trying to remove keys: %s", err.Error())
	}
}

func (r *RedisClient) DeleteScanMatch(ctx context.Context, filter string) bool {
	if !Connected() {
		return false
	}

	c := singleton()
	keys := r.GetKeys(ctx, filter)
	cnt := 0
	for _, key := range keys {
		log.Infof("Deleting: %s", key)
		if err := c.Del(ctx, key).Err(); err != nil {
			log.Errorf("Error trying to delete key: %s - %s", key, err.Error())
			continue
		}
		cnt++
	}
	log.Infof("Deleted %d keys", cnt)
	return true
}

func (r *RedisClient) AddToSortedSet(ctx context.Context, key string, value string, score float64) {
	if !Connected() {
		return
	}

	cacheKey := r.cacheKey(key)
	log.Debug("Pushing key and value to sorted set: %s %s", key, value)
	member := redis.Z{Score: score, Member: value}
	if err := singleton().ZAdd(ctx, cacheKey, member).Err(); err != nil {
		log.Errorw(
			"ZADD command failed",
			"keyName", key,
			"error", err.Error(),
		)
	}
}

func (r *RedisClient) GetSortedSetRange(ctx context.Context, key string, from string, to string) ([]string, []float64, error) {
	if !Connected() {
		return nil, nil, ErrRedisIsDown
	}
	cacheKey := r.cacheKey(key)
	log.Debugw("Getting sorted set range", "keyname", key, "scoreForm", from, "scoreTo", to)
	args := redis.ZRangeBy{Min: from, Max: to}
	values, err := singleton().ZRangeByScoreWithScores(ctx, cacheKey, &args).Result()
	if err != nil {
		log.Errorw("ZRANGEBYSCORE command failed", "keyname", key, "scoreForm", from, "scoreTo", to, "error", err.Error())
	}
	if len(values) == 0 {
		return nil, nil, nil
	}

	elements := make([]string, len(values))
	scores := make([]float64, len(values))

	for i, v := range values {
		elements[i] = fmt.Sprint(v.Member)
		scores[i] = v.Score
	}

	return elements, scores, nil
}

func (r *RedisClient) RemoveSortedSetRange(ctx context.Context, key string, from string, to string) error {
	if !Connected() {
		return ErrRedisIsDown
	}
	cacheKey := r.cacheKey(key)
	log.Debugw("Removing sorted set range", "keyname", key, "scoreFrom", from, "scoreTo", to)
	if err := singleton().ZRemRangeByScore(ctx, cacheKey, from, to).Err(); err != nil {
		log.Debugw(
			"ZREMRANGEBYSCORE command failed",
			"keyName", key,
			"fixedKey", cacheKey,
			"scoreFrom", from,
			"scoreTo", to,
			"error", err.Error(),
		)

		return err
	}

	return nil
}

func (r *RedisClient) GetListRange(ctx context.Context, key string, from int64, to int64) ([]string, error) {
	if !Connected() {
		return nil, ErrRedisIsDown
	}
	log.Debugf("Getting key form list: %s %d %d", key, from, to)
	values, err := singleton().LRange(ctx, r.cacheKey(key), from, to).Result()
	if err != nil {
		log.Debugw("Error getting key form list", "keyname", key, "error", err.Error())
		return nil, err
	}
	return values, nil
}

// 从列表中删除一个值
func (r *RedisClient) RemoveFromList(ctx context.Context, key string, value string) error {
	if !Connected() {
		return ErrRedisIsDown
	}
	log.Debugf("Deleting form list: %s", key)
	if err := singleton().LRem(ctx, r.cacheKey(key), 0, value).Err(); err != nil {
		log.Errorw("Error deleting value from list", "keyname", key, "error", err.Error())
		return err
	}
	return nil
}

// 添加值到列表中
func (r *RedisClient) AppendToList(ctx context.Context, key string, value string) {
	if !Connected() {
		return
	}

	log.Debug("Pushing to list: %s", key)
	if err := singleton().RPush(ctx, r.cacheKey(key), value).Err(); err != nil {
		log.Errorf("Error trying to append to set keys: %s", err.Error())
	}
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	if !Connected() {
		return false, ErrRedisIsDown
	}
	log.Debugw("Checking if exists", "keyName", key)
	exists, err := singleton().Exists(ctx, r.cacheKey(key)).Result()
	if err != nil {
		log.Errorf("Error trying to check if key exists: %s", err.Error())
		return false, err
	}
	if exists == 1 {
		return true, nil
	}

	return false, nil
}

func (r *RedisClient) SetExp(ctx context.Context, key string, expire int64) error {
	if !Connected() {
		return ErrRedisIsDown
	}

	log.Debugw("Trying to set expire for key: %s", key)
	if r.RandomExp {
		expire += int64(float64(expire) * RANDOM_PERCENT)
	}
	if err := singleton().Expire(ctx, r.cacheKey(key), time.Duration(expire)).Err(); err != nil {
		log.Errorf("Could not EXPIRE key: %s", err.Error())
		return err
	}
	return nil
}

func (r *RedisClient) SetKey(ctx context.Context, key string, value string, expire int64) error {
	if !Connected() {
		return ErrRedisIsDown
	}

	log.Debugw("Trying to set key: %s", key)
	if r.RandomExp {
		expire += int64(float64(expire) * RANDOM_PERCENT)
	}
	if err := singleton().Set(ctx, r.cacheKey(key), value, time.Duration(expire)).Err(); err != nil {
		log.Errorw("Error to set key", "keyname", key, "error", err.Error())
		return err
	}
	return nil
}

func (r *RedisClient) SetRawKey(ctx context.Context, key string, value string, expire int64) error {
	if !Connected() {
		return ErrRedisIsDown
	}

	cacheKey := r.KeyPrefix + key
	log.Debugw("Trying to set key: %s", key)
	if r.RandomExp {
		expire += int64(float64(expire) * RANDOM_PERCENT)
	}
	if err := singleton().Set(ctx, cacheKey, value, time.Duration(expire)).Err(); err != nil {
		log.Errorw("Error to set key", "keyname", key, "error", err.Error())
		return err
	}
	return nil
}

func (r *RedisClient) AppendToSetPipelined(ctx context.Context, key string, values []string) {
	if len(values) == 0 {
		return
	}
	if !r.Connected() {
		log.Debugln(ErrRedisIsDown)
		return
	}
	cacheKey := r.cacheKey(key)
	c := singleton()
	pipe := c.Pipeline()
	for _, val := range values {
		pipe.RPush(ctx, cacheKey, val)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		log.Errorf("Error trying to append to set keys: %s", err.Error())
	}

	if storageExpTime := int64(viper.GetDuration("logRecord.storage_expiration_time")); storageExpTime >= 0 {
		if storageExpTime == 0 {
			storageExpTime = EXPIRATION
		}
		// If there is no expiry on the analytics set, we should set it.
		exp, _ := r.GetExp(ctx, cacheKey)
		if exp == -1 {
			_ = r.SetExp(ctx, cacheKey, utils.SecondToNano(storageExpTime))
		}
	}
}

func (r *RedisClient) Publish(ctx context.Context, channel string, message string) error {
	if !Connected() {
		return ErrRedisIsDown
	}

	if err := singleton().Publish(ctx, channel, message).Err(); err != nil {
		log.Errorf("Error trying to set value: %s", err.Error())
		return err
	}
	return nil
}

func (r *RedisClient) StartPubSubHandler(ctx context.Context, channel string, callback func(any)) error {
	if !Connected() {
		return ErrRedisIsDown
	}

	pubsub := singleton().Subscribe(ctx, channel)
	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
		log.Errorf("Error while receiving pubsub message: %s", err.Error())

		return err
	}

	for msg := range pubsub.Channel() {
		callback(msg)
	}

	return nil
}

func (r *RedisClient) GetAndDelete(ctx context.Context, key string) (string, error) {
	if !Connected() {
		return "", ErrRedisIsDown
	}
	c := singleton()
	cacheKey := r.cacheKey(key)
	var cmd *redis.StringCmd

	_, err := c.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		cmd = pipe.Get(ctx, cacheKey)
		pipe.Del(ctx, cacheKey)
		return nil
	})

	return cmd.Val(), err
}

func (r *RedisClient) AddToHash(ctx context.Context, key string, val map[string]any) error {
	if !Connected() {
		return ErrRedisIsDown
	}
	c := singleton()
	cacheKey := r.cacheKey(key)
	return c.HSet(ctx, cacheKey, val).Err()
}

func (r *RedisClient) MGetFromHash(ctx context.Context, keys []string) ([]map[string]string, error) {
	if !Connected() {
		return []map[string]string{}, ErrRedisIsDown
	}
	c := singleton()

	var results []map[string]string
	for _, key := range keys {
		cacheKey := r.cacheKey(key)
		result, err := c.HGetAll(ctx, cacheKey).Result()
		// 需要判断空
		if err != nil || len(result) == 0 {
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *RedisClient) RemoveSortedSet(ctx context.Context, key string, members ...any) {
	if !Connected() {
		return
	}
	c := singleton()
	cacheKey := r.cacheKey(key)
	if err := c.ZRem(ctx, cacheKey, members...).Err(); err != nil {
		log.Errorf("Remove sorted set %s members error: %s", key, err.Error())
	}
}

func (r *RedisClient) LowLevel() any {
	return singleton()
}
