package v1

import (
	"gin-user-management/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetAll()
	Create(ctx *gin.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update()
	Delete()
}
