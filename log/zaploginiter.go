package log

import (
	"fmt"
	"github.com/JoeyZeYi/source/log/zap"
	"github.com/JoeyZeYi/source/log/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"time"
)

type zapLogInter interface {
	logInit(*appZapLogConf) (zap.AtomicLevel, *zap.Logger, error)
}

type macZapLogInit struct {
}

func (m *macZapLogInit) logInit(config *appZapLogConf) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapConfig zap.Config
		level     zap.AtomicLevel
		zapLog    *zap.Logger
		err       error
	)
	zapConfig = zap.NewProductionConfig(config.Level)

	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.TimeKey = "timestamp"                   //"@timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //epochSecondTimeEncoder //RFC3339TimeEncoder
	zapConfig.Encoding = "console"
	zapLog, err = zapConfig.Build()
	level = zapConfig.Level
	return level, zapLog, err
}

type winZapLogInit struct {
}

func (w *winZapLogInit) logInit(config *appZapLogConf) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapConfig zap.Config
		level     zap.AtomicLevel
		zapLog    *zap.Logger
		err       error
	)

	zapConfig = zap.NewProductionConfig(config.Level)

	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.TimeKey = "timestamp"                   //"@timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //epochSecondTimeEncoder //RFC3339TimeEncoder
	zapConfig.Encoding = "console"
	zapLog, err = zapConfig.Build()
	level = zapConfig.Level
	return level, zapLog, err
}

type unixLikeZapLogInit struct {
}

func epochMillisTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	nanos := t.UnixNano()
	millis := nanos / int64(time.Millisecond)
	enc.AppendInt64(millis)
}

func epochSecondTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

func epochFullTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (u *unixLikeZapLogInit) logInit(config *appZapLogConf) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		level  zap.AtomicLevel
		zapLog *zap.Logger
	)

	writers := []zapcore.WriteSyncer{os.Stderr}
	output := zapcore.NewMultiWriteSyncer(writers...)
	if len(config.logPath) != 0 {
		output = zapcore.AddSync(&lumberjack.Logger{
			Filename: config.logPath,
			MaxSize:  500, // megabytes
			MaxAge:   5,   // days
		})
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"             //"@timestamp"
	encoderConfig.EncodeTime = epochFullTimeEncoder //epochSecondTimeEncoder //RFC3339TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	level = config.Level

	zapLog = zap.New(zapcore.NewCore(encoder, output, level), zap.AddCaller(), zap.AddStacktrace(zapcore.DPanicLevel))
	return level, zapLog, nil
}

func zapLogInit(config *appZapLogConf) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapInit zapLogInter
		level   zap.AtomicLevel
		logger  *zap.Logger
		err     error
	)

	if runtime.GOOS == "darwin" {
		zapInit = &macZapLogInit{}
	} else if runtime.GOOS == "windows" {
		zapInit = &winZapLogInit{}
	} else {
		zapInit = &unixLikeZapLogInit{}
	}

	if level, logger, err = zapInit.logInit(config); err != nil {
		fmt.Printf("loginit err:%v", err)
		return level, logger, err
	}

	if config.withPid {
		logger = logger.With(zap.Int("pid", os.Getpid()))
	}

	if config.HostName != "" {
		logger = logger.With(zap.String("hostname", config.HostName))
	}

	if config.ElkTemplateName != "" {
		if runtime.GOOS == "windows" {
			logger = logger.With(zap.String("service", "windows"))
		} else {
			logger = logger.With(zap.String("service", config.ElkTemplateName))
		}

	}
	return level, logger, nil
}
