package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type ZapGormLogger struct {
	zapLogger *zap.Logger
}

func NewZapGormLogger(zapLogger *zap.Logger) logger.Interface {
	return &ZapGormLogger{zapLogger: zapLogger}
}

func (l *ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.zapLogger.Sugar().Infof(msg, data...)
}

func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.zapLogger.Sugar().Warnf(msg, data...)
}

func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.zapLogger.Sugar().Errorf(msg, data...)
}

func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.zapLogger.Error("SQL 执行失败",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			zap.Error(err),
		)
	} else {
		l.zapLogger.Info("SQL 执行成功",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	}
}
