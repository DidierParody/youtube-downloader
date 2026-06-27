package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DuckDBPath    string
	S3Endpoint    string
	S3AccessKey   string
	S3SecretKey   string
	S3Bucket      string
	S3Region      string
	S3UseSSL      bool
	KafkaBrokers  string
}

func Load() *Config {
	_ = godotenv.Load(".env")

	return &Config{
		Port:         getEnv("EXPORT_SERVICE_PORT", "3002"),
		DuckDBPath:   getEnv("DUCKDB_DATABASE_URL", "./analytics.duckdb"),
		S3Endpoint:   getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3AccessKey:  getEnv("SKey_endpoint_ACCESS_KEY_ID", "minioadmin"),
		S3SecretKey:  getEnv("S3_SECRET_ACCESS_KEY", "minioadmin"),
		S3Bucket:     getEnv("S3_BUCKET_EXPORTS", "ytd-exports"),
		S3Region:     getEnv("AWS_REGION", "us-east-1"),
		S3UseSSL:     getEnv("S3_USE_SSL", "false") == "true",
		KafkaBrokers: getEnv("REDPANDA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
