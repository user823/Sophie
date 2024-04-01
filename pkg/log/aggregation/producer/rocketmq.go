package producer

import (
	"context"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/mq"
	"github.com/user823/Sophie/pkg/utils"
)

type RocketMQProducer struct {
	producer rmq.Producer
}

func NewRocketMQProducer(endpoints, accessKey, accessSecret string) aggregation.RecordProducer {
	cfg := &mq.RocketMQConfig{
		Endpoint:     endpoints,
		NameSpace:    aggregation.Recordkey,
		AccessSecret: accessSecret,
		AccessKey:    accessKey,
		MaxAttempts:  3,
		Topics:       []string{aggregation.Recordkey},
	}

	producer, err := mq.NewRocketMQProducer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &RocketMQProducer{producer}
}

func (p *RocketMQProducer) Connect() bool {
	return p.producer.Start() == nil
}

func (p *RocketMQProducer) AppendToSet(ctx context.Context, messages []string) {
	for i := range messages {
		msg := &rmq.Message{
			Topic: aggregation.Recordkey,
			Body:  utils.S2b(messages[i]),
		}
		msg.SetTag(aggregation.MessageTag)
		_, err := p.producer.Send(ctx, msg)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

func (p *RocketMQProducer) Stop() error {
	return p.producer.GracefulStop()
}
