package zlogger

import (
	"fmt"
	"wallet/pkg/common/config"
	"wallet/pkg/safego"
	"wallet/pkg/zlogger/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zLog *zap.Logger

	// 逻辑相关 内部日志接口
	zLogInner *zap.Logger
	sLogInner *zap.SugaredLogger

	// 外部追加属性
	appendOptions []zap.Option
)

type logConfig struct {
	// 日志等级. debug/info/warn/error
	Level string
	// 日志追踪等级
	StackTrace string
	// 输出
	Output string
	// 是否添加caller
	Caller bool
	// 配置读取前缀
	prefix string
}

func newLogConfig() *logConfig {
	cfg := &logConfig{
		Level:      config.Config.Log.LogLevel,
		StackTrace: "panic",
		Output:     "stdout",
		Caller:     false,
	}
	return cfg
}

var logFileName string
var alertCore zapcore.Core
var globalFields map[string]string

func InitLogConfig(logName string) {
	// 导入落地功能
	lumberjack.InitLumberjackLogger()
	// 默认配置设置
	logFileName = logName
	// 日志配置
	var logicConfig = newLogConfig()
	logicConfig.Output = logFileName
	zLog, _ = zap.NewDevelopment()
	sLogInner = zLog.WithOptions(zap.AddCallerSkip(1)).Sugar()
	zLogInner = zLog.WithOptions(zap.AddCallerSkip(1))
	err := initLog(logicConfig)
	if err != nil {
		return
	}
}

func getLogLevel(lvl string) zap.AtomicLevel {
	lv := zap.NewAtomicLevel()
	switch lvl {
	case "panic":
		lv.SetLevel(zap.PanicLevel)
	case "fatal":
		lv.SetLevel(zap.FatalLevel)
	case "error":
		lv.SetLevel(zap.ErrorLevel)
	case "info":
		lv.SetLevel(zap.InfoLevel)
	case "debug", "trace":
		lv.SetLevel(zap.DebugLevel)
	case "warn":
		lv.SetLevel(zap.WarnLevel)
	}
	return lv
}

func SetEmptyLogger() {
	zLogInner = zap.NewNop()
	sLogInner = zLogInner.Sugar()
	zLog = zLogInner
}

// 新建日志接口
func newLogger(dev bool, logCfg *logConfig) (*zap.Logger, error) {
	var cfg zap.Config
	// if dev {
	// 	cfg = zap.NewDevelopmentConfig()
	// } else {
	cfg = zap.NewProductionConfig()
	//}
	cfg.Level = getLogLevel(logCfg.Level)
	cfg.DisableStacktrace = true
	var opts []zap.Option
	if logCfg.StackTrace != "" {
		opts = append(opts, zap.AddStacktrace(getLogLevel(logCfg.StackTrace).Level()))
	}
	if logCfg.Caller {
		opts = append(opts, zap.AddCaller())
	}
	if logFileName != "" {
		logCfg.Output = logFileName
	}
	cfg.OutputPaths = []string{logCfg.Output, "stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05]")
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	if alertCore != nil {
		opts = append(opts, zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, alertCore)
		}))
	}
	// 追加属性
	if len(appendOptions) > 0 {
		opts = append(opts, appendOptions...)
	}

	return cfg.Build(opts...)
}

// 从配置文件加载Log配置，并初始化Log Hook
func initLog(logicLogCfg *logConfig) error {
	fmt.Println("init logger")
	var dev bool

	// 生成日志接口
	logic, err := newLogger(dev, logicLogCfg)
	if err != nil {
		return err
	}

	// 本地日志接口
	zLog = logic
	sLogInner = logic.WithOptions(zap.AddCallerSkip(1)).Sugar()
	zLogInner = logic.WithOptions(zap.AddCallerSkip(1))

	safego.PanicCatchFunc = func(name string, p interface{}) {
		zLog.Error("receive panic", zap.String("name", name), zap.Any("p", p))
	}

	if globalFields != nil {
		setGlobalFields(globalFields)
	}

	return nil
}

func SetGlobalFields(fields map[string]string) {
	globalFields = fields
}

func setGlobalFields(fields map[string]string) {
	nf := make(map[string]interface{})
	for k, v := range fields {
		nf[k] = v
	}
	zf := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zf = append(zf, zap.String(k, v))
	}

	zLog = zLog.With(zf...)
	sLogInner = zLog.WithOptions(zap.AddCallerSkip(1)).Sugar()
	zLogInner = zLog.WithOptions(zap.AddCallerSkip(1))
}

func Fatal(format ...interface{}) {
	sLogInner.Fatal(format...)
}
func Debug(args ...interface{}) {
	sLogInner.Debug(args...)
}
func Error(args ...interface{}) {
	sLogInner.Error(args...)
}
func Info(args ...interface{}) {
	sLogInner.Info(args...)
}
func Warn(args ...interface{}) {
	sLogInner.Warn(args...)
}
func Panic(args ...interface{}) {
	sLogInner.Panic(args...)
}
func Warnf(format string, args ...interface{}) {
	sLogInner.Warnf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	sLogInner.Panicf(format, args...)
}
func Infof(format string, args ...interface{}) {
	sLogInner.Infof(format, args...)
}
func Errorf(format string, args ...interface{}) {
	sLogInner.Errorf(format, args...)
}
func Debugf(format string, args ...interface{}) {
	sLogInner.Debugf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	sLogInner.Fatalf(format, args...)
}
func Fatalln(args ...interface{}) {
	sLogInner.Fatal(args...)
}
func Debugln(args ...interface{}) {
	sLogInner.Debug(args...)
}
func Errorln(args ...interface{}) {
	sLogInner.Error(args...)
}
func Infoln(args ...interface{}) {
	sLogInner.Info(args...)
}
func Warnln(args ...interface{}) {
	sLogInner.Warn(args...)
}
func Panicln(args ...interface{}) {
	sLogInner.Panic(args...)
}

func Fatalw(msg string, fields ...zap.Field) {
	zLogInner.Fatal(msg, fields...)
}
func Debugw(msg string, fields ...zap.Field) {
	zLogInner.Debug(msg, fields...)
}
func Errorw(msg string, fields ...zap.Field) {
	zLogInner.Error(msg, fields...)
}
func Infow(msg string, fields ...zap.Field) {
	zLogInner.Info(msg, fields...)
}
func Warnw(msg string, fields ...zap.Field) {
	zLogInner.Warn(msg, fields...)
}
func Panicw(msg string, fields ...zap.Field) {
	zLogInner.Panic(msg, fields...)
}
