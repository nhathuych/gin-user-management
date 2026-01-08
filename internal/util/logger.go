package util

import (
	"gin-user-management/pkg/logger"

	"github.com/rs/zerolog"
)

func NewLogger(path, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:      level,
		Filename:   path,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
	}
	return logger.NewLogger(config)
}
