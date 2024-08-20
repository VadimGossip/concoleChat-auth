package server

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
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
	if len(cfg.host) == 0 {
		return fmt.Errorf("swaggerConfig host not found")
	}

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

	logrus.Infof("swaggerConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *swaggerConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
