package util

import (
	"os"
	"runtime/debug"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	isDev       bool
	projectName string // Cached project module name
)

func InitEnv() {
	// Load environment variables from a .env file if it exists.
	// The compiled Go binary does NOT embed any environment variables.
	// All configuration is always read from the process environment at runtime.
	// In containerized or production environments, variables are typically injected by Docker, Kubernetes, or the operating system instead of a .env file.
	godotenv.Load()

	isDev = GetEnv("APP_ENV", "production") == "development"

	// Cache the project module name from go.mod once
	if info, ok := debug.ReadBuildInfo(); ok {
		projectName = info.Main.Path
	}
}

func IsDevelopment() bool {
	return isDev
}

func GetProjectName() string {
	return projectName
}

func GetEnv(key, defautValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defautValue
}

func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return val
}
