package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	userKafkaTopic = "USER_KAFKA_TOPIC"
)

type userKafkaServiceConfig struct {
	userTopic string
}

func (cfg *userKafkaServiceConfig) setFromEnv() error {
	cfg.userTopic = os.Getenv(userKafkaTopic)
	if len(cfg.userTopic) == 0 {
		return fmt.Errorf("userKafkaServiceConfig topic not found")
	}
	return nil
}

func NewUserKafkaServiceConfig() (*userKafkaServiceConfig, error) {
	cfg := &userKafkaServiceConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("userKafkaServiceConfig set from env err: %s", err)
	}

	logrus.Infof("userKafkaServiceConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *userKafkaServiceConfig) UserTopic() string {
	return cfg.userTopic
}
