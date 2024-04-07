package producer

import (
	"context"
	"time"

	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

type RedisProducer struct {
	store      kv.RedisStore
	expireTime time.Duration
}

func NewRedisProducer(store kv.RedisStore, expireTime time.Duration) aggregation.RecordProducer {
	store.SetKeyPrefix(aggregation.RecordPrefix)
	return &RedisProducer{store, expireTime}
}

func (p *RedisProducer) Connect() bool {
	return p.store.Connected()
}

func (p *RedisProducer) AppendToSet(ctx context.Context, msg []string) {
	p.store.AppendToSetPipelined(ctx, aggregation.Recordkey, msg)
	p.store.SetExp(ctx, aggregation.Recordkey, p.expireTime.Nanoseconds())
}

func (p *RedisProducer) Stop() error {
	return nil
}
