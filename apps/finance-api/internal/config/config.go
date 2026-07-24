package config

import "os"

type Config struct {
	Port           string
	AllowedOrigins string
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: getEnv("FINANCE_API_ALLOWED_ORIGINS", ""),
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
