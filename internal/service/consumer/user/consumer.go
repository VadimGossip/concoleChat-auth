package user

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/client/kafka"
	"github.com/VadimGossip/concoleChat-auth/internal/config"
	defService "github.com/VadimGossip/concoleChat-auth/internal/service"
	def "github.com/VadimGossip/concoleChat-auth/internal/service/consumer"
)

var _ def.UserConsumerService = (*service)(nil)

type service struct {
	kafkaServiceConfig config.UserKafkaServiceConfig
	consumer           kafka.Consumer
	userService        defService.UserService
}

func NewService(
	kafkaServiceConfig config.UserKafkaServiceConfig,
	consumer kafka.Consumer,
	userService defService.UserService,
) *service {
	return &service{
		kafkaServiceConfig: kafkaServiceConfig,
		consumer:           consumer,
		userService:        userService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, s.kafkaServiceConfig.UserTopic(), s.CreateUserHandler)
	}()

	return errChan
}
