package repository

import (
	"context"
	"fmt"
	"gin-user-management/internal/db"
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

func (sur *SqlUserRepository) CountUsers(ctx context.Context, search string, deleted bool) (int64, error) {
	total, err := sur.db.CountUsers(ctx, sqlc.CountUsersParams{
		Search:  search,
		Deleted: &deleted,
	})
	if err != nil {
		return 0, err
	}
	return total, nil
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

func (sur *SqlUserRepository) GetAllV2(ctx context.Context, search, orderBy, sort string, limit, offset int32, deleted bool) ([]sqlc.User, error) {
	query := `
		SELECT *
		FROM users
		WHERE
			(
				$1::text IS NULL OR
				$1::text = '' OR
				email    ILIKE '%' || $1 || '%' OR
				fullname ILIKE '%' || $1 || '%'
			)
	`

	if deleted {
		query += " AND deleted_at IS NOT NULL"
	} else {
		query += " AND deleted_at IS NULL"
	}

	order := "ASC"
	if sort == "desc" {
		order = "DESC"
	}

	if orderBy == "id" {
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)
	} else {
		query += " ORDER BY id ASC"
	}

	query += " LIMIT $2 OFFSET $3"

	rows, err := db.DBPool.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []sqlc.User{}
	for rows.Next() {
		var i sqlc.User
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Email,
			&i.Password,
			&i.Fullname,
			&i.Age,
			&i.Status,
			&i.Role,
			&i.DeletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (sur *SqlUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := sur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (sur *SqlUserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	user, err := sur.db.GetUser(ctx, uuid)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

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
