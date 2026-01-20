package auth

import (
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTGenerator struct{}

type Claims struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	Role  int32  `json:"role"`
	jwt.RegisteredClaims
}

const (
	AccessTokenTTL = 15 * time.Minute
)

var (
	jwtSecret = []byte(util.GetEnv("JWT_SECRET", "gin-user-management-default-secret"))
)

func NewJWTGenerator() *JWTGenerator {
	return &JWTGenerator{}
}

func (js *JWTGenerator) GenerateAccessToken(user sqlc.User) (string, error) {
	claims := &Claims{
		UUID:  user.Uuid.String(),
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (js *JWTGenerator) GenerateRefreshToken() {
}
