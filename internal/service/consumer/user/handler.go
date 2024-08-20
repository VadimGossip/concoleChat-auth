package user

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) CreateUserHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	info := &model.UserInfo{}

	if err := json.Unmarshal(msg.Value, info); err != nil {
		return err
	}

	_, err := s.userService.Create(ctx, info)
	if err != nil {
		return err
	}

	return nil
}
