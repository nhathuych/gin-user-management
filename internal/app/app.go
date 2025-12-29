package app

import (
	"gin-user-management/internal/config"
	"gin-user-management/internal/route"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() route.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
}

func NewApplication(cfg *config.Config) *Application {
	r := gin.Default()

	modules := []Module{
		NewUserModule(),
		// Add new module here
	}

	route.RegisterRoutes(r, getModuleRoutes(modules)...)

	return &Application{
		config: cfg,
		router: r,
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
