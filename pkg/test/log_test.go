package test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/log/aggregation/producer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	redisStore kv.RedisStore
)

func InitRedis() {
	ctx := context.Background()
	cfg := &kv.RedisConfig{
		Addrs:                 []string{"127.0.0.1:6379"},
		Username:              "sophie",
		Password:              "123456",
		Database:              0,
		MasterName:            "",
		MaxIdle:               2000,
		MaxActive:             4000,
		Timeout:               0,
		EnableCluster:         false,
		UseSSL:                false,
		SSLInsecureSkipVerify: false,
	}
	go kv.KeepConnection(ctx, cfg)
	redisStore = kv.NewKVStore("redis", nil).(kv.RedisStore)
	time.Sleep(2 * time.Second)
}

func TestLog(t *testing.T) {
	logger := log.Default()
	logger.Infof("hhh %s", "xxx")
}

func TestZap(t *testing.T) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建 Logger，使用 New 方法并传入 EncoderConfig
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zap.DebugLevel,
	))

	// 输出日志消息
	logger.Sugar().Infof("This is an info message %s", "hhh")
}

func TestAggregation(t *testing.T) {
	pr := producer.NewRedisProducer(redisStore, 1000*time.Second)
	aggregation.NewAnalytics(aggregation.NewAnalyticsOptions(), pr)

	// 开启日志聚合
	analytics := aggregation.GetAnalytics()
	analytics.Start()

	var wg sync.WaitGroup
	wg.Add(100)
	// 多线程测试
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			log.Infof("this is %d", i)
		}(i)
	}
	wg.Wait()

	// 关闭日志聚合刷新缓存
	analytics.Stop()
}

func TestAggregationSub(t *testing.T) {
	InitRedis()

	t.Run("test-aggregation", TestAggregation)
}
