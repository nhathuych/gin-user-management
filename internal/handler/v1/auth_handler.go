package v1

import (
	serviceV1 "gin-user-management/internal/service/v1"
	"gin-user-management/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service serviceV1.AuthService
}

func NewAuthHandler(service serviceV1.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	util.ResponseMessage(ctx, http.StatusOK, "Logged in successfully.")
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	util.ResponseMessage(ctx, http.StatusOK, "Logged out successfully.")
}
