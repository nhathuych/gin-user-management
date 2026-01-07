package v1

import (
	dtoV1 "gin-user-management/internal/dto/v1"
	serviceV1 "gin-user-management/internal/service/v1"
	"gin-user-management/internal/util"
	"gin-user-management/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service serviceV1.UserService
}

func NewUserHandler(service serviceV1.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type GetUserByUuidParam struct {
	Uuid string `uri:"uuid" binding:"uuid"`
}

type GetUsersParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lte=100"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

// =========== ENDPOINTS ===========

func (uh *UserHandler) GetAll(ctx *gin.Context) {
	var params GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	util.ResponseSuccess(ctx, http.StatusOK, "")
}

func (uh *UserHandler) Create(ctx *gin.Context) {
	var input dtoV1.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapCreateInputToModel()

	createdUser, err := uh.service.Create(ctx, user)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	dto := dtoV1.MapUserToDTO(createdUser)
	util.ResponseSuccess(ctx, http.StatusCreated, dto)
}

func (uh *UserHandler) GetByUUID(ctx *gin.Context) {
	var params GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	util.ResponseSuccess(ctx, http.StatusOK, "")
}

func (uh *UserHandler) Update(ctx *gin.Context) {
	var params GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var input dtoV1.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	util.ResponseSuccess(ctx, http.StatusOK, "")
}

func (uh *UserHandler) Delete(ctx *gin.Context) {
	var params GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	util.ResponseStatusCode(ctx, http.StatusNoContent)
}
