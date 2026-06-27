package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/twmb/franz-go/pkg/kgo"
)

type kafkaPublisher struct {
	client *kgo.Client
}

// NewEventPublisher creates a new Redpanda/Kafka event publisher using franz-go.
func NewEventPublisher(brokers []string) (ports.EventPublisher, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, err
	}
	return &kafkaPublisher{client: cl}, nil
}

func (p *kafkaPublisher) Publish(ctx context.Context, topic string, message []byte) error {
	if p.client == nil {
		return errors.New("kafka client not initialized")
	}
	p.client.Produce(ctx, &kgo.Record{
		Topic: topic,
		Value: message,
	}, nil)
	return nil
}

func (p *kafkaPublisher) Close() error {
	if p.client != nil {
		p.client.Close()
	}
	return nil
}
