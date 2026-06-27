package config

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server
	ServerHost string
	ServerPort string

	// PostgreSQL
	PostgresDSN string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Kafka
	KafkaBrokers []string

	// JWT
	JWTSecret        string
	JWTRefreshSecret string

	// Logging
	LogLevel string
}

// NewConfig loads configuration from environment variables.
func NewConfig() *Config {
	logLevel := getEnv("LOG_LEVEL", "info")
	slog.LevelDebug
	slog.Info("loading configuration", "log_level", logLevel)

	return &Config{
		ServerHost:       getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		PostgresDSN:      getEnv("POSTGRES_DSN", "postgres://user:password@localhost:5432/ytd?sslmode=disable"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnvAsInt("REDIS_DB", 0),
		KafkaBrokers:     getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret"),
		LogLevel:         logLevel,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvAsSlice(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return fallback
}
