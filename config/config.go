package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   int
	SecretKey string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppPort:   envInt("APP_PORT", 3000),
		SecretKey: mustEnv("SECRET_KEY"),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing env: " + key)
	}
	return v
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic("invalid int env: " + key)
	}
	return i
}

func (cfg Config) AppPortString() string {
	return strconv.Itoa(cfg.AppPort)
}
