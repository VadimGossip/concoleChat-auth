package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	"github.com/pkg/errors"
)

const (
	prometheusHostEnvName = "APP_PROMETHEUS_HOST"
	prometheusPortEnvName = "APP_PROMETHEUS_PORT"
)

type prometheusConfig struct {
	host string
	port int
}

func (cfg *prometheusConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(prometheusHostEnvName)
	portStr := os.Getenv(prometheusPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("prometheusConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse prometheusConfig port")
	}
	return nil
}

func NewPrometheusConfig() (*prometheusConfig, error) {
	cfg := &prometheusConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("prometheusConfig set from env err: %s", err)
	}

	logger.Infof("prometheusConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *prometheusConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
