package repository

import (
	"context"
	"gin-user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	CountUsers(ctx context.Context, search string, deleted bool) (int64, error)
	GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error)
	GetAllV2(ctx context.Context, search, orderBy, sort string, limit, offset int32, deleted bool) ([]sqlc.User, error)
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	GetByEmail(ctx context.Context, email string) (sqlc.User, error)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
}
