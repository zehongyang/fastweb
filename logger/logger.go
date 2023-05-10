package logger

import (
	"github.com/zehongyang/fastweb/config"
	"go.uber.org/zap"
)

func D(fields ...zap.Field) {
	config.GLogger.Debug("", fields...)
}

func I(fields ...zap.Field) {
	config.GLogger.Info("", fields...)
}

func W(fields ...zap.Field) {
	config.GLogger.Warn("", fields...)
}

func E(fields ...zap.Field) {
	config.GLogger.Error("", fields...)
}

func F(fields ...zap.Field) {
	config.GLogger.Fatal("", fields...)
}
