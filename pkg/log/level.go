package log

import "go.uber.org/zap/zapcore"

type Level = zapcore.Level

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)
