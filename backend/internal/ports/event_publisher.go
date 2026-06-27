package ports

import "context"

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Close() error
}
