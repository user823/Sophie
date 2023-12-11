package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Config struct {
	encoderConfig     zapcore.EncoderConfig
	outputWss         []zapcore.WriteSyncer
	errorWss          []zapcore.WriteSyncer
	level             zapcore.Level
	development       bool
	disableCaller     bool
	disableStacktrace bool
	aggregation       bool
	skipCaller        int
	// 初始logger的名字，一般为根模块的名字
	Name string
}

func createConfigFromOptions(opts *Option) (*Config, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		return nil, err
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var outputWss []zapcore.WriteSyncer
	for _, path := range opts.OutputPaths {
		if path == stdout {
			outputWss = append(outputWss, os.Stdout)
			continue
		}
		f, err := os.Create(path)
		if err != nil {
			continue
		}
		outputWss = append(outputWss, f)
	}
	if len(outputWss) == 0 {
		outputWss = append(outputWss, os.Stdout)
	}

	var errorWss []zapcore.WriteSyncer
	for _, path := range opts.ErrorOutputPaths {
		if path == stderr {
			errorWss = append(errorWss, os.Stderr)
			continue
		}
		f, err := os.Create(path)
		if err != nil {
			continue
		}
		errorWss = append(errorWss, f)
	}
	if len(errorWss) == 0 {
		errorWss = append(errorWss, os.Stderr)
	}

	return &Config{
		Name:              opts.Name,
		encoderConfig:     encoderConfig,
		outputWss:         outputWss,
		errorWss:          errorWss,
		level:             zapLevel,
		development:       opts.Development,
		disableCaller:     opts.DisableCaller,
		disableStacktrace: opts.DisableStacktrace,
		skipCaller:        opts.SkipCaller,
	}, nil
}

func (c *Config) build() *zapLogger {
	L := &zapLogger{env: newEnvironment(), level: c.level}

	// 添加buf通道
	c.AddOutput(zapcore.AddSync(&logbuf))

	jsonEncoder := zapcore.NewJSONEncoder(c.encoderConfig)
	output := zapcore.NewMultiWriteSyncer(c.outputWss...)
	errput := zapcore.NewMultiWriteSyncer(c.errorWss...)
	allowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= L.level && level < zap.ErrorLevel
	})
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, output, allowPriority),
		zapcore.NewCore(jsonEncoder, errput, highPriority),
	)

	var opts []zap.Option
	if !c.disableCaller {
		opts = append(opts, zap.AddCaller())
		opts = append(opts, zap.AddCallerSkip(c.skipCaller))
	}
	if !c.disableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.PanicLevel))
	}

	logger := zap.New(core, opts...).Named(c.Name)
	L.zlogger = logger

	// 设置启用日志聚合
	L.SetAggregation(c.aggregation)
	return L
}

func (c *Config) SetOutput(w io.Writer) error {
	if err := c.CleanOutputWriters(); err != nil {
		return err
	}
	c.outputWss = append(c.outputWss, zapcore.AddSync(w))
	return nil
}

func (c *Config) AddOutput(w io.Writer) {
	c.outputWss = append(c.outputWss, zapcore.AddSync(w))
}

// 全部清除
func (c *Config) CleanOutputWriters() error {
	for _, ws := range c.outputWss {
		if closer, ok := ws.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	c.outputWss = []zapcore.WriteSyncer{}
	return nil
}

// 全部清除
func (c *Config) CleanErrorWriters() error {
	for _, ws := range c.errorWss {
		if closer, ok := ws.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	c.errorWss = []zapcore.WriteSyncer{}
	return nil
}
