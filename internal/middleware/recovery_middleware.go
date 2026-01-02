package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(recoveryLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Printf("panic recovered: %s", stack)

				recoveryLogger.Error().
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Str("panic_error", fmt.Sprintf("%v", err)).
					Str("stack_at", ExtractFirstAppStackLine(stack)).
					Str("stack", string(stack)).
					Msg("panic_recovered")

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Oops! Something went wrong. Please try again later.",
				})
			}
		}()

		ctx.Next()
	}
}

func ExtractFirstAppStackLine(stack []byte) string {
	lines := strings.Split(string(stack), "\n")

	for i := 0; i < len(lines)-1; i++ {
		if strings.Contains(lines[i], ".go:") &&
			!strings.Contains(lines[i], "/runtime/") &&
			!strings.Contains(lines[i], "/debug/") &&
			!strings.Contains(lines[i], "recovery_middleware.go") &&
			!strings.Contains(lines[i], "gin@") {

			pathLine := strings.TrimSpace(lines[i])
			if nextLine := strings.TrimSpace(lines[i+1]); strings.HasPrefix(nextLine, "\t") {
				pathLine = strings.TrimSpace(nextLine)
			}

			if idx := strings.Index(pathLine, " +0x"); idx > 0 {
				pathLine = pathLine[:idx]
			}

			return pathLine
		}
	}

	return ""
}
