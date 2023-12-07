package log

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// 定义兼容标准库的logger类型
// 调用格式化方法时，用户自定义类型要实现zap.ObjectMarshaler 接口
type Logger interface {
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Fatalw(msg string, keysAndValues ...any)
	Panic(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
	Panicw(msg string, keysAndValues ...any)
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
	Printw(msg string, keysAndValues ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Infoln(v ...any)
	Infow(msg string, keysAndValues ...any)
	Debug(v ...any)
	Debugf(format string, v ...any)
	Debugln(v ...any)
	Debugw(msg string, keysAndValues ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Warnln(v ...any)
	Warnw(msg string, keysAndValues ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Errorln(v ...any)
	Errorw(msg string, keysAndValues ...any)
	Write(p []byte) (n int, err error)

	// SetAggregation 开启/关闭 日志聚合模式
	// 默认是关闭状态
	SetAggregation()

	// WithValues 添加环境信息
	WithValues(keysAndValues ...interface{})

	// WithName Logger 通过Name进行分层控制，使用'.'分割
	// 每层Name命名使用字母、数字、下划线
	WithName(name string) Logger

	// SetLevel 设置输出级别
	SetLevel(lvl string)
	Name() string
	Flush()
}

type Level = zapcore.Level

// 基于zap实现logger
type zapLogger struct {
	zlogger     *zap.Logger
	aggregation bool
	config      *Config
	buf         *bytes.Buffer
	env         environment
}

// 默认使用的logger
var (
	std *zapLogger
	mu  sync.Mutex
	// 写入到redis中
	writeToRedis = func(msg string) {
		//fmt.Println("write_to_redis", msg)
	}
	_ Logger = &zapLogger{}
)

func Default() Logger {
	if std == nil {
		Init(nil)
	}
	return std
}

func Init(opts *Option) {
	mu.Lock()
	defer mu.Unlock()
	l, err := New(opts)
	if err != nil {
		panic(err)
	}
	std = l.(*zapLogger)
	zap.RedirectStdLog(std.zlogger)
}

func New(opts *Option) (Logger, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	if errs := opts.Validate(); len(errs) != 0 {
		return nil, errs[0]
	}

	if opts.DisableLogger {
		return &noopLogger{}, nil
	}

	config, err := createConfigFromOptions(opts)
	if err != nil {
		return nil, err
	}

	// 为config添加buf缓冲区
	var buf bytes.Buffer
	config.AddOutput(&buf)

	l := config.build()
	return &zapLogger{
		zlogger:     l,
		aggregation: opts.Aggregation,
		config:      config,
		buf:         &buf,
		env:         newEnvironment(),
	}, nil
}

/*
	实现zapLogger
*/

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		return additional
	}
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); i += 2 {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}

		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
	}
	return append(fields, additional...)
}

func (l *zapLogger) Fatal(v ...any) {
	l.zlogger.Sugar().Fatal(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Fatalf(format string, v ...any) {
	l.zlogger.Sugar().Fatalf(format, v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Fatalln(v ...any) {
	l.zlogger.Sugar().Fatalln(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...any) {
	l.zlogger.Sugar().Fatalw(msg, keysAndValues...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Panic(v ...any) {
	l.zlogger.Sugar().Panic(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Panicf(format string, v ...any) {
	l.zlogger.Sugar().Panicf(format, v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Panicln(v ...any) {
	l.zlogger.Sugar().Panicln(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...any) {
	l.zlogger.Sugar().Panicw(msg, keysAndValues...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Info(v ...any) {
	l.zlogger.Sugar().Info(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Infof(format string, v ...any) {
	l.zlogger.Sugar().Infof(format, v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Infoln(v ...any) {
	l.zlogger.Sugar().Infoln(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Infow(msg string, keysAndValues ...any) {
	l.zlogger.Sugar().Infow(msg, keysAndValues...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Print(v ...any) {
	l.Info(v...)
}

func (l *zapLogger) Printf(format string, v ...any) {
	l.Infof(format, v...)
}

func (l *zapLogger) Println(v ...any) {
	l.Infoln(v...)
}

func (l *zapLogger) Printw(msg string, keysAndValues ...any) {
	l.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Debug(v ...any) {
	l.zlogger.Sugar().Debug(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Debugln(v ...any) {
	l.zlogger.Sugar().Debugln(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Debugf(format string, v ...any) {
	l.zlogger.Sugar().Debugf(format, v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...any) {
	l.zlogger.Sugar().Debugw(msg, keysAndValues...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Warn(v ...any) {
	l.zlogger.Sugar().Warnln(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Warnf(format string, v ...any) {
	l.zlogger.Sugar().Warnf(format, v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Warnln(v ...any) {
	l.zlogger.Sugar().Warnln(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...any) {
	l.zlogger.Sugar().Warnw(msg, keysAndValues...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Error(v ...any) {
	l.zlogger.Sugar().Error(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Errorf(format string, v ...any) {
	l.zlogger.Sugar().Errorf(format, v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Errorln(v ...any) {
	l.zlogger.Sugar().Errorln(v...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...any) {
	l.zlogger.Sugar().Errorw(msg, keysAndValues...)
	if l.aggregation {
		writeToRedis(l.buf.String())
	}
}

func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zlogger.Info(string(p))
	return len(p), nil
}

func (l *zapLogger) SetAggregation() {
	l.aggregation = !l.aggregation
}

func (l *zapLogger) WithValues(keysAndValues ...interface{}) {
	l.env.addValues(handleFields(l.zlogger, keysAndValues)...)
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zlogger.Named(name)
	return &zapLogger{
		zlogger:     newLogger,
		aggregation: l.aggregation,
		config:      l.config,
		buf:         &bytes.Buffer{},
	}
}

func (l *zapLogger) SetLevel(lvl string) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(lvl)); err != nil {
		l.zlogger.DPanic("Invalid output level be set", zap.String("setting_level", lvl))
	}
	l.config.level = zapLevel
}

func (l *zapLogger) Name() string {
	return l.zlogger.Name()
}

func (l *zapLogger) Flush() {
	_ = l.zlogger.Sync()
}
