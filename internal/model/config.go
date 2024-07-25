package model

type NetServerConfig struct {
	Port int
}

type DbCfg struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}
type Config struct {
	AppGrpcServer NetServerConfig
	Db            DbCfg
}
