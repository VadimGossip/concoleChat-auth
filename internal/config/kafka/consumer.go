package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

func (cfg *kafkaConsumerConfig) setFromEnv() error {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return fmt.Errorf("kafkaConsumerConfig kafka brokers address not found")
	}

	cfg.brokers = strings.Split(brokersStr, ",")

	cfg.groupID = os.Getenv(groupIDEnvName)
	if len(cfg.groupID) == 0 {
		return fmt.Errorf("kafkaConsumerConfig kafk group id not found")
	}

	return nil
}

func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	cfg := &kafkaConsumerConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("kafkaConsumerConfig set from env err: %s", err)
	}

	logrus.Infof("kafkaConsumerConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *kafkaConsumerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
