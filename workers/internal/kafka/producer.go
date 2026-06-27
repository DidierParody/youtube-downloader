package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Producer is a generic Kafka producer for publishing events

type Producer struct {
	writer *kafka.Writer
	logger *logrus.Logger
}

// NewProducer creates a new Producer
func NewProducer(brokers []string, logger *logrus.Logger) *Producer {
	return &Producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      brokers,
			BatchTimeout: 10 * time.Millisecond,
		}),
		logger: logger,
	}
}

// Publish sends an event to a Kafka topic
func (p *Producer) Publish(ctx context.Context, topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Value: payload,
		Key:   []byte(""), // Optionally set a key for partitioning
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", topic, err)
	}

	p.logger.Debugf("Published event to topic %s", topic)
	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}
