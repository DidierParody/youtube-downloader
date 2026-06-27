package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DuckDBPath    string
	RedisAddr     string
	RedisPassword string
	IcebergPath   string
}

func Load() *Config {
	_ = godotenv.Load(".env")

	return &Config{
		Port:          getEnv("ANALYTICS_API_PORT", "3001"),
		DuckDBPath:    getEnv("DUCKDB_DATABASE_URL", "./analytics.duckdb"),
		RedisAddr:     getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		IcebergPath:   getEnv("ICEBERG_CATALOG_PATH", "./iceberg"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
