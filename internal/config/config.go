package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	JWTSecret       string
	AccessExpireMin int
	RefreshExpireHr int
}

func LoadConfig() *Config {
	accessExpireStr := os.Getenv("JWT_ACCESS_EXPIRE_MINUTES")
	refreshExpireStr := os.Getenv("JWT_REFRESH_EXPIRE_HOURS")

	accessExpire, err := strconv.Atoi(accessExpireStr)
	if err != nil || accessExpire <= 0 {
		log.Println(" invalid JWT_ACCESS_EXPIRE_MINUTES, fallback to 15")
		accessExpire = 15
	}

	refreshExpire, err := strconv.Atoi(refreshExpireStr)
	if err != nil || refreshExpire <= 0 {
		log.Println(" invalid JWT_REFRESH_EXPIRE_HOURS, fallback to 168")
		refreshExpire = 168
	}

	return &Config{
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AccessExpireMin: accessExpire,
		RefreshExpireHr: refreshExpire,
	}
}
