package server

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

const (
	httpEnvPrefix = "app_http"
)

type httpConfig struct {
	host string
	port int
}

func (cfg *httpConfig) setFromEnv() error {
	return envconfig.Process(httpEnvPrefix, cfg)
}

func NewHTTPConfig() (*httpConfig, error) {
	cfg := &httpConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("http server config set from env err: %s", err)
	}

	logrus.Infof("http server config: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *httpConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
