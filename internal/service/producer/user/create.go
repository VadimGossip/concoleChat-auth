package user

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) ProduceCreate(info *model.UserInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v\n", err)
	}

	err = s.producer.Produce(&sarama.ProducerMessage{
		Topic: s.kafkaServiceConfig.UserTopic(),
		Value: sarama.StringEncoder(data),
	})

	if err != nil {
		return err
	}

	return nil
}
