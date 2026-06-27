package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Brokers []string
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg ProducerConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) PublishEvent(ctx context.Context, topic string, eventType string, payload interface{}) error {
	event := map[string]interface{}{
		"event_type": eventType,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"payload":    payload,
	}
	return p.Publish(ctx, topic, eventType, event)
}
