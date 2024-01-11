package config

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewLogger Creates new logger instance
func NewLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	slogger := logger.Sugar()

	return slogger
}

// LoggerModule
var LoggerModule = fx.Options(
	fx.Provide(NewLogger),
)
