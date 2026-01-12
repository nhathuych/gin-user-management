package repository

import (
	"context"
	"fmt"
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

func (sur *SqlUserRepository) GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error) {
	var (
		users []sqlc.User
		err   error
	)

	switch orderBy {
	case "id":
		switch sort {
		case "asc":
			users, err = sur.db.ListUsersOrderByIdASC(ctx, sqlc.ListUsersOrderByIdASCParams{
				Limit: limit, Offset: offset, Search: search,
			})
		case "desc":
			users, err = sur.db.ListUsersOrderByIdDESC(ctx, sqlc.ListUsersOrderByIdDESCParams{
				Limit: limit, Offset: offset, Search: search,
			})
		default:
			return users, fmt.Errorf("invalid sort: %s", sort)
		}
	default:
		return users, fmt.Errorf("invalid orderBy: %s", orderBy)
	}

	return users, err
}

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
