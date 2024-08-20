package db

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

const (
	pgEnvPrefix = "pg"
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
	return envconfig.Process(pgEnvPrefix, cfg)
}

func NewPGConfig() (*pgConfig, error) {
	cfg := &pgConfig{}
	if err := cfg.setFromEnv(); err != nil {
		return nil, fmt.Errorf("posgress config set from env err: %s", err)
	}
	logrus.Infof("posgress config: [%+v]", *cfg)
	return cfg, nil
}

func (cfg *pgConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", cfg.host, cfg.port, cfg.name, cfg.username, cfg.password, cfg.sslMode)
}
