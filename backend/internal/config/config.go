package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken  string
	ServerPort        string
	PostgresHost      string
	PostgresUser      string
	PostgresPassword  string
	PostgresDB        string
	PostgresPort      string
	AvitoClientId     string
	AvitoClientSecret string
	CookieEncryptionKey string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	return &Config{
		TelegramBotToken:  getEnv("BOT_TOKEN", ""),
		ServerPort:        getEnv("BACKEND_PORT", "8080"),
		PostgresHost:      getEnv("POSTGRES_HOST", "localhost"),
		PostgresUser:      getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword:  getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:        getEnv("POSTGRES_DB", "miniapp"),
		PostgresPort:      getEnv("POSTGRES_PORT", "5432"),
		AvitoClientId:     getEnv("AVITO_CLIENT_ID", ""),
		AvitoClientSecret: getEnv("AVITO_CLIENT_SECRET", ""),
		CookieEncryptionKey: getEnv("COOKIE_ENCRYPTION_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
