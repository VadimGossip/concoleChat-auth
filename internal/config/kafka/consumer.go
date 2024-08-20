package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const (
	consumerBrokersEnvName = "KAFKA_BROKERS"
	consumerGroupIDEnvName = "KAFKA_GROUP_ID"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

func (cfg *kafkaConsumerConfig) setFromEnv() error {
	brokersStr := os.Getenv(consumerBrokersEnvName)
	if len(brokersStr) == 0 {
		return fmt.Errorf("kafkaConsumerConfig kafka brokers address not found")
	}

	cfg.brokers = strings.Split(brokersStr, ",")

	cfg.groupID = os.Getenv(consumerGroupIDEnvName)
	if len(cfg.groupID) == 0 {
		return fmt.Errorf("kafkaConsumerConfig kafka group id not found")
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
