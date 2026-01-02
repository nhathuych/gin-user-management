package repository

import "gin-user-management/internal/db/sqlc"

type SqlUserRepository struct {
	db *sqlc.Queries
}

func NewSqlUserRepository(DB *sqlc.Queries) UserRepository {
	return &SqlUserRepository{
		db: DB,
	}
}

func (sur *SqlUserRepository) GetAll() {}

func (sur *SqlUserRepository) Create() {}

func (sur *SqlUserRepository) GetByUUID() {}

func (sur *SqlUserRepository) Update() {}

func (sur *SqlUserRepository) Delete() {}
