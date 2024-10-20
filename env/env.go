package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/shrimpsizemoose/trekker/logger"
)

type EnvError struct {
	Key     string
	Message string
}

// загрузить переменные из .env файла, молча игнорируя если файла нет
func LoadEnv() {
	_ = godotenv.Load()
}

// GetEnvOrDefault("PORT", "8000, "использую порт по умолчанию")
// -> "PORT не задан, использую порт по умолчанию: 8000"
func GetEnvOrDefault(key, defaultValue string, defaultComment string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultComment != "" {
		logger.Warn.Printf("%s не задан, %s: %s", key, defaultComment, defaultValue)
	}
	return defaultValue
}

// проверяет, что нужные переменные окружения (ключи мапы) есть
// и выводит сообщения об ошибках (значения мапы) если нет
func RequireEnv(requirements map[string]string) {
	var errors []EnvError

	for k, v := range requirements {
		if _, exists := os.LookupEnv(k); !exists {
			errors = append(errors, EnvError{Key: k, Message: v})
		}
	}

	if len(errors) > 0 {
		logger.Error.Printf("Мне не хватает некоторых переменных окружения 👇")
		for _, err := range errors {
			if err.Message == "" {
				logger.Error.Printf("Нет переменной окружения: %v", err.Key)
			} else {
				logger.Error.Printf("%s: %s", err.Key, err.Message)
			}
		}
		os.Exit(1)
	}
}
