package server

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

const (
	grpcEnvPrefix = "app_grpc"
)

type grpcConfig struct {
	host string
	port int
}

func (cfg *grpcConfig) setFromEnv() error {
	return envconfig.Process(grpcEnvPrefix, cfg)
}

func NewGRPCConfig() (*grpcConfig, error) {
	cfg := &grpcConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("grpc server config set from env err: %s", err)
	}

	logrus.Infof("grpc server config: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *grpcConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
