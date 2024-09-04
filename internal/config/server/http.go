package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/VadimGossip/concoleChat-auth/internal/logger"
	"github.com/pkg/errors"
)

const (
	httpHostEnvName = "APP_HTTP_HOST"
	httpPortEnvName = "APP_HTTP_PORT"
)

type httpConfig struct {
	host string
	port int
}

func (cfg *httpConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(httpHostEnvName)
	portStr := os.Getenv(httpPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("httpConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse httpConfig port")
	}
	return nil
}

func NewHTTPConfig() (*httpConfig, error) {
	cfg := &httpConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("httpConfig set from env err: %s", err)
	}

	logger.Infof("httpConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *httpConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
