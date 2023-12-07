package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

// 创建zapLogger 配置，用来构建zap.core
// 字段不可访问，只能通过api修改
type Config struct {
	encoderConfig     zapcore.EncoderConfig
	outputWss         []zapcore.WriteSyncer
	errorWss          []zapcore.WriteSyncer
	level             zapcore.Level
	development       bool
	disableCaller     bool
	disableStacktrace bool
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
	}, nil
}

func (c *Config) build() *zap.Logger {
	jsonEncoder := zapcore.NewJSONEncoder(c.encoderConfig)
	output := zapcore.NewMultiWriteSyncer(c.outputWss...)
	errput := zapcore.NewMultiWriteSyncer(c.errorWss...)
	allowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= c.level && level < zap.ErrorLevel
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
		opts = append(opts, zap.AddCallerSkip(1))
	}
	if !c.disableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.PanicLevel))
	}
	l := zap.New(core, opts...)
	return l.Named(c.Name)
}

func (c *Config) SetOutput(w io.Writer) {
	c.CleanOutputWriters()
	c.outputWss = append(c.outputWss, zapcore.AddSync(w))
}

func (c *Config) AddOutput(w io.Writer) {
	c.outputWss = append(c.outputWss, zapcore.AddSync(w))
}

// 全部清除
func (c *Config) CleanOutputWriters() {
	for _, ws := range c.outputWss {
		if closer, ok := ws.(io.Closer); ok {
			closer.Close()
		}
	}
	c.outputWss = []zapcore.WriteSyncer{}
}

// 全部清除
func (c *Config) CleanErrorWriters() {
	for _, ws := range c.errorWss {
		if closer, ok := ws.(io.Closer); ok {
			closer.Close()
		}
	}
	c.errorWss = []zapcore.WriteSyncer{}
}
