package repository

import (
	"context"
	"gin-user-management/internal/db/sqlc"
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

func (sur *SqlUserRepository) Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := sur.db.CreateUser(ctx, userParams)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (sur *SqlUserRepository) GetByUUID() {}

func (sur *SqlUserRepository) Update() {}

func (sur *SqlUserRepository) Delete() {}
