package route

import (
	"gin-user-management/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	handler *handler.UserHandler
}

func NewUserRoute(handler *handler.UserHandler) *UserRoute {
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
		users.PUT("/:uuid", ur.handler.Update)
		users.DELETE("/:uuid", ur.handler.Delete)
	}
}
