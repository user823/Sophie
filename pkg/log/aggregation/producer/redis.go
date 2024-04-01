package producer

import (
	"context"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

type RedisProducer struct {
	store kv.RedisStore
}

func NewRedisProducer(store kv.RedisStore) aggregation.RecordProducer {
	store.SetKeyPrefix(aggregation.RecordPrefix)
	return &RedisProducer{store}
}

func (p *RedisProducer) Connect() bool {
	return p.store.Connected()
}

func (p *RedisProducer) AppendToSet(ctx context.Context, msg []string) {
	p.store.AppendToSetPipelined(ctx, aggregation.Recordkey, msg)
}

func (p *RedisProducer) Stop() error {
	return nil
}
