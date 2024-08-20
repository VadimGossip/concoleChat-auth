package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

const (
	consumerEnvPrefix = "kafka_consumer"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

func (cfg *kafkaConsumerConfig) setFromEnv() error {
	return envconfig.Process(consumerEnvPrefix, cfg)
}

func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	cfg := &kafkaConsumerConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("kafka consumer config set from env err: %s", err)
	}

	logrus.Infof("kafka consumer config: [%+v]", *cfg)
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
