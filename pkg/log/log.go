package log

import (
	"go.uber.org/zap"
	"sync"
)

// 定义兼容标准库的logger类型
// 调用格式化方法时，用户自定义类型要实现zap.ObjectMarshaler 接口以完成正确编码
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
	ZapLogger() *zap.Logger

	// WithValues 添加环境信息
	WithValues(keysAndValues ...any) Logger

	// WithName Logger 通过Name进行分层控制，使用'.'分割
	// 每层Name命名使用字母、数字、下划线
	WithName(name string) Logger

	// SetLevel 设置输出级别
	SetLevel(lvl Level)
	Name() string
	Flush()
}

// 适配标准log包
type StdLoggerAdapter interface {
	RedirectToStd()
}

// 基于zap实现logger
type zapLogger struct {
	zlogger *zap.Logger
	level   Level
	flags   uint32
}

const (
	aggregationFlag uint32 = 1 << iota
)

// 默认使用的logger
var (
	std Logger
	mu  sync.Mutex
)

// log 包导入即可用，使用默认的配置
func init() {
	// 初始化stdLogger
	opts := DefaultOptions()
	l, err := New(opts)
	if err != nil {
		panic(err)
	}
	std = l

	if a, ok := std.(StdLoggerAdapter); ok {
		a.RedirectToStd()
	}

	zap.ReplaceGlobals(ZapLogger())
}

func Default() Logger {
	return std
}

func SetGlobal(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	std = l
	if a, ok := std.(StdLoggerAdapter); ok {
		a.RedirectToStd()
	}

	zap.ReplaceGlobals(ZapLogger())
}

func New(opts *Options) (Logger, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	if errs := opts.Validate(); len(errs) != 0 {
		return nil, errs[0]
	}

	if opts.DisableLogger {
		return &noopLogger{}, nil
	}

	config, err := CreateConfigFromOptions(opts)
	if err != nil {
		return nil, err
	}

	return config.Build(), nil
}

/*
	实现zapLogger
*/

func (l *zapLogger) Fatal(v ...any) {
	l.zlogger.Sugar().Fatal(v...)
}

func (l *zapLogger) Fatalf(format string, v ...any) {

	l.zlogger.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalln(v ...any) {

	l.zlogger.Sugar().Fatalln(v...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...any) {

	l.zlogger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Panic(v ...any) {

	l.zlogger.Sugar().Panic(v...)
}

func (l *zapLogger) Panicf(format string, v ...any) {

	l.zlogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panicln(v ...any) {

	l.zlogger.Sugar().Panicln(v...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...any) {

	l.zlogger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Info(v ...any) {

	l.zlogger.Sugar().Info(v...)
}

func (l *zapLogger) Infof(format string, v ...any) {

	l.zlogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infoln(v ...any) {

	l.zlogger.Sugar().Infoln(v...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...any) {

	l.zlogger.Sugar().Infow(msg, keysAndValues...)
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
}

func (l *zapLogger) Debugln(v ...any) {

	l.zlogger.Sugar().Debugln(v...)
}

func (l *zapLogger) Debugf(format string, v ...any) {

	l.zlogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...any) {

	l.zlogger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Warn(v ...any) {

	l.zlogger.Sugar().Warnln(v...)
}

func (l *zapLogger) Warnf(format string, v ...any) {

	l.zlogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnln(v ...any) {

	l.zlogger.Sugar().Warnln(v...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...any) {

	l.zlogger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(v ...any) {

	l.zlogger.Sugar().Error(v...)
}

func (l *zapLogger) Errorf(format string, v ...any) {

	l.zlogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorln(v ...any) {

	l.zlogger.Sugar().Errorln(v...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...any) {

	l.zlogger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Write(p []byte) (n int, err error) {

	l.zlogger.Info(string(p))
	return len(p), nil
}

func (l *zapLogger) WithValues(keysAndValues ...any) Logger {
	if len(keysAndValues)%2 != 0 {
		std.Error("WithValues add odd parameters into environment")
		return nil
	}

	return &ezapLogger{
		logger: l,
		env:    keysAndValues,
	}
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zlogger.Named(name)
	return &zapLogger{
		zlogger: newLogger,
		flags:   l.flags,
	}
}

func (l *zapLogger) SetLevel(lvl Level) {
	l.level = lvl
}

func (l *zapLogger) Name() string {
	return l.zlogger.Name()
}

func (l *zapLogger) Flush() {
	_ = l.zlogger.Sync()
}

func (l *zapLogger) RedirectToStd() {
	zap.RedirectStdLog(l.zlogger)
}

func (l *zapLogger) ZapLogger() *zap.Logger {
	return l.zlogger
}

func Fatal(v ...any) {
	std.Fatal(v...)
}

func Fatalf(format string, v ...any) {
	std.Fatalf(format, v...)
}

func Fatalln(v ...any) {
	std.Fatalln(v...)
}

func Fatalw(msg string, keysAndValues ...any) {
	std.Fatalw(msg, keysAndValues...)
}

func Panic(v ...any) {
	std.Panic(v...)
}

func Panicf(format string, v ...any) {
	std.Panicf(format, v...)
}

func Panicln(v ...any) {
	std.Panicln(v...)
}

func Panicw(msg string, keysAndValues ...any) {
	std.Panicw(msg, keysAndValues...)
}

func Info(v ...any) {
	std.Info(v...)
}

func Infof(format string, v ...any) {
	std.Infof(format, v...)
}

func Infoln(v ...any) {
	std.Infoln(v...)
}

func Infow(msg string, keysAndValues ...any) {
	std.Infow(msg, keysAndValues...)
}

func Print(v ...any) {
	std.Info(v...)
}

func Printf(format string, v ...any) {
	std.Infof(format, v...)
}

func Println(v ...any) {
	std.Infoln(v...)
}

func Printw(msg string, keysAndValues ...any) {
	std.Infow(msg, keysAndValues...)
}

func Debug(v ...any) {
	std.Debug(v...)
}

func Debugln(v ...any) {
	std.Debugln(v...)
}

func Debugf(format string, v ...any) {
	std.Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...any) {
	std.Debugw(msg, keysAndValues...)
}

func Warn(v ...any) {
	std.Warnln(v...)
}

func Warnf(format string, v ...any) {
	std.Warnf(format, v...)
}

func Warnln(v ...any) {
	std.Warnln(v...)
}

func Warnw(msg string, keysAndValues ...any) {
	std.Warnw(msg, keysAndValues...)
}

func Error(v ...any) {
	std.Error(v...)
}

func Errorf(format string, v ...any) {
	std.Errorf(format, v...)
}

func Errorln(v ...any) {
	std.Errorln(v...)
}

func Errorw(msg string, keysAndValues ...any) {
	std.Errorw(msg, keysAndValues...)
}

func WithValues(keysAndValues ...any) Logger {
	return std.WithValues(keysAndValues...)
}

func Flush() {
	std.Flush()
}

func ZapLogger() *zap.Logger {
	return std.ZapLogger()
}
