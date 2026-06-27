package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration values for the workers

type Config struct {
	// Redpanda (Kafka) settings
	KafkaBrokers []string
	KafkaGroupID string

	// Database settings
	DatabaseURL string

	// Object storage (MinIO/R2) settings
	StorageEndpoint  string
	StorageAccessKey string
	StorageSecretKey string
	StorageBucket    string
	StorageRegion    string
	StorageSecure    bool

	// Application settings
	WorkerName string
	LogLevel   string
}

// LoadEnv loads configuration from environment variables
func LoadEnv() (*Config, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092" // default
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "youtube-downloader-workers"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return &Config{
		KafkaBrokers:     strings.Split(brokers, ","),
		KafkaGroupID:     groupID,
		DatabaseURL:      databaseURL,
		StorageEndpoint:  os.Getenv("STORAGE_ENDPOINT"),
		StorageAccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		StorageSecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		StorageBucket:    os.Getenv("STORAGE_BUCKET"),
		StorageRegion:    os.Getenv("STORAGE_REGION"),
		StorageSecure:    os.Getenv("STORAGE_SECURE") == "true",
		WorkerName:       os.Getenv("WORKER_NAME"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
	}, nil
}
