package main

import (
	"gin-user-management/internal/app"
	"gin-user-management/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from a .env file if it exists.
	// The compiled Go binary does NOT embed any environment variables.
	// All configuration is always read from the process environment at runtime.
	// In containerized or production environments, variables are typically injected by Docker, Kubernetes, or the operating system instead of a .env file.
	godotenv.Load()

	cfg := config.NewConfig()
	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		panic(err)
	}
}
