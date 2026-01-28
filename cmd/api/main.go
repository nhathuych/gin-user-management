package main

import (
	"gin-user-management/internal/app"
	"gin-user-management/internal/config"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/logger"
)

func main() {
	util.InitEnv()

	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   "log/app.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      util.IsDevelopment(),
	})

	cfg := config.NewConfig()
	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		panic(err)
	}
}
