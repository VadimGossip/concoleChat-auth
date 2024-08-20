package kafka

import (
	"context"
	"github.com/IBM/sarama"

	"github.com/VadimGossip/concoleChat-auth/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) error
	Close() error
}

type Producer interface {
	Produce(ctx context.Context, msg *sarama.ProducerMessage) error
	Close() error
}
