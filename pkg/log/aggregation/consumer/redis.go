package consumer

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/utils"
)

type RedisConsumer struct {
	store kv.RedisStore
	// 使用redis的分布式锁互斥消费消息
	mu *redsync.Mutex
}

func NewRedisConsumer(store kv.RedisStore) aggregation.RecordConsumer {
	cli, ok := store.LowLevel().(redis.UniversalClient)
	if !ok {
		log.Fatal("get redis client failed")
	}
	pool := goredis.NewPool(cli)
	rs := redsync.New(pool)
	mutex := rs.NewMutex(aggregation.MessageTag)

	store.SetKeyPrefix(aggregation.RecordPrefix)
	return &RedisConsumer{
		store: store,
		mu:    mutex,
	}
}

func (c *RedisConsumer) Connect() bool {
	return c.store.Connected()
}

func (c *RedisConsumer) GetAndDeleteSet(ctx context.Context) []aggregation.LogRecord {
	if err := c.mu.Lock(); err != nil {
		log.Errorf("redis lock error: %s", err.Error())
		return []aggregation.LogRecord{}
	}

	// 解锁，尝试3次
	defer func() {
		if ok, err := c.mu.Unlock(); !ok || err != nil {
			log.Errorf("redis unlock error， trying unlock again")
			attemptTimes := 2
			failTimes := 0
			for failTimes < attemptTimes {
				if ok, err = c.mu.Unlock(); !ok || err != nil {
					failTimes++
					continue
				}
				return
			}
		}
	}()

	vals := c.store.GetAndDeleteSet(ctx, aggregation.Recordkey)
	records := make([]aggregation.LogRecord, len(vals))
	cnt := 0
	for i := range vals {
		var tmp map[string]any

		if err := jsoniter.Unmarshal(utils.S2b(vals[i]), &tmp); err != nil {
			log.Warn("json decode failed: %s", err.Error())
			continue
		}
		var record aggregation.LogRecord
		if err := mapstructure.Decode(tmp, &record); err != nil {
			log.Warn("mapstructure decode failed: %s", err.Error())
			continue
		}

		records[cnt] = record
		cnt++
	}
	records = records[:cnt]
	return records
}

func (c *RedisConsumer) Stop() error {
	return nil
}
