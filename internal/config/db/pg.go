package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	pgHostEnvName     = "PG_HOST"
	pgPortEnvName     = "PG_PORT"
	pgUsernameEnvName = "PG_USERNAME"
	pgNameEnvName     = "PG_NAME"
	pgSSLModeEnvName  = "PG_SSLMODE"
	pgPasswordEnvName = "PG_PASSWORD"
)

type pgConfig struct {
	host     string
	port     int
	username string
	name     string
	sslMode  string
	password string
}

func (cfg *pgConfig) setFromEnv() error {
	var err error
	cfg.host = os.Getenv(pgHostEnvName)
	if len(cfg.host) == 0 {
		return fmt.Errorf("pgConfig host not found")
	}

	portStr := os.Getenv(pgPortEnvName)
	if len(portStr) == 0 {
		return fmt.Errorf("pgConfig port not found")
	}

	cfg.port, err = strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse pgConfig port")
	}

	cfg.username = os.Getenv(pgUsernameEnvName)
	if len(cfg.username) == 0 {
		return fmt.Errorf("pgConfig username not found")
	}

	cfg.name = os.Getenv(pgNameEnvName)
	if len(cfg.username) == 0 {
		return fmt.Errorf("pgConfig name not found")
	}

	cfg.sslMode = os.Getenv(pgSSLModeEnvName)
	if len(cfg.sslMode) == 0 {
		return fmt.Errorf("pgConfig ssl mode not found")
	}

	cfg.password = os.Getenv(pgPasswordEnvName)
	if len(cfg.password) == 0 {
		return fmt.Errorf("pgConfig password not found")
	}

	return nil
}

func NewPGConfig() (*pgConfig, error) {
	cfg := &pgConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("pgConfig set from env err: %s", err)
	}
	logrus.Infof("pgConfig: [%+v]", *cfg)

	return cfg, nil
}

func (cfg *pgConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", cfg.host, cfg.port, cfg.name, cfg.username, cfg.password, cfg.sslMode)
}
