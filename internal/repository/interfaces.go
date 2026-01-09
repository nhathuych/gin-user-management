package repository

import (
	"context"
	"gin-user-management/internal/db/sqlc"
)

type UserRepository interface {
	GetAll()
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	Delete()
}
