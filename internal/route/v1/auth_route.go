package v1

import (
	handlerV1 "gin-user-management/internal/handler/v1"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	handler *handlerV1.AuthHandler
}

func NewAuthRoute(handler *handlerV1.AuthHandler) *AuthRoute {
	return &AuthRoute{
		handler: handler,
	}
}

func (ar *AuthRoute) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", ar.handler.Login)
		auth.POST("/logout", ar.handler.Logout)
		auth.POST("/refresh-token", ar.handler.RefreshToken)
		auth.POST("/password/forgot", ar.handler.RequestPasswordReset)
		auth.POST("/password/reset", ar.handler.ResetPassword)
	}
}
