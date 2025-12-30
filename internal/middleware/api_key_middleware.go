// ApiKeyMiddleware is used to protect internal or service-to-service APIs.
// It validates requests using a static API key provided in the `X-API-Key` header.
//
// This middleware is NOT intended for user authentication (JWT).
// It is typically used for:
// - Internal APIs
// - Microservice-to-microservice communication
// - Cron jobs
// - Webhooks or trusted external services
//
// Requests without a valid API key will be rejected before reaching the handler.

package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	expectedKey := os.Getenv("API_KEY")
	if expectedKey == "" {
		expectedKey = "your-api-key"
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("x-api-key")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing API key",
			})
			return
		}

		if apiKey != expectedKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			return
		}

		ctx.Set("username", "huy")

		ctx.Next()
	}
}
