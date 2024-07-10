package model

type NetServerConfig struct {
	Port int
}

type Config struct {
	AppGrpcServer NetServerConfig
}
