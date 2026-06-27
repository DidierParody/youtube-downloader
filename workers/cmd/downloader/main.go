package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/DidierParody/youtube-downloader/workers/internal/config"
	"github.com/DidierParody/youtube-downloader/workers/internal/domain"
	"github.com/DidierParody/youtube-downloader/workers/internal/kafka"
	"github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	TopicDownloadEvents = "download.events"
	TopicAuditEvents    = "audit.events"
	TopicDLQ            = "download.dlq"
)

// EventEnvelope represents the generic event structure used by the consumer
type EventEnvelope struct {
	Header  domain.EventHeader `json:"header"`
	Payload json.RawMessage    `json:"payload"`
}

func main() {
	logger := logrus.New()
	logger.Info("Starting downloader worker")

	// Load environment config
	cfg, err := config.LoadEnv()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load config")
	}

	// Initialize MinIO client (R2 compatible)
	ctx := context.Background()
	minioClient, err := minio.New(cfg.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV2(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
		Secure: cfg.StorageSecure,
		Region: cfg.StorageRegion,
	})
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize MinIO client")
	}

	// Initialize producer for downstream events
	producer := kafka.NewProducer(cfg.KafkaBrokers, logger)
	defer producer.Close()

	// Process function for the consumer
	processFunc := func(ctx context.Context, msg kafka.Message) error {
		var envelope EventEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			logger.WithError(err).Error("Failed to unmarshal event")
			return err
		}

		// Check if this is the correct event type
		if envelope.Header.EventType != "DownloadRequested" {
			logger.Debugf("Skipping event of type %s", envelope.Header.EventType)
			return nil
		}

		// Parse the payload
		var downloadReq domain.DownloadRequested
		if err := json.Unmarshal(envelope.Payload, &downloadReq); err != nil {
			logger.WithError(err).Error("Failed to unmarshal DownloadRequested payload")
			return err
		}

		logger.Infof("Processing DownloadRequested: %s, URL: %s", envelope.Header.EventID, downloadReq.URL)

		// Idempotency check: check if already processed
		if isAlreadyProcessed(ctx, downloadReq.EventID) {
			logger.Infof("Event %s already processed, skipping", downloadReq.EventID)
			return nil
		}

		// Mark as started
		if err := publishEvent(ctx, producer, TopicDownloadEvents, domain.DownloadStarted{
			EventHeader: domain.EventHeader{
				EventID:       generateEventID(),
				CorrelationID: envelope.Header.CorrelationID,
				CausationID:   envelope.Header.EventID,
				EventType:     "DownloadStarted",
				EventVersion:  "1.0",
				Producer:      "downloader-worker",
				OccurredAt:    time.Now().UTC(),
			},
			DownloadID: downloadReq.EventID,
			URL:        downloadReq.URL,
			Status:     "started",
		}); err != nil {
			logger.WithError(err).Error("Failed to publish DownloadStarted event")
			return err
		}

		// Download the video using yt-dlp
		downloadedFile, err := downloadVideo(ctx, downloadReq.URL, logger)
		if err != nil {
			logger.WithError(err).Error("Failed to download video")
			if err := publishEvent(ctx, producer, TopicDownloadEvents, domain.DownloadFailed{
				EventHeader: domain.EventHeader{
					EventID:       generateEventID(),
					CorrelationID: envelope.Header.CorrelationID,
					CausationID:   envelope.Header.EventID,
					EventType:     "DownloadFailed",
					EventVersion:  "1.0",
					Producer:      "downloader-worker",
					OccurredAt:    time.Now().UTC(),
				},
				DownloadID: downloadReq.EventID,
				URL:        downloadReq.URL,
				Error:      err.Error(),
			}); err != nil {
				logger.WithError(err).Error("Failed to publish DownloadFailed event")
			}
			return err
		}

		// Upload to MinIO (R2)
		bucket := cfg.StorageBucket
		key := filepath.Join("downloads", filepath.Base(downloadedFile))
		_, err = minioClient.FPutObject(ctx, bucket, key, downloadedFile, minio.PutObjectOptions{ContentType: "video/mp4"})
		if err != nil {
			logger.WithError(err).Error("Failed to upload to MinIO")
			return err
		}

		// Get file info for the completed event
		fileInfo, err := os.Stat(downloadedFile)
		if err != nil {
			logger.WithError(err).Error("Failed to get file info")
			return err
		}

		// Publish DownloadCompleted event
		if err := publishEvent(ctx, producer, TopicDownloadEvents, domain.DownloadCompleted{
			EventHeader: domain.EventHeader{
				EventID:       generateEventID(),
				CorrelationID: envelope.Header.CorrelationID,
				CausationID:   envelope.Header.EventID,
				EventType:     "DownloadCompleted",
				EventVersion:  "1.0",
				Producer:      "downloader-worker",
				OccurredAt:    time.Now().UTC(),
			},
			DownloadID: downloadReq.EventID,
			URL:        downloadReq.URL,
			Bucket:     bucket,
			Key:        key,
			Format:     "mp4",
			FileSize:   fileInfo.Size(),
	PEG		
		)return fmt.Errorf("failed to publish DownloadCompleted event: %w", err) PEG		}

		// Publish to audit.events
		if err := publishEvent(ctx, producer, TopicAuditEvents, domain.DownloadCompleted{...}); err != nil {
			logger.WithError(err).Error("Failed to publish to audit.events")
		}

		return nil
		}

	// Create and start the consumer
	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID, TopicDownloadEvents, TopicDLQ, processFunc, logger)
	if err := consumer.Start(ctx); err != nil {
		logger.WithError(err).Fatal("Consumer failed")
	}
}

// downloadVideo uses yt-dlp to download the video
func downloadVideo(ctx context.Context, url string, logger *logrus.Logger) (string, error) {
	logger.Infof("Downloading video from URL: %s", url)

	// Use yt-dlp to download the video
	outputFile := filepath.Join(os.TempDir(), fmt.Sprintf("video_%d.mp4", time.Now().Unix()))
	cmd := exec.CommandContext(ctx, "yt-dlp", "-f", "best", "-o", outputFile, url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.WithError(err).Errorf("yt-dlp failed: %s", string(output))
		return "", fmt.Errorf("yt-dlp failed: %w, output: %s", err, string(output))
	}
	logger.Infof("Video downloaded to: %s", outputFile)
	return outputFile, nil
}

// publishEvent is a helper to publish an event to a specific topic
func publishEvent(ctx context.Context, producer *kafka.Producer, topic string, event interface{}) error {
	return producer.Publish(ctx, topic, event)
}

// generateEventID generates a simple unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// isAlreadyProcessed checks if the event has already been processed (idempotency)
func isAlreadyProcessed(ctx context.Context, eventID string) bool {
	// In a real implementation, check a database or cache (e.g., Redis) for the event ID.
	// For this example, we assume no events are pre-processed.
	return false
}
