package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
