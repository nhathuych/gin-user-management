package v1

import (
	"gin-user-management/internal/repository"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/auth"
	"gin-user-management/pkg/cache"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo           repository.UserRepository
	redisCache     cache.RedisCacheService
	tokenGenerator auth.TokenGenerator
}

func NewAuthService(repo repository.UserRepository, redisCache cache.RedisCacheService, tokenGenerator auth.TokenGenerator) AuthService {
	return &authService{
		repo:           repo,
		redisCache:     redisCache,
		tokenGenerator: tokenGenerator,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()

	email = util.NormalizeString(email)
	user, err := as.repo.GetByEmail(context, email)
	if err != nil {
		return "", "", 0, util.NewError("Invalid credentials.", util.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", 0, util.NewError("Invalid credentials.", util.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to generate access token.", util.ErrCodeInternal)
	}

	refreshToken, err := as.tokenGenerator.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to generate refresh token.", util.ErrCodeInternal)
	}

	if err := as.tokenGenerator.StoreRefreshToken(context, refreshToken); err != nil {
		return "", "", 0, util.WrapError(err, "Failed to store refresh token.", util.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context, refreshToken string) error {
	context := ctx.Request.Context()

	// 1. Revoke refresh token (core)
	if err := as.tokenGenerator.RevokeRefreshToken(context, refreshToken); err != nil {
		return util.WrapError(err, "Failed to revoke refresh token.", util.ErrCodeInternal)
	}

	// 2. Optional: blacklist the access token if present
	authHeader := ctx.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, err := as.tokenGenerator.ParseWithClaims(accessToken); err == nil {
			jti := claims.ID             // jti
			exp := claims.ExpiresAt.Time // exp
			ttl := time.Until(exp)

			if ttl > 0 {
				key := cache.BlacklistAccessTokenKey(jti)
				_ = as.redisCache.Set(context, key, "revoked", ttl)
			}
		}
	}

	return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, oldRefreshToken string) (string, string, int, error) {
	context := ctx.Request.Context()

	refreshTokenInfo, err := as.tokenGenerator.ValidateRefreshToken(context, oldRefreshToken)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to validate refresh token.", util.ErrCodeUnauthorized)
	}

	userUUID, err := uuid.Parse(refreshTokenInfo.UUID)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to parse user UUID from refresh token.", util.ErrCodeInternal)
	}
	user, err := as.repo.GetByUUID(context, userUUID)
	if err != nil {
		return "", "", 0, util.NewError("User not found.", util.ErrCodeUnauthorized)
	}

	newAccessToken, err := as.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to generate access token.", util.ErrCodeInternal)
	}

	newRefreshToken, err := as.tokenGenerator.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, util.WrapError(err, "Failed to generate refresh token.", util.ErrCodeInternal)
	}

	if err := as.tokenGenerator.RevokeRefreshToken(context, oldRefreshToken); err != nil {
		return "", "", 0, util.WrapError(err, "Failed to revoke refresh token.", util.ErrCodeInternal)
	}

	if err := as.tokenGenerator.StoreRefreshToken(context, newRefreshToken); err != nil {
		return "", "", 0, util.WrapError(err, "Failed to store refresh token.", util.ErrCodeInternal)
	}

	return newAccessToken, newRefreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
