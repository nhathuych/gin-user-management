package v1

import (
	dtoV1 "gin-user-management/internal/dto/v1"
	serviceV1 "gin-user-management/internal/service/v1"
	"gin-user-management/internal/util"
	"gin-user-management/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service serviceV1.UserService
}

func NewUserHandler(service serviceV1.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// =========== ENDPOINTS ===========

func (uh *UserHandler) GetAll(ctx *gin.Context) {
	var params dtoV1.GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	users, total, err := uh.service.GetAll(ctx, params.Search, params.Order, params.Sort, params.Page, params.Limit)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	dtos := dtoV1.MapUsersToDTOs(users)
	paginatedUsers := util.NewPaginationResponse(dtos, params.Page, params.Limit, total)
	util.ResponseSuccess(ctx, http.StatusOK, "Users retrieved successfully.", paginatedUsers)
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
	util.ResponseSuccess(ctx, http.StatusCreated, "User created successfully.", dto)
}

func (uh *UserHandler) GetByUUID(ctx *gin.Context) {
	var params dtoV1.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	util.ResponseSuccess(ctx, http.StatusOK, "")
}

func (uh *UserHandler) Update(ctx *gin.Context) {
	var params dtoV1.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	uuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	var input dtoV1.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapUpdateInputToModel(uuid)

	updatedUser, err := uh.service.Update(ctx, user)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	dto := dtoV1.MapUserToDTO(updatedUser)
	util.ResponseSuccess(ctx, http.StatusOK, "User updated successfullt.", dto)
}

func (uh *UserHandler) SoftDeleteUser(ctx *gin.Context) {
	var params dtoV1.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	uuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	deletedUser, err := uh.service.SoftDeleteUser(ctx, uuid)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	dto := dtoV1.MapUserToDTO(deletedUser)
	util.ResponseSuccess(ctx, http.StatusOK, "User deleted successfully.", dto)
}

func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var params dtoV1.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	uuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	restoredUser, err := uh.service.RestoreUser(ctx, uuid)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	dto := dtoV1.MapUserToDTO(restoredUser)
	util.ResponseSuccess(ctx, http.StatusOK, "User restored successfully.", dto)
}

func (uh *UserHandler) HardDeleteUser(ctx *gin.Context) {
	var params dtoV1.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		util.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	uuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	err = uh.service.HardDeleteUser(ctx, uuid)
	if err != nil {
		util.ResponseError(ctx, err)
		return
	}

	util.ResponseStatusCode(ctx, http.StatusNoContent)
}
