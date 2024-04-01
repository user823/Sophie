package log

import (
	"github.com/user823/Sophie/pkg/utils"
	"go.uber.org/zap/zapcore"
	"time"
)

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	enc.AppendString(utils.Time2Str(t))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}
