package server

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
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
	if len(cfg.host) == 0 {
		return fmt.Errorf("httpConfig host not found")
	}

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

	logrus.Infof("httpConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *httpConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
