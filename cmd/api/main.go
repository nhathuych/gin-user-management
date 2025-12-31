package main

import (
	"gin-user-management/internal/app"
	"gin-user-management/internal/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	cfg := config.NewConfig()
	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		panic(err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("No .env file found.")
	}
}
