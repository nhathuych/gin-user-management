package middleware

import (
	"gin-user-management/internal/util"
	"gin-user-management/pkg/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtGenerator auth.TokenGenerator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			util.ResponseError(ctx, util.NewError("Invalid or missing Authorization header", util.ErrCodeUnauthorized))
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtGenerator.ParseWithClaims(tokenString)
		if err != nil {
			util.ResponseError(ctx, util.NewError("Invalid or missing Authorization header", util.ErrCodeUnauthorized))
			ctx.Abort()
			return
		}

		ctx.Set("user_uuid", claims.UUID)
		ctx.Set("user_role", claims.Role)

		// uuid, _ := ctx.Get("user_uuid")
		// role, _ := ctx.Get("user_role")
		// log.Printf(`authenticated user: {uuid: %s, role: %d}`, uuid, role)

		ctx.Next()
	}
}
