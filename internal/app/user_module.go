package app

import (
	handlerV1 "gin-user-management/internal/handler/v1"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/route"
	routeV1 "gin-user-management/internal/route/v1"
	serviceV1 "gin-user-management/internal/service/v1"
	"gin-user-management/pkg/auth"
)

type UserModule struct {
	routes route.Route
}

func NewUserModule(ctx *ModuleContext, jwtGenerator auth.TokenGenerator) *UserModule {
	userRepo := repository.NewSqlUserRepository(ctx.DB)
	userService := serviceV1.NewUserService(userRepo, ctx.RedisCache)
	userHandler := handlerV1.NewUserHandler(userService)
	userRoute := routeV1.NewUserRoute(userHandler, ctx.RedisCache, jwtGenerator)

	return &UserModule{
		routes: userRoute,
	}
}

func (m *UserModule) Routes() route.Route {
	return m.routes
}
