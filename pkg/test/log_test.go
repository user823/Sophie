package test

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/db/kv/redis"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

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
	ctx := context.Background()
	go redis.KeepConnection(ctx, &redis.RedisConfig{
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "123456",
		Database: 0,
	})
	redisClient := kv.NewKVStore("redis").(kv.RedisStore)
	aggregation.NewAnalytics(aggregation.NewAnalyticsOptions(), redisClient, log.GetRecordMagager())

	// 开启日志聚合
	analytics := aggregation.GetAnalytics()
	analytics.Start()
	log.SetAggregation(true)

	// 附带环境信息
	log.Infow("my first log~", "a", "b")

	// 关闭日志聚合刷新缓存
	analytics.Stop()
	result := redisClient.GetAndDeleteSet(ctx, aggregation.RecordkeyName)
	fmt.Println(result)
}
