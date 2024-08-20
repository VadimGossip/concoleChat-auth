package server

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

const (
	swaggerEnvPrefix = "app_swagger"
)

type swaggerConfig struct {
	host string
	port int
}

func (cfg *swaggerConfig) setFromEnv() error {
	return envconfig.Process(swaggerEnvPrefix, cfg)
}

func NewSwaggerConfig() (*swaggerConfig, error) {
	cfg := &swaggerConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("swagger server config set from env err: %s", err)
	}

	logrus.Infof("swagger server config: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *swaggerConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
