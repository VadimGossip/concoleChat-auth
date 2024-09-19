package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	"github.com/pkg/errors"
)

const (
	swaggerHostEnvName = "APP_SWAGGER_HOST"
	swaggerPortEnvName = "APP_SWAGGER_PORT"
)

type swaggerConfig struct {
	host string
	port int
}

func (cfg *swaggerConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(swaggerHostEnvName)
	portStr := os.Getenv(swaggerPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("swaggerConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse swaggerConfig port")
	}
	return nil
}

func NewSwaggerConfig() (*swaggerConfig, error) {
	cfg := &swaggerConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("swaggerConfig set from env err: %s", err)
	}

	logger.Infof("swaggerConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *swaggerConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
