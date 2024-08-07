package model

import "time"

type NetServerConfig struct {
	Port int
}

type PGDbCfg struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

type RedisDbConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Db           int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Config struct {
	AppGrpcServer NetServerConfig
	PgDb          PGDbCfg
	RedisDb       RedisDbConfig
}
