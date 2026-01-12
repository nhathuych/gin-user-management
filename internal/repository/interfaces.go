package repository

import (
	"context"
	"gin-user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error)
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
}
