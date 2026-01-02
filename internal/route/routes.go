package route

import (
	"gin-user-management/internal/middleware"
	"gin-user-management/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	httpLogger := newLogger("logs/http.log", "info")
	recoveryLogger := newLogger("logs/recovery.log", "info")

	r.Use(
		middleware.RateLimiterMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
	)

	apiV1 := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(apiV1)
	}
}

func newLogger(path, level string) *zerolog.Logger {
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
