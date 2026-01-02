package config

import (
	"fmt"
	"gin-user-management/internal/util"
)

type DBConfig struct {
	DBName   string
	Host     string
	Port     string
	User     string
	Password string
	SSLMode  string
}

type Config struct {
	ServerAddress string
	DB            DBConfig
}

func NewConfig() *Config {
	return &Config{
		ServerAddress: fmt.Sprintf(":%s", util.GetEnv("PORT", "8080")),
		DB: DBConfig{
			DBName:   util.GetEnv("POSTGRES_DB", "gin_user_management_development"),
			Host:     util.GetEnv("POSTGRES_HOST", "localhost"),
			Port:     util.GetEnv("POSTGRES_PORT", "5432"),
			User:     util.GetEnv("POSTGRES_USER", "postgres"),
			Password: util.GetEnv("POSTGRES_PASSWORD", "postgres"),
			SSLMode:  util.GetEnv("POSTGRES_SSLMODE", "disable"),
		},
	}
}

func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.SSLMode)
}
