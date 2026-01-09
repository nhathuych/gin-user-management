package v1

import (
	handlerV1 "gin-user-management/internal/handler/v1"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	handler *handlerV1.UserHandler
}

func NewUserRoute(handler *handlerV1.UserHandler) *UserRoute {
	return &UserRoute{
		handler: handler,
	}
}

func (ur *UserRoute) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("", ur.handler.GetAll)
		users.POST("", ur.handler.Create)
		users.GET("/:uuid", ur.handler.GetByUUID)
		users.PATCH("/:uuid", ur.handler.Update)
		users.DELETE("/:uuid", ur.handler.Delete)
	}
}
