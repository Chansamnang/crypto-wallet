package zaplog

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	ptr *zap.SugaredLogger
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.ptr.Fatal(args...)
}
func (l *ZapLogger) Debug(args ...interface{}) {
	l.ptr.Debug(args...)
}
func (l *ZapLogger) Error(args ...interface{}) {
	l.ptr.Error(args...)
}
func (l *ZapLogger) Info(args ...interface{}) {
	l.ptr.Info(args...)
}
func (l *ZapLogger) Warn(args ...interface{}) {
	l.ptr.Warn(args...)
}
func (l *ZapLogger) Panic(args ...interface{}) {
	l.ptr.Panic(args...)
}
func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.ptr.Panicf(format, args...)
}
func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.ptr.Warnf(format, args...)
}
func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.ptr.Infof(format, args...)
}
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.ptr.Errorf(format, args...)
}
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.ptr.Debugf(format, args...)
}
func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.ptr.Fatalf(format, args...)
}
func (l *ZapLogger) Fatalln(args ...interface{}) {
	l.ptr.Fatal(args...)
}
func (l *ZapLogger) Debugln(args ...interface{}) {
	l.ptr.Debug(args...)
}
func (l *ZapLogger) Errorln(args ...interface{}) {
	l.ptr.Error(args...)
}
func (l *ZapLogger) Infoln(args ...interface{}) {
	l.ptr.Info(args...)
}
func (l *ZapLogger) Warnln(args ...interface{}) {
	l.ptr.Warn(args...)
}
func (l *ZapLogger) Panicln(args ...interface{}) {
	l.ptr.Panic(args...)
}

func (l *ZapLogger) DPanic(v ...interface{}) {
	l.ptr.DPanic(v...)
}
func (l *ZapLogger) DPanicf(format string, v ...interface{}) {
	l.ptr.DPanicf(format, v...)
}

func (l *ZapLogger) WithFields(fields map[string]interface{}) *ZapLogger {
	last := l.ptr.Desugar()
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	new := last.With(zapFields...)
	return &ZapLogger{
		ptr: new.Sugar(),
	}
}
func (l *ZapLogger) WithField(key string, value interface{}) *ZapLogger {
	last := l.ptr.Desugar()
	new := last.With(zap.Any(key, value))
	return &ZapLogger{
		ptr: new.Sugar(),
	}
}

func (l *ZapLogger) WithError(err error) *ZapLogger {
	last := l.ptr.Desugar()
	new := last.With(zap.Error(err))
	return &ZapLogger{
		ptr: new.Sugar(),
	}
}

func NewZapLogger(cfg zap.Config) *ZapLogger {
	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &ZapLogger{
		ptr: logger.Sugar(),
	}
}

func NewFromLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		ptr: logger.WithOptions(zap.AddCallerSkip(1)).Sugar(),
	}
}

func NewDefaultZapLogger() *ZapLogger {
	cfg := zap.NewProductionConfig()
	return NewZapLogger(cfg)
}
