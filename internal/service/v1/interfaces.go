package v1

import (
	"gin-user-management/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetAll()
	Create(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	Delete()
}
