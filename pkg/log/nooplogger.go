package log

type noopLogger struct{}

func (l *noopLogger) Fatal(v ...any)                          {}
func (l *noopLogger) Fatalf(format string, v ...any)          {}
func (l *noopLogger) Fatalln(v ...any)                        {}
func (l *noopLogger) Fatalw(msg string, keysAndValues ...any) {}
func (l *noopLogger) Panic(v ...any)                          {}
func (l *noopLogger) Panicf(format string, v ...any)          {}
func (l *noopLogger) Panicln(v ...any)                        {}
func (l *noopLogger) Panicw(msg string, keysAndValues ...any) {}
func (l *noopLogger) Print(v ...any)                          {}
func (l *noopLogger) Printf(format string, v ...any)          {}
func (l *noopLogger) Println(v ...any)                        {}
func (l *noopLogger) Printw(msg string, keysAndValues ...any) {}
func (l *noopLogger) Info(v ...any)                           {}
func (l *noopLogger) Infof(format string, v ...any)           {}
func (l *noopLogger) Infoln(v ...any)                         {}
func (l *noopLogger) Infow(msg string, keysAndValues ...any)  {}
func (l *noopLogger) Debug(v ...any)                          {}
func (l *noopLogger) Debugf(format string, v ...any)          {}
func (l *noopLogger) Debugln(v ...any)                        {}
func (l *noopLogger) Debugw(msg string, keysAndValues ...any) {}
func (l *noopLogger) Warn(v ...any)                           {}
func (l *noopLogger) Warnf(format string, v ...any)           {}
func (l *noopLogger) Warnln(v ...any)                         {}
func (l *noopLogger) Warnw(msg string, keysAndValues ...any)  {}
func (l *noopLogger) Error(v ...any)                          {}
func (l *noopLogger) Errorf(format string, v ...any)          {}
func (l *noopLogger) Errorln(v ...any)                        {}
func (l *noopLogger) Errorw(msg string, keysAndValues ...any) {}
func (l *noopLogger) Write(p []byte) (n int, err error)       { return 0, nil }
func (l *noopLogger) SetAggregation(bool)                     {}
func (l *noopLogger) SetLevel(lvl string)                     {}
func (l *noopLogger) WithValues(keysAndValues ...string)      {}
func (l *noopLogger) WithName(name string) Logger             { return nil }
func (l *noopLogger) Name() string                            { return "" }
func (l *noopLogger) Flush()                                  {}
