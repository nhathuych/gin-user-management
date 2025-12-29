package app

import (
	"gin-user-management/internal/handler"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/route"
	"gin-user-management/internal/service"
)

type UserModule struct {
	routes route.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoute := route.NewUserRoute(userHandler)

	return &UserModule{
		routes: userRoute,
	}
}

func (m *UserModule) Routes() route.Route {
	return m.routes
}
