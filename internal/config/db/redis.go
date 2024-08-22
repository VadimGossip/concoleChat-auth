package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	redisHostEnvName         = "REDIS_HOST"
	redisPortEnvName         = "REDIS_PORT"
	redisUsernameEnvName     = "REDIS_USERNAME"
	redisPasswordEnvName     = "REDIS_PASSWORD"
	redisDBNameEnvName       = "REDIS_DB"
	redisReadTimeoutEnvName  = "REDIS_READ_TIMEOUT_SEC"
	redisWriteTimeoutEnvName = "REDIS_WRITE_TIMEOUT_SEC"
)

type redisConfig struct {
	host         string
	port         int
	username     string
	password     string
	db           int
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (cfg *redisConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(redisHostEnvName)
	if len(cfg.host) == 0 {
		return fmt.Errorf("redisConfig host not found")
	}

	portStr := os.Getenv(redisPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("redisConfig port not found")
	}
	cfg.username = os.Getenv(redisUsernameEnvName)
	cfg.password = os.Getenv(redisPasswordEnvName)

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse redisConfig port")
	}

	dbStr := os.Getenv(redisDBNameEnvName)
	if len(dbStr) == 0 {
		return fmt.Errorf("redisConfig db not found")
	}

	cfg.db, err = strconv.Atoi(dbStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse redisConfig db")
	}

	readTimeoutStr := os.Getenv(redisReadTimeoutEnvName)
	if len(readTimeoutStr) == 0 {
		return fmt.Errorf("redisConfig read timeout not found")
	}

	readTimeoutSec, err := strconv.ParseInt(readTimeoutStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse redisConfig read timeout")
	}
	cfg.readTimeout = time.Duration(readTimeoutSec) * time.Second

	writeTimeoutStr := os.Getenv(redisWriteTimeoutEnvName)
	if len(writeTimeoutStr) == 0 {
		return fmt.Errorf("redis write timeout not found")
	}

	writeTimeoutSec, err := strconv.ParseInt(writeTimeoutStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse redisConfig write timeout")
	}
	cfg.writeTimeout = time.Duration(writeTimeoutSec) * time.Second

	return nil
}

func NewRedisConfig() (*redisConfig, error) {
	cfg := &redisConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("redisConfig set from env err: %s", err)
	}
	logrus.Infof("redisConfig: [%+v]", *cfg)
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
	return cfg.readTimeout
}

func (cfg *redisConfig) WriteTimeoutSec() time.Duration {
	return cfg.writeTimeout
}
