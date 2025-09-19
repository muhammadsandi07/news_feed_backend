package config

import (
	"os"
)

type Config struct {
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
