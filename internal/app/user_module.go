package app

import (
	handlerV1 "gin-user-management/internal/handler/v1"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/route"
	routeV1 "gin-user-management/internal/route/v1"
	serviceV1 "gin-user-management/internal/service/v1"
)

type UserModule struct {
	routes route.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewSqlUserRepository()
	userService := serviceV1.NewUserService(userRepo)
	userHandler := handlerV1.NewUserHandler(userService)
	userRoute := routeV1.NewUserRoute(userHandler)

	return &UserModule{
		routes: userRoute,
	}
}

func (m *UserModule) Routes() route.Route {
	return m.routes
}
