package handler

import (
	"gin-user-management/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAll(ctx *gin.Context) {
	uh.service.GetAll()
	ctx.JSON(http.StatusOK, "ok")
}

func (uh *UserHandler) Create(ctx *gin.Context) {
	uh.service.Create()
}

func (uh *UserHandler) GetByUUID(ctx *gin.Context) {
	uh.service.GetByUUID()
}

func (uh *UserHandler) Update(ctx *gin.Context) {
	uh.service.Update()
}

func (uh *UserHandler) Delete(ctx *gin.Context) {
	uh.service.Delete()
}
