package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg ConsumerConfig) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     cfg.Brokers,
			Topic:       cfg.Topic,
			GroupID:     cfg.GroupID,
			StartOffset: kafka.FirstOffset,
			MaxWait:     500 * time.Millisecond,
		}),
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader difference
}

func (c *Consumer) UnmarshalNext(ctx context.Context, dest interface{}) (kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return msg, err
	}
	if err := json.Unmarshal(msg.Value, dest); err != nil {
		return msg, err
	}
	return msg, nil
}
