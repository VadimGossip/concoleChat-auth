package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	"github.com/pkg/errors"
)

const (
	userCacheExpireSec = "USER_CACHE_EXPIRE_SEC"
)

type userCacheConfig struct {
	expire time.Duration
}

func (cfg *userCacheConfig) setFromEnv() error {
	expireStr := os.Getenv(userCacheExpireSec)
	if len(expireStr) == 0 {
		return fmt.Errorf("userCacheConfig expire not found")
	}

	expireSec, err := strconv.ParseInt(expireStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse userCacheConfig expire")
	}
	cfg.expire = time.Duration(expireSec) * time.Second
	return nil
}

func NewUserCacheConfig() (*userCacheConfig, error) {
	cfg := &userCacheConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("userCacheConfig set from env err: %s", err)
	}

	logger.Infof("userCacheConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *userCacheConfig) Expire() time.Duration {
	return cfg.expire
}
