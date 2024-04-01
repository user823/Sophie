package consumer

import (
	"context"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/mq"
	"strings"
	"time"
)

const (
	// 每次最多拉取100 条消息
	maxMsgNum = 10
	// 要 >= 20s
	invisibleDuration = 20 * time.Second
)

type RocketMQConsumer struct {
	consumer rmq.SimpleConsumer
}

func NewRocketMQConsumer(endpoints, accessKey, accessSecret string) aggregation.RecordConsumer {
	cfg := &mq.RocketMQConfig{
		Endpoint:                endpoints,
		NameSpace:               aggregation.Recordkey,
		AccessSecret:            accessSecret,
		AccessKey:               accessKey,
		ConsumerGroup:           aggregation.Recordkey,
		MaxAttempts:             3,
		Topics:                  []string{aggregation.Recordkey},
		AwaitDuration:           3 * time.Second,
		SubscriptionExpressions: aggregation.MessageTag,
	}

	consumer, err := mq.NewRocketMQConsumer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &RocketMQConsumer{consumer}
}

func (c *RocketMQConsumer) Connect() bool {
	return c.consumer.Start() == nil
}

func (c *RocketMQConsumer) GetAndDeleteSet(ctx context.Context) []aggregation.LogRecord {
	mvs, err := c.consumer.Receive(ctx, maxMsgNum, invisibleDuration)
	if err != nil {
		if !strings.Contains(err.Error(), "MESSAGE_NOT_FOUND") {
			log.Warnf("rocketmq get message error: %s", err.Error())
		}
		return []aggregation.LogRecord{}
	}

	records := make([]aggregation.LogRecord, len(mvs))
	cnt := 0

	for _, mv := range mvs {
		// 确认消息
		if err = c.consumer.Ack(context.TODO(), mv); err != nil {
			log.Warnf("rocket mq ack message error: %s", err.Error())
			continue
		}

		var tmp map[string]any
		if err = jsoniter.Unmarshal(mv.GetBody(), &tmp); err != nil {
			log.Warnf("json decode failed: %s", err.Error())
			continue
		}
		var record aggregation.LogRecord
		if err = mapstructure.Decode(tmp, &record); err != nil {
			log.Warnf("mapstructure decode failed: %s", err.Error())
			continue
		}

		records[cnt] = record
	}
	records = records[:cnt]
	return records
}

func (c *RocketMQConsumer) Stop() error {
	return c.consumer.GracefulStop()
}
