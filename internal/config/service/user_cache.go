package service

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	userCacheEnvPrefix = "user_cache"
)

type userCacheConfig struct {
	expireSec int64 `envconfig:"EXPIRE_SEC"`
}

func (cfg *userCacheConfig) setFromEnv() error {
	return envconfig.Process(userCacheEnvPrefix, cfg)
}

func NewUserCacheConfig() (*userCacheConfig, error) {
	cfg := &userCacheConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("user cache config set from env err: %s", err)
	}

	logrus.Infof("user cache config: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *userCacheConfig) Expire() time.Duration {
	return time.Duration(cfg.expireSec) * time.Second
}
