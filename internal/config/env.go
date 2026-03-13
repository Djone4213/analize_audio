package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() *Config {
	_ = godotenv.Load(".env")

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Name:     getEnv("DB_NAME", "p"),
			User:     getEnv("DB_USER", "p"),
			Password: getEnv("DB_PASSWORD", "p"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
		},
		App: AppConfig{
			Dir: getEnv("DIR", "/public/"),
		},
		Bot: BotConfig{
			Token: getEnv("BOT_TOKEN", ""),
			URL:   getEnv("BOT_URL", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
