package repository

import (
	"context"
	"gin-user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type SqlUserRepository struct {
	db sqlc.Querier
}

func NewSqlUserRepository(db sqlc.Querier) UserRepository {
	return &SqlUserRepository{
		db: db,
	}
}

func (sur *SqlUserRepository) GetAll() {}

func (sur *SqlUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := sur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (sur *SqlUserRepository) GetByUUID() {}

func (sur *SqlUserRepository) Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := sur.db.UpdateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (sur *SqlUserRepository) SoftDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := sur.db.SoftDeleteUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (sur *SqlUserRepository) RestoreUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := sur.db.RestoreUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (sur *SqlUserRepository) HardDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := sur.db.HardDeleteUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}
