package app

import (
	"gin-user-management/internal/config"
	"gin-user-management/internal/db"
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/route"
	"gin-user-management/internal/validation"
	"gin-user-management/pkg/auth"
	"gin-user-management/pkg/cache"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	DB    sqlc.Querier
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) *Application {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Init validator failed: %v", err)
	}

	r := gin.Default()

	redisClient := config.NewRedisClient()
	redisCacheService := cache.NewRedisCacheService(redisClient)

	tokenGenerator := auth.NewJWTGenerator(redisCacheService)

	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
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
