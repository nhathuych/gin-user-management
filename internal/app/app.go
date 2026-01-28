package app

import (
	"gin-user-management/internal/config"
	"gin-user-management/internal/db"
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/route"
	"gin-user-management/internal/validation"
	"gin-user-management/pkg/auth"
	"gin-user-management/pkg/cache"
	"gin-user-management/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() route.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB         sqlc.Querier
	RedisCache cache.RedisCacheService
}

func NewApplication(cfg *config.Config) *Application {
	if err := db.InitDB(); err != nil {
		logger.AppLogger.Fatal().Err(err).Msg("🔴 Failed to connect to database")
	}

	if err := validation.InitValidator(); err != nil {
		logger.AppLogger.Fatal().Err(err).Msg("🔴 Init validator failed")
	}

	r := gin.Default()

	redisClient := config.NewRedisClient()
	redisCacheService := cache.NewRedisCacheService(redisClient)

	tokenGenerator := auth.NewJWTGenerator(redisCacheService)

	ctx := &ModuleContext{
		DB:         db.DB,
		RedisCache: redisCacheService,
	}

	modules := []Module{
		NewUserModule(ctx, tokenGenerator),
		NewAuthModule(ctx, tokenGenerator),
		// Add new module here
	}

	route.RegisterRoutes(r, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}

func getModuleRoutes(modules []Module) []route.Route {
	routeList := make([]route.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
