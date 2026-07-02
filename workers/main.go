package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	brokers := getEnv("KAFKA_BROKERS", "localhost:9092")

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers),
		kgo.ConsumerGroup("workers"),
		kgo.ConsumeTopics("download-events", "metadata-events", "ai-events", "media-events"),
	)
	if err != nil {
		slog.Error("failed to create kafka client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	slog.Info("workers started", "brokers", brokers)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fetches := client.PollFetches(ctx)
		if ctx.Err() != nil {
			return
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, record := range p.Records {
				processRecord(client, record)
			}
		})
	}
}

func processRecord(client *kgo.Client, record *kgo.Record) {
	slog.Info("processing",
		"topic", record.Topic,
		"partition", record.Partition,
		"offset", record.Offset,
	)

	switch record.Topic {
	case "download-events":
		handleDownload(client, record.Value)
	case "metadata-events":
		handleMetadata(client, record.Value)
	case "ai-events":
		handleAI(client, record.Value)
	case "media-events":
		handleMedia(client, record.Value)
	}
}

func handleDownload(client *kgo.Client, data []byte) {
	slog.Info("DOWNLOAD: processing", "data", string(data))
	time.Sleep(1 * time.Second)
	slog.Info("DOWNLOAD: completed")
	client.Produce(context.Background(), &kgo.Record{
		Topic: "download-completed",
		Value: data,
	}, nil)
}

func handleMetadata(client *kgo.Client, data []byte) {
	slog.Info("METADATA: processing", "data", string(data))
	time.Sleep(500 * time.Millisecond)
	slog.Info("METADATA: completed")
}

func handleAI(client *kgo.Client, data []byte) {
	slog.Info("AI: processing", "data", string(data))
	time.Sleep(2 * time.Second)
	slog.Info("AI: completed")
}

func handleMedia(client *kgo.Client, data []byte) {
	slog.Info("MEDIA: processing", "data", string(data))
	time.Sleep(1 * time.Second)
	slog.Info("MEDIA: completed")
}
