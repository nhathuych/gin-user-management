package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

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
	IsDev      bool
}

type contextKey string

const TraceIdKey contextKey = "trace_id"

var AppLogger *zerolog.Logger

func InitLogger(config LoggerConfig) {
	if AppLogger != nil {
		return
	}

	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	AppLogger = NewLogger(config)
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	writers := []io.Writer{
		&lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		},
	}

	if config.IsDev {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatMessage: func(i interface{}) string {
				if i == nil {
					return ""
				}
				return fmt.Sprintf("%s", i)
			},
		})
	}

	logger := zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()
	return &logger
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIdKey).(string); ok {
		return traceID
	}
	return ""
}
