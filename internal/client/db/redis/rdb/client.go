package rdb

import (
	"time"

	"github.com/redis/go-redis/v9"

	db "github.com/VadimGossip/concoleChat-auth/internal/client/db/redis"
)

type odbClient struct {
	masterDBC db.DB
}

type ClientOptions struct {
	Addr         string
	Username     string
	Password     string
	DB           int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New(options *ClientOptions) db.Client {
	dbc := redis.NewClient(&redis.Options{
		Addr:         options.Addr,
		Username:     options.Username,
		Password:     options.Password,
		DB:           options.DB,
		ReadTimeout:  options.ReadTimeout,
		WriteTimeout: options.WriteTimeout,
	})
	return &odbClient{
		masterDBC: NewDB(dbc),
	}
}

func (c *odbClient) DB() db.DB {
	return c.masterDBC
}

func (c *odbClient) Close() error {
	if c.masterDBC != nil {
		return c.masterDBC.Close()
	}

	return nil
}
