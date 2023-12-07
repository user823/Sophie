package test

import (
	"github.com/user823/Sophie/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	log.Init(nil)
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
