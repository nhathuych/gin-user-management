package auth

import (
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTGenerator struct{}

type CustomClaims struct {
	UUID string `json:"sub"`
	Role int32  `json:"role"`
	jwt.RegisteredClaims
}

const (
	AccessTokenTTL = 15 * time.Minute
)

var (
	jwtSecret = []byte(util.GetEnv("JWT_SECRET", "gin-user-management-default-secret"))
)

func NewJWTGenerator() TokenGenerator {
	return &JWTGenerator{}
}

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

func (jg *JWTGenerator) GenerateRefreshToken() {
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
