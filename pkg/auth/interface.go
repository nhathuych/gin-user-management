package auth

import (
	"gin-user-management/internal/db/sqlc"
)

type TokenGenerator interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken()
	ParseWithClaims(tokenString string) (*CustomClaims, error)
}
