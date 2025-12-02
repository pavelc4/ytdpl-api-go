package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	CookiePath string
	APIVersion string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using environment variables")
	}

	return &Config{
		Port:       getEnv("PORT", "5000"),
		CookiePath: getEnv("COOKIE_PATH", ""),
		APIVersion: getEnv("API_VERSION", "v1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
