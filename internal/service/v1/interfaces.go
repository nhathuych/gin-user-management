package v1

import (
	"gin-user-management/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAll()
	Create(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx *gin.Context, uuid uuid.UUID) error
}
