package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/VadimGossip/concoleChat-auth/internal/logger"
)

const (
	producerBrokersEnvName = "KAFKA_BROKERS"
)

type kafkaProducerConfig struct {
	brokers []string
}

func (cfg *kafkaProducerConfig) setFromEnv() error {
	brokersStr := os.Getenv(producerBrokersEnvName)
	if len(brokersStr) == 0 {
		return fmt.Errorf("kafkaProducerConfig kafka brokers address not found")
	}

	cfg.brokers = strings.Split(brokersStr, ",")

	return nil
}

func NewKafkaProducerConfig() (*kafkaProducerConfig, error) {
	cfg := &kafkaProducerConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("kafkaProducerConfig set from env err: %s", err)
	}

	logger.Infof("kafkaProducerConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *kafkaProducerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	return config
}
