package log

import (
	"fmt"
	"go.uber.org/zap"
)

type environment []string

type ezapLogger struct {
	logger Logger
	env    []any
}

func (l *ezapLogger) Fatal(v ...any) {
	l.logger.Fatalw(fmt.Sprint(v...), l.env...)
}

func (l *ezapLogger) Fatalf(format string, v ...any) {
	l.logger.Fatalw(fmt.Sprintf(format, v...), l.env...)
}

func (l *ezapLogger) Fatalln(v ...any) {
	l.logger.Fatalw(fmt.Sprintln(v...), l.env...)
}

func (l *ezapLogger) Fatalw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, l.env...)
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *ezapLogger) Panic(v ...any) {
	l.logger.Panicw(fmt.Sprint(v...), l.env...)
}

func (l *ezapLogger) Panicf(format string, v ...any) {
	l.logger.Panicw(fmt.Sprintf(format, v...), l.env...)
}

func (l *ezapLogger) Panicln(v ...any) {
	l.logger.Panicw(fmt.Sprintln(v...), l.env...)
}

func (l *ezapLogger) Panicw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, l.env...)
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *ezapLogger) Print(v ...any) {
	l.logger.Printw(fmt.Sprint(v...), l.env...)
}

func (l *ezapLogger) Printf(format string, v ...any) {
	l.logger.Printw(fmt.Sprintf(format, v...), l.env...)
}

func (l *ezapLogger) Println(v ...any) {
	l.logger.Printw(fmt.Sprintln(v...), l.env...)
}

func (l *ezapLogger) Printw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, l.env...)
	l.logger.Printw(msg, keysAndValues...)
}

func (l *ezapLogger) Info(v ...any) {
	l.logger.Infow(fmt.Sprint(v...), l.env...)
}

func (l *ezapLogger) Infof(format string, v ...any) {
	l.logger.Infow(fmt.Sprintf(format, v...), l.env...)
}

func (l *ezapLogger) Infoln(v ...any) {
	l.logger.Infow(fmt.Sprintln(v...), l.env...)
}

func (l *ezapLogger) Infow(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, l.env...)
	l.logger.Infow(msg, keysAndValues...)
}

func (l *ezapLogger) Debug(v ...any) {
	l.logger.Debugw(fmt.Sprint(v...), l.env...)
}

func (l *ezapLogger) Debugf(format string, v ...any) {
	l.logger.Debugw(fmt.Sprintf(format, v...), l.env...)
}

func (l *ezapLogger) Debugln(v ...any) {
	l.logger.Debugw(fmt.Sprintln(v...), l.env...)
}

func (l *ezapLogger) Debugw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, l.env...)
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *ezapLogger) Warn(v ...any) {
	l.logger.Warnw(fmt.Sprint(v...), l.env...)
}

func (l *ezapLogger) Warnf(format string, v ...any) {
	l.logger.Warnw(fmt.Sprintf(format, v...), l.env...)
}

func (l *ezapLogger) Warnln(v ...any) {
	l.logger.Warnw(fmt.Sprintln(v...), l.env...)
}

func (l *ezapLogger) Warnw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, l.env...)
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *ezapLogger) Error(v ...any) {
	l.logger.Errorw(fmt.Sprint(v...), l.env...)
}

func (l *ezapLogger) Errorf(format string, v ...any) {
	l.logger.Errorw(fmt.Sprintf(format, v...), l.env...)
}

func (l *ezapLogger) Errorln(v ...any) {
	l.logger.Errorw(fmt.Sprintln(v...), l.env...)
}

func (l *ezapLogger) Errorw(msg string, keysAndValues ...any) {
	keysAndValues = append(keysAndValues, l.env...)
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *ezapLogger) Write(p []byte) (n int, err error) {
	return l.logger.Write(p)
}

func (l *ezapLogger) WithValues(keysAndValues ...any) Logger {
	l.env = keysAndValues
	return l
}

func (l *ezapLogger) WithName(name string) Logger {
	l.logger = l.logger.WithName(name)
	return l
}

func (l *ezapLogger) SetLevel(lvl Level) {
	l.logger.SetLevel(lvl)
}

func (l *ezapLogger) Name() string {
	return l.logger.Name()
}

func (l *ezapLogger) Flush() {
	l.logger.Flush()
}

func (l *ezapLogger) RedirectToStd() {
	if a, ok := l.logger.(StdLoggerAdapter); ok {
		a.RedirectToStd()
	}
}

func (l *ezapLogger) ZapLogger() *zap.Logger {
	return l.logger.ZapLogger()
}
