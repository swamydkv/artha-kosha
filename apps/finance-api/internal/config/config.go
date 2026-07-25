package config

import (
	"os"
	"strconv"
	"time"
)

const (
	DefaultPort                 = "8080"
	DefaultOutboxRetentionDays  = 7
	DefaultArchiveRetentionDays = 30
	DefaultSessionTTL           = 2 * time.Hour
	DefaultWorkerInterval       = 5 * time.Second
	DefaultPruneInterval        = 24 * time.Hour
	DefaultRequestTimeout       = 30 * time.Second
)

type Config struct {
	Port                 string
	AllowedOrigins       string
	OutboxRetentionDays  int
	ArchiveRetentionDays int
	SessionTTL           time.Duration
	WorkerInterval       time.Duration
	PruneInterval        time.Duration
	RequestTimeout       time.Duration
}

func Load() Config {
	return Config{
		Port:                 getEnv("PORT", DefaultPort),
		AllowedOrigins:       getEnv("FINANCE_API_ALLOWED_ORIGINS", ""),
		OutboxRetentionDays:  getEnvInt("OUTBOX_RETENTION_DAYS", DefaultOutboxRetentionDays),
		ArchiveRetentionDays: getEnvInt("ARCHIVE_RETENTION_DAYS", DefaultArchiveRetentionDays),
		SessionTTL:           getEnvDuration("SESSION_TTL", DefaultSessionTTL),
		WorkerInterval:       getEnvDuration("WORKER_INTERVAL", DefaultWorkerInterval),
		PruneInterval:        getEnvDuration("PRUNE_INTERVAL", DefaultPruneInterval),
		RequestTimeout:       getEnvDuration("REQUEST_TIMEOUT", DefaultRequestTimeout),
	}
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}
