package model

import "time"

type NetServerCfg struct {
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

type RedisDbCfg struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Db              int
	ReadTimeoutSec  int
	WriteTimeoutSec int
}

type UserCacheCfg struct {
	Expire time.Duration
}

type Config struct {
	AppGrpcServer NetServerCfg
	PgDb          PGDbCfg
	RedisDb       RedisDbCfg
	UserCache     UserCacheCfg
}
