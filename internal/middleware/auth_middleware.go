package middleware

import (
	"gin-user-management/internal/util"
	"gin-user-management/pkg/auth"
	"gin-user-management/pkg/cache"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtGenerator auth.TokenGenerator, redisCache cache.RedisCacheService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			util.ResponseError(ctx, util.NewError("Unauthorized.", util.ErrCodeUnauthorized))
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtGenerator.ParseWithClaims(tokenString)
		if err != nil {
			util.ResponseError(ctx, util.WrapError(err, "Invalid accesstoken token.", util.ErrCodeUnauthorized))
			ctx.Abort()
			return
		}

		jti := claims.ID // jti
		key := cache.BlacklistAccessTokenKey(jti)
		exists, err := redisCache.Exists(ctx, key)
		if exists && err == nil {
			util.ResponseError(ctx, util.NewError("Token has been revoked.", util.ErrCodeUnauthorized))
			ctx.Abort()
			return
		}

		ctx.Set("user_uuid", claims.UUID)
		ctx.Set("user_role", claims.Role)

		// uuid, _ := ctx.Get("user_uuid")
		// role, _ := ctx.Get("user_role")
		// logger.AppLogger.Info().Msgf(`authenticated user: {uuid: %s, role: %d}`, uuid, role)

		ctx.Next()
	}
}
