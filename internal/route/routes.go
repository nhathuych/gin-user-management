package route

import (
	"gin-user-management/internal/middleware"
	"gin-user-management/internal/util"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	httpLogger := util.NewLogger("logs/http.log", "info")
	recoveryLogger := util.NewLogger("logs/recovery.log", "info")

	r.Use(
		middleware.RateLimiterMiddleware(),
		middleware.TraceMiddleware(),
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
