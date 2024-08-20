package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	grpcHostEnvName = "APP_GRPC_HOST"
	grpcPortEnvName = "APP_GRPC_PORT"
)

type grpcConfig struct {
	host string
	port int
}

func (cfg *grpcConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(grpcHostEnvName)
	portStr := os.Getenv(grpcPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("grpcConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse grpcConfig port")
	}
	return nil
}

func NewGRPCConfig() (*grpcConfig, error) {
	cfg := &grpcConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("grpcConfig set from env err: %s", err)
	}

	logrus.Infof("grpcConfig: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *grpcConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}
