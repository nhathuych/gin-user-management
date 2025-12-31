package route

import (
	"gin-user-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	r.Use(
		middleware.LoggerMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
		middleware.RateLimiterMiddleware(),
	)

	apiV1 := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(apiV1)
	}
}
