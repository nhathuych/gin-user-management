package main

import (
	"gin-user-management/internal/config"
	"gin-user-management/internal/handler"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/route"
	"gin-user-management/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoute := route.NewUserRoute(userHandler)

	r := gin.Default()
	route.RegisterRoutes(r, userRoute)

	if err := r.Run(cfg.ServerAddress); err != nil {
		panic(err)
	}
}
