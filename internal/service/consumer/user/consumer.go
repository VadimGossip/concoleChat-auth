package user

import (
	"context"
	"github.com/VadimGossip/concoleChat-auth/internal/client/kafka"
	defService "github.com/VadimGossip/concoleChat-auth/internal/service"
	def "github.com/VadimGossip/concoleChat-auth/internal/service/consumer"
)

var _ def.UserConsumerService = (*service)(nil)

type service struct {
	topic       string
	consumer    kafka.Consumer
	userService defService.UserService
}

func NewService(
	topic string,
	consumer kafka.Consumer,
	userService defService.UserService,
) *service {
	return &service{
		topic:       topic,
		consumer:    consumer,
		userService: userService,
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

		errChan <- s.consumer.Consume(ctx, s.topic, s.CreateUserHandler)
	}()

	return errChan
}
