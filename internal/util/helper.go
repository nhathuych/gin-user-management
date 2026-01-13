package util

import (
	"os"
	"strconv"
)

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
