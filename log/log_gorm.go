package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoeyZeYi/source/log/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Logger struct {
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	Level                     gormlogger.LogLevel
	SendMsg                   ISendMsg //可用于sql告警信息发送 企业微信/钉钉/飞书等
	SqlLog                    bool     //是否输出sql日志
}

type ISendMsg interface {
	// SendError err/错误信息  errSqlFileAdd/错误sql发生的文件地址
	SendError(ctx context.Context, err error, errSqlFileAdd, sql string)
	// SendTimeout  timeoutSqlFileAdd/超时sql发生的文件地址
	SendTimeout(ctx context.Context, timeoutSqlFileAdd, sql, time string)
}

type IGormLogOption interface {
	apply(*Logger)
}

type gormLogOption func(*Logger)

func (app gormLogOption) apply(option *Logger) {
	app(option)
}

// GormLoggerSlowThreshold sql耗时超过多少写日志
func GormLoggerSlowThreshold(slowThreshold time.Duration) IGormLogOption {
	return gormLogOption(func(option *Logger) {
		option.SlowThreshold = slowThreshold
	})
}

func GormLoggerLevel(level gormlogger.LogLevel) IGormLogOption {
	return gormLogOption(func(option *Logger) {
		option.Level = level
	})
}

func GormSendMsg(sendMsg ISendMsg) IGormLogOption {
	return gormLogOption(func(option *Logger) {
		option.SendMsg = sendMsg
	})
}

func GormSqlLog(sqlLog bool) IGormLogOption {
	return gormLogOption(func(option *Logger) {
		option.SqlLog = true
	})
}

func NewGormLogger(options ...IGormLogOption) *Logger {
	logger := &Logger{}
	for _, option := range options {
		option.apply(logger)
	}
	return logger
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.Level = level
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, params ...interface{}) {
	return
}

func (l *Logger) Warn(ctx context.Context, msg string, params ...interface{}) {
	return
}

func (l *Logger) Error(ctx context.Context, msg string, params ...interface{}) {
	return
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		Error("SQLOutput-Error", zap.Any("code_addr", utils.FileWithLineNum()), zap.String("sql", sql), zap.Error(err), zap.Int64("rows", rows), zap.Any("耗时", elapsed.String()))
		if l.SendMsg != nil {
			l.SendMsg.SendError(ctx, err, utils.FileWithLineNum(), sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		Warn("SQLOutput-Warn", zap.Any("code_addr", utils.FileWithLineNum()), zap.String("sql", sql), zap.Int64("rows", rows), zap.Any("耗时", elapsed.String()), zap.Any("SLOW", slowLog))
		if l.SendMsg != nil {
			l.SendMsg.SendTimeout(ctx, utils.FileWithLineNum(), sql, elapsed.String())
		}
	case l.SqlLog:
		sql, rows := fc()
		Info("SQLOutput-Info", zap.Any("code_addr", utils.FileWithLineNum()), zap.String("sql", sql), zap.Int64("rows", rows), zap.Any("耗时", elapsed.String()))
	}
}
