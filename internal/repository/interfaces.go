package repository

import (
	"context"
	"gin-user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll()
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
}
