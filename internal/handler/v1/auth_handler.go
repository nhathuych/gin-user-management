package v1

import (
	dtoV1 "gin-user-management/internal/dto/v1"
	serviceV1 "gin-user-management/internal/service/v1"
	"gin-user-management/internal/util"
	"gin-user-management/internal/validation"
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
	var input dtoV1.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiresIn, err := ah.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	response := dtoV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	util.ResponseSuccess(ctx, http.StatusOK, "Logged in successfully.", response, nil)
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	util.ResponseMessage(ctx, http.StatusOK, "Logged out successfully.")
}

func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var input dtoV1.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiresIn, err := ah.service.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	response := dtoV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	util.ResponseSuccess(ctx, http.StatusOK, "Token refreshed successfully.", response, nil)
}
