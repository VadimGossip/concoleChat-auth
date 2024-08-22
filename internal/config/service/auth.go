package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

const (
	refreshTokenSecretKey  = "REFRESH_TOKEN_SECRET_KEY"
	accessTokenSecretKey   = "ACCESS_TOKEN_SECRET_KEY"
	refreshTokenExpiration = "REFRESH_TOKEN_EXPIRATION_MIN"
	accessTokenExpiration  = "ACCESS_TOKEN_EXPIRATION_MIN"
)

type authServiceConfig struct {
	refreshSecretKey  string
	accessSecretKey   string
	refreshExpiration time.Duration
	accessExpiration  time.Duration
}

func (cfg *authServiceConfig) setFromEnv() error {
	cfg.refreshSecretKey = os.Getenv(refreshTokenSecretKey)
	if len(cfg.refreshSecretKey) == 0 {
		return fmt.Errorf("authServiceConfig refreshSecretKey not found")
	}

	cfg.accessSecretKey = os.Getenv(accessTokenSecretKey)
	if len(cfg.accessSecretKey) == 0 {
		return fmt.Errorf("authServiceConfig accessSecretKey not found")
	}

	refreshExpirationStr := os.Getenv(refreshTokenExpiration)
	if len(refreshExpirationStr) == 0 {
		return fmt.Errorf("authServiceConfig refreshExpiration not found")
	}

	refreshExpirationMin, err := strconv.ParseInt(refreshExpirationStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse authServiceConfig refreshExpiration")
	}
	cfg.refreshExpiration = time.Duration(refreshExpirationMin) * time.Second

	accessExpirationStr := os.Getenv(accessTokenExpiration)
	if len(accessExpirationStr) == 0 {
		return fmt.Errorf("authServiceConfig accessExpiration not found")
	}

	accessExpirationMin, err := strconv.ParseInt(accessExpirationStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse authServiceConfig accessExpiration")
	}
	cfg.accessExpiration = time.Duration(accessExpirationMin) * time.Second
	return nil
}

func NewAuthServiceConfig() (*authServiceConfig, error) {
	cfg := &authServiceConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("authServiceConfig set from env err: %s", err)
	}

	logrus.Infof("authServiceConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *authServiceConfig) RefreshTokenSecretKey() string {
	return cfg.refreshSecretKey
}

func (cfg *authServiceConfig) AccessTokenSecretKey() string {
	return cfg.accessSecretKey
}

func (cfg *authServiceConfig) RefreshTokenExpiration() time.Duration {
	return cfg.refreshExpiration
}

func (cfg *authServiceConfig) AccessTokenExpiration() time.Duration {
	return cfg.accessExpiration
}
