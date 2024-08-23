package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	refreshTokenSecretKey  = "REFRESH_TOKEN_SECRET_KEY"
	accessTokenSecretKey   = "ACCESS_TOKEN_SECRET_KEY"
	refreshTokenExpiration = "REFRESH_TOKEN_EXPIRATION_MIN"
	accessTokenExpiration  = "ACCESS_TOKEN_EXPIRATION_MIN"
)

type tokenConfig struct {
	refreshSecretKey  string
	accessSecretKey   string
	refreshExpiration time.Duration
	accessExpiration  time.Duration
}

func (cfg *tokenConfig) setFromEnv() error {
	cfg.refreshSecretKey = os.Getenv(refreshTokenSecretKey)
	if len(cfg.refreshSecretKey) == 0 {
		return fmt.Errorf("tokenConfig refreshSecretKey not found")
	}

	cfg.accessSecretKey = os.Getenv(accessTokenSecretKey)
	if len(cfg.accessSecretKey) == 0 {
		return fmt.Errorf("tokenConfig accessSecretKey not found")
	}

	refreshExpirationStr := os.Getenv(refreshTokenExpiration)
	if len(refreshExpirationStr) == 0 {
		return fmt.Errorf("tokenConfig refreshExpiration not found")
	}

	refreshExpirationMin, err := strconv.ParseInt(refreshExpirationStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse tokenConfig refreshExpiration")
	}
	cfg.refreshExpiration = time.Duration(refreshExpirationMin) * time.Minute

	accessExpirationStr := os.Getenv(accessTokenExpiration)
	if len(accessExpirationStr) == 0 {
		return fmt.Errorf("tokenConfig accessExpiration not found")
	}

	accessExpirationMin, err := strconv.ParseInt(accessExpirationStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse tokenConfig accessExpiration")
	}
	cfg.accessExpiration = time.Duration(accessExpirationMin) * time.Minute
	return nil
}

func NewTokenConfig() (*tokenConfig, error) {
	cfg := &tokenConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("tokenConfig set from env err: %s", err)
	}

	logrus.Infof("tokenConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *tokenConfig) RefreshTokenSecretKey() string {
	return cfg.refreshSecretKey
}

func (cfg *tokenConfig) AccessTokenSecretKey() string {
	return cfg.accessSecretKey
}

func (cfg *tokenConfig) RefreshTokenExpiration() time.Duration {
	return cfg.refreshExpiration
}

func (cfg *tokenConfig) AccessTokenExpiration() time.Duration {
	return cfg.accessExpiration
}
