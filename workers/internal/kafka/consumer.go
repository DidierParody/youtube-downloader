package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// EventHeader is shared across all events to ensure idempotency and traceability
type EventHeader struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	CausationID     string    `json:"causation_id"`
	EventType     string    `json:"event_type"`
	EventVersion  string    `json:"event_version"`
	Producer      string    `json:"producer"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// EventInterface is the minimal contract for events
func (h *EventHeader) Topic(topic string) string { return topic }
func (h *EventHeader) Type() string              { return h.EventType }

// Generic event envelope for dispatching

type EventEnvelope struct {
	Header  EventHeader       `json:"header"`
	Payload json.RawMessage   `json:"payload"`
}

// Consumer is a generic Kafka consumer with consumer groups and graceful shutdown

type Consumer struct {
	reader        *kafka.Reader
	logger        *logrus.Logger
	groupID       string
	topic         string
	dlqTopic      string
	processFunc   func(ctx context.Context, msg kafka.Message) error
	dlqWriter     *kafka.Writer
}

// NewConsumer creates a new Consumer
func NewConsumer(brokers []string, groupID, topic, dlqTopic string, processFunc func(ctx context.Context, msg kafka.Message) error, logger *logrus.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MaxWait:  1 * time.Second,
			Logger:   kafka.LoggerFunc(func(_ string, _ ...interface{}) {}),
		}),
		logger:      logger,
		groupID:     groupID,
		topic:       topic,
		dlqTopic:    dlqTopic,
		processFunc: processFunc,
		dlqWriter:   kafka.NewWriter(kafka.WriterConfig{Brokers: brokers, Topic: dlqTopic}),
	}
}

// Start begins consuming messages with graceful shutdown on SIGTERM/SIGINT
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Infof("Starting consumer for topic: %s, group: %s", c.topic, c.groupID)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Listen for shutdown signals
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigch:
			c.logger.Info("Shutdown signal received, stopping consumer...")
			cancel()
		case <-ctx.Done():
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return c.reader.Close()
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return c.reader.Close()
			}
			c.logger.WithError(err).Error("Failed to fetch message")
			continue
		}

		// Process message
		if err := c.processFunc(ctx, msg); err != nil {
			c.logger.WithError(err).Error("Message processing failed, sending to DLQ")
			if derr := c.sendToDLQ(ctx, msg); derr != nil {
				c.logger.WithError(derr).Error("Failed to send message to DLQ")
			}
			// Commit the offset even on failure to avoid infinite loops
			if cerr := c.reader.CommitMessages(ctx, msg); cerr != nil {
				c.logger.WithError(cerr).Error("Failed to commit message after DLQ")
			}
			continue
		}

		// Commit the offset after successful processing
		if cerr := c.reader.CommitMessages(ctx, msg); cerr != nil {
			c.logger.WithError(cerr).Error("Failed to commit message after successful processing")
			continue
		}

		c.logger.Debugf("Processed and committed message: %s", string(msg.Value))
	}
}

// sendToDLQ sends the failed message to the DLQ topic
func (c *Consumer) sendToDLQ(ctx context.Context, msg kafka.Message) error {
	err := c.dlqWriter.WriteMessages(ctx, msg)
	return err
}

// Close closes the consumer
func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.dlqWriter.Close()
}
