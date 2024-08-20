package user

import (
	"github.com/VadimGossip/concoleChat-auth/internal/client/kafka"
	"github.com/VadimGossip/concoleChat-auth/internal/config"
)

type service struct {
	kafkaServiceConfig config.UserKafkaServiceConfig
	producer           kafka.Producer
}

func NewService(
	kafkaServiceConfig config.UserKafkaServiceConfig,
	producer kafka.Producer,
) *service {
	return &service{
		kafkaServiceConfig: kafkaServiceConfig,
		producer:           producer,
	}
}
