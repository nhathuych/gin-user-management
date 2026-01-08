package logger

import (
	"context"
	"io"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type contextKey string

const TraceIdKey contextKey = "trace_id"

func NewLogger(config LoggerConfig) *zerolog.Logger {
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var writer io.Writer
	writer = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxAge,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIdKey).(string); ok {
		return traceID
	}
	return ""
}
