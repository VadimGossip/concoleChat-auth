package kafka

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) error
	Close() error
}

type Producer interface {
	Produce(topic string, msg any) error
	Close() error
}
