package v1

import (
	handlerV1 "gin-user-management/internal/handler/v1"
	"gin-user-management/internal/middleware"
	"gin-user-management/pkg/auth"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	handler      *handlerV1.UserHandler
	jwtGenerator auth.TokenGenerator
}

func NewUserRoute(handler *handlerV1.UserHandler, jwtGenerator auth.TokenGenerator) *UserRoute {
	return &UserRoute{
		handler:      handler,
		jwtGenerator: jwtGenerator,
	}
}

func (ur *UserRoute) Register(r *gin.RouterGroup) {
	users := r.Group("/users")

	// 🌐 PUBLIC
	users.GET("", ur.handler.GetAll)
	users.GET("/:uuid", ur.handler.GetByUUID)

	// 🔒 PROTECTED
	protected := users.Group("")
	protected.Use(middleware.AuthMiddleware(ur.jwtGenerator))
	{
		protected.GET("/deleted", ur.handler.GetDeletedUsers)
		protected.POST("", ur.handler.Create)
		protected.PATCH("/:uuid", ur.handler.Update)

		protected.DELETE("/:uuid", ur.handler.SoftDeleteUser)
		protected.PUT("/:uuid/restore", ur.handler.RestoreUser)
		protected.DELETE("/:uuid/force", ur.handler.HardDeleteUser)
	}
}
