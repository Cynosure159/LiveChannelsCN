package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init 初始化 Logger
// mode: "dev" (开发模式，控制台高亮) | "prod" (生产模式，JSON)
func Init(mode string) {
	once.Do(func() {
		var config zap.Config
		if mode == "prod" {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// 自定义时间格式
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		var err error
		log, err = config.Build()
		if err != nil {
			panic(err)
		}
	})
}

// Get 返回全局 Logger
func Get() *zap.Logger {
	if log == nil {
		Init("dev") // 默认初始化
	}
	return log
}

// Sync 刷新缓冲
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Info 快捷方式
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error 快捷方式
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal 快捷方式
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Warn 快捷方式
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Debug 快捷方式
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}
