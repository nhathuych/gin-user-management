package repository

import (
	"context"
	"gin-user-management/internal/db/sqlc"
)

type UserRepository interface {
	GetAll()
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update()
	Delete()
}
