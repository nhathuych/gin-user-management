package middleware

import (
	"context"
	"gin-user-management/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		// add to request context
		contextValue := context.WithValue(ctx.Request.Context(), logger.TraceIdKey, traceID)
		ctx.Request = ctx.Request.WithContext(contextValue)

		// add to gin context
		ctx.Set(logger.TraceIdKey, traceID)

		// add to response header
		ctx.Header("X-Trace-Id", traceID)

		ctx.Next()
	}
}
