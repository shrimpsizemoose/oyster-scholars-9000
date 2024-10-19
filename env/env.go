package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/shrimpsizemoose/trekker/logger"
)

func LoadEnv() {
	_ = godotenv.Load()
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func RequireEnv(keys ...string) {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			logger.Error.Fatalf("Переменная окружения %s на задана, не могу продолжать", key)
		}
	}
}
