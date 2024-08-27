package user

import (
	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) ProduceCreate(info *model.UserInfo) error {
	return s.producer.Produce(s.kafkaServiceConfig.UserTopic(), info)
}
