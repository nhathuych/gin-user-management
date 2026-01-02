package app

import (
	"gin-user-management/internal/config"
	"gin-user-management/internal/db"
	"gin-user-management/internal/route"
	"gin-user-management/internal/validation"
	"log"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() route.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module // no need yet
}

func NewApplication(cfg *config.Config) *Application {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Init validator failed: %v", err)
	}

	r := gin.Default()

	modules := []Module{
		NewUserModule(),
		// Add new module here
	}

	route.RegisterRoutes(r, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules, // no need yet
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
