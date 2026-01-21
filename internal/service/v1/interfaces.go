package v1

import (
	"gin-user-management/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAll(ctx *gin.Context, search, orderBy, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error)
	Create(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	Update(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx *gin.Context, uuid uuid.UUID) error
}

type AuthService interface {
	Login(ctx *gin.Context, email, password string) (string, string, int, error)
	Logout(ctx *gin.Context) error
}
