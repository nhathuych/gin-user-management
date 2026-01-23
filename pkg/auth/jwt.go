package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/cache"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JWTGenerator struct {
	redis cache.RedisCacheService
}

type CustomClaims struct {
	UUID string `json:"sub"`
	Role int32  `json:"role"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UUID      string    `json:"sub"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

var (
	jwtSecret = []byte(util.GetEnv("JWT_SECRET", "gin-user-management-default-secret"))
)

func NewJWTGenerator(redis cache.RedisCacheService) TokenGenerator {
	return &JWTGenerator{
		redis: redis,
	}
}

/**************** ACCESS TOKEN ****************/

func (jg *JWTGenerator) GenerateAccessToken(user sqlc.User) (string, error) {
	claims := CustomClaims{
		UUID: user.Uuid.String(),
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // jti
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (jg *JWTGenerator) ParseWithClaims(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, util.NewError("Invalid or expired token", util.ErrCodeUnauthorized)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, util.NewError("Invalid claims type", util.ErrCodeUnauthorized)
	}

	return claims, nil
}

/**************** REFRESH TOKEN ****************/

func (jg *JWTGenerator) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	return RefreshToken{
		Token:     token,
		UUID:      user.Uuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil
}

func (jg *JWTGenerator) StoreRefreshToken(ctx context.Context, token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return util.NewError("Refresh token already expired.", util.ErrCodeBadRequest)
	}

	return jg.redis.Set(ctx, cacheKey, token, ttl)
}

func (jg *JWTGenerator) ValidateRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	err := jg.redis.Get(ctx, cacheKey, &refreshToken)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return RefreshToken{}, util.WrapError(err, "Refresh token not found.", util.ErrCodeUnauthorized)
		}
		return RefreshToken{}, util.WrapError(err, "Cannot retrieve refresh token.", util.ErrCodeInternal)
	}
	if refreshToken.Revoked {
		return RefreshToken{}, util.NewError("Refresh token is revoked.", util.ErrCodeBadRequest)
	}
	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		return RefreshToken{}, util.NewError("Refresh token is expired.", util.ErrCodeBadRequest)
	}

	return refreshToken, nil
}

// RevokeRefreshToken marks a refresh token as revoked.
// The token is not deleted from Redis to prevent replay attacks and to safely handle concurrent refresh requests.
// This operation overwrites the existing Redis key with the revoked state.
func (jg *JWTGenerator) RevokeRefreshToken(ctx context.Context, oldToken string) error {
	cacheKey := "refresh_token:" + oldToken

	var refreshToken RefreshToken
	if err := jg.redis.Get(ctx, cacheKey, &refreshToken); err != nil {
		return util.WrapError(err, "Refresh token not found.", util.ErrCodeInternal)
	}
	if refreshToken.Revoked {
		return util.NewError("Refresh token already revoked.", util.ErrCodeBadRequest)
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return util.NewError("Refresh token expired.", util.ErrCodeUnauthorized)
	}

	refreshToken.Revoked = true
	// Overwrite the existing Redis key to persist the revoked state until expiration
	return jg.redis.Set(ctx, cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}
