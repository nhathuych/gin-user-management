package auth

import (
	"context"
	"gin-user-management/internal/db/sqlc"
)

type TokenGenerator interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	ParseWithClaims(tokenString string) (*CustomClaims, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	StoreRefreshToken(ctx context.Context, token RefreshToken) error
	ValidateRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, oldToken string) error
}
