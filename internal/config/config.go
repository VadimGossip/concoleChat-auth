package config

import (
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/spf13/viper"
)

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func unmarshal(cfg *model.Config) error {
	if err := viper.UnmarshalKey("app_grpc", &cfg.AppGrpcServer); err != nil {
		return err
	}
	return nil
}

func Init(configDir string) (*model.Config, error) {
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}
	cfg := &model.Config{}
	return cfg, unmarshal(cfg)
}
