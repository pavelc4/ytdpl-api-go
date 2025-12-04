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
	R2Config   R2Config
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Endpoint        string
	PublicURL       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using environment variables")
	}

	return &Config{
		Port:       getEnv("PORT", "5000"),
		CookiePath: getEnv("COOKIE_PATH", ""),
		APIVersion: getEnv("API_VERSION", "v1"),
		R2Config: R2Config{
			AccountID:       getEnv("R2_ACCOUNT_ID", ""),
			AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
			BucketName:      getEnv("R2_BUCKET_NAME", ""),
			Endpoint:        getEnv("R2_ENDPOINT", ""),
			PublicURL:       getEnv("R2_PUBLIC_URL", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
