package config

import (
	"os"

	"github.com/joho/godotenv"
)

const defaultDBURL = "postgres://postgres:postgres@localhost:5432/his?sslmode=disable"
const defaultAppPort = "8080"

type Config struct {
	AppPort   string
	DBURL     string
	AppSecret string
}

func Load() Config {
	// Load .env file if present (silently ignore if not found)
	_ = godotenv.Load()

	return Config{
		AppPort:   getEnv("APP_PORT", defaultAppPort),
		DBURL:     getEnv("DB_URL", defaultDBURL),
		AppSecret: os.Getenv("APP_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
