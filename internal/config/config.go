package config

import "os"

type Config struct {
	ServerAddress string
}

func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		ServerAddress: ":" + port,
	}
}
