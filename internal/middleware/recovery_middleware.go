package middleware

import (
	"fmt"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/logger"
	"net/http"
	"runtime"
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
				stackAt := GetCallerFrame(3)

				if util.IsDevelopment() {
					logger.AppLogger.Error().Msgf("PANIC RECOVERED: %v\n\x1b[31m%s\x1b[0m", err, string(stack))
				}

				recoveryLogger.Error().
					Str("trace_id", logger.GetTraceID(ctx.Request.Context())).
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Interface("panic_error", err).
					Str("stack_at", stackAt).
					// Str("full_stack", string(stack)). // Full stack is optional; usually stored in a field for deep debugging
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

// GetCallerFrame scans the stack to find the first frame originating from the app code.
func GetCallerFrame(skip int) string {
	for i := skip; i < skip+15; i++ { // Increased depth to 15 for complex call chains
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if isAppFile(file) {
			return formatFilePath(file, line)
		}
	}
	return "unknown"
}

// isAppFile identifies if a file path belongs to the application logic.
func isAppFile(file string) bool {
	// Filter out Go internals and external libraries
	if strings.HasPrefix(file, runtime.GOROOT()) ||
		strings.Contains(file, "/pkg/mod/") ||
		strings.Contains(file, "/vendor/") ||
		strings.Contains(file, "recovery_middleware.go") {
		return false
	}
	return true
}

// formatFilePath cleans up absolute paths to be relative to the project root.
func formatFilePath(file string, line int) string {
	if projectName := util.GetProjectName(); projectName != "" {
		if idx := strings.Index(file, projectName); idx != -1 {
			return fmt.Sprintf("%s:%d", file[idx:], line)
		}
	}

	parts := strings.Split(file, "/")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s/%s:%d", parts[len(parts)-2], parts[len(parts)-1], line)
	}

	return fmt.Sprintf("%s:%d", file, line)
}
