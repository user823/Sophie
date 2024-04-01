package mq

import (
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/user823/Sophie/pkg/errors"
	"go.uber.org/zap"
	"os"
	"time"
)

type RocketMQConfig struct {
	// 多个proxy格式"ip:port;ip:port"
	Endpoint string
	// 名称空间用来区分话题
	NameSpace     string
	ConsumerGroup string
	AccessKey     string
	AccessSecret  string
	// 投递消息最大尝试次数
	MaxAttempts int
	Topics      []string
	// 消费消息延时
	AwaitDuration time.Duration
	// 消费消息tag
	SubscriptionExpressions string
}

func init() {
	// log to console
	os.Setenv(rmq.ENABLE_CONSOLE_APPENDER, "false")
	os.Setenv(rmq.CLIENT_LOG_LEVEL, zap.ErrorLevel.String())
	os.Setenv(rmq.CLIENT_LOG_ROOT, os.Getenv("HOME"))
	rmq.InitLogger()
}

// 获取rocketmq 生产者
func NewRocketMQProducer(config any) (rmq.Producer, error) {
	cfg, ok := config.(*RocketMQConfig)
	if !ok {
		return nil, errors.New("cannot correctly convert the config type")
	}

	producer, err := rmq.NewProducer(&rmq.Config{
		Endpoint:  cfg.Endpoint,
		NameSpace: cfg.NameSpace,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    cfg.AccessKey,
			AccessSecret: cfg.AccessSecret,
		},
	}, rmq.WithMaxAttempts(int32(cfg.MaxAttempts)), rmq.WithTopics(cfg.Topics...))
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// 获取rocketmq 消费者
func NewRocketMQConsumer(config any) (rmq.SimpleConsumer, error) {
	cfg, ok := config.(*RocketMQConfig)
	if !ok {
		return nil, errors.New("cannot correctly convert the config type")
	}

	filter := make(map[string]*rmq.FilterExpression)
	for i := range cfg.Topics {
		filter[cfg.Topics[i]] = rmq.NewFilterExpression(cfg.SubscriptionExpressions)
	}

	simpleConsumer, err := rmq.NewSimpleConsumer(&rmq.Config{
		Endpoint:      cfg.Endpoint,
		NameSpace:     cfg.NameSpace,
		ConsumerGroup: cfg.ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    cfg.AccessKey,
			AccessSecret: cfg.AccessSecret,
		},
	}, rmq.WithAwaitDuration(cfg.AwaitDuration), rmq.WithSubscriptionExpressions(filter))

	if err != nil {
		return nil, err
	}
	return simpleConsumer, nil
}
