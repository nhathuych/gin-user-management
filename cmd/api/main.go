package main

import (
	"gin-user-management/internal/app"
	"gin-user-management/internal/config"
)

func main() {
	cfg := config.NewConfig()

	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		panic(err)
	}
}
