package config

import (
	"fmt"
	"gin-user-management/internal/util"
)

type Config struct {
	ServerAddress string
}

func NewConfig() *Config {
	return &Config{
		ServerAddress: fmt.Sprintf(":%s", util.GetEnv("PORT", "8080")),
	}
}
