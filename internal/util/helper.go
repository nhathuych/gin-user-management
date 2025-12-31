package util

import (
	"os"
)

func GetEnv(key, defautValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defautValue
}
