package config

import (
	"time"

	"github.com/IBM/sarama"
)

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type SwaggerConfig interface {
	Address() string
}

type PrometheusConfig interface {
	Address() string
}

type PgConfig interface {
	DSN() string
}

type RedisConfig interface {
	Address() string
	Username() string
	Password() string
	DB() int
	ReadTimeoutSec() time.Duration
	WriteTimeoutSec() time.Duration
}

type TokenConfig interface {
	RefreshTokenSecretKey() string
	AccessTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type UserCacheConfig interface {
	Expire() time.Duration
}

type UserKafkaServiceConfig interface {
	UserTopic() string
}

type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

type KafkaProducerConfig interface {
	Brokers() []string
	Config() *sarama.Config
}
