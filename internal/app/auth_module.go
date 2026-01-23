package app

import (
	handlerV1 "gin-user-management/internal/handler/v1"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/route"
	routeV1 "gin-user-management/internal/route/v1"
	serviceV1 "gin-user-management/internal/service/v1"
	"gin-user-management/pkg/auth"
)

type AuthModule struct {
	routes route.Route
}

func NewAuthModule(ctx *ModuleContext, tokenGenerator auth.TokenGenerator) *AuthModule {
	userRepo := repository.NewSqlUserRepository(ctx.DB)
	authService := serviceV1.NewAuthService(userRepo, ctx.RedisCache, tokenGenerator)
	authHandler := handlerV1.NewAuthHandler(authService)
	authRoute := routeV1.NewAuthRoute(authHandler)

	return &AuthModule{
		routes: authRoute,
	}
}

func (m *AuthModule) Routes() route.Route {
	return m.routes
}
