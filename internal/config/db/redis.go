package db

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	redisEnvPrefix = "redis"
)

type redisConfig struct {
	host            string
	port            int
	username        string
	password        string
	db              int
	readTimeoutSec  int
	writeTimeoutSec int
}

func (cfg *redisConfig) setFromEnv() error {
	return envconfig.Process(redisEnvPrefix, cfg)
}

func NewRedisConfig() (*redisConfig, error) {
	cfg := &redisConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("redis config set from env err: %s", err)
	}
	logrus.Infof("redis config: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *redisConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}

func (cfg *redisConfig) Username() string {
	return cfg.username
}

func (cfg *redisConfig) Password() string {
	return cfg.password
}

func (cfg *redisConfig) DB() int {
	return cfg.db
}

func (cfg *redisConfig) ReadTimeoutSec() time.Duration {
	return time.Duration(cfg.readTimeoutSec) * time.Second
}

func (cfg *redisConfig) WriteTimeoutSec() time.Duration {
	return time.Duration(cfg.writeTimeoutSec) * time.Second
}
