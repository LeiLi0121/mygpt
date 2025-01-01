package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger // 全局日志对象

// InitLogger 初始化日志模块
func InitLogger(debug bool) {
	var cfg zap.Config

	if debug {
		// 开发环境：输出到控制台，带有详细的调用信息
		cfg = zap.NewDevelopmentConfig()
	} else {
		// 生产环境：JSON 格式，输出到文件或日志收集系统
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp" // 自定义时间字段名
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic("无法初始化日志模块: " + err.Error())
	}
}

// CloseLogger 优雅关闭日志
func CloseLogger() {
	_ = Log.Sync()
}
