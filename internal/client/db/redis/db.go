package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client interface {
	DB() DB
	Close() error
}

type Execer interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	GetWithDur(ctx context.Context, keys string) (string, time.Duration, error)
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
	MGetWithDur(ctx context.Context, keys ...string) ([]any, time.Duration, error)
	HGetAll(ctx context.Context, key string, dest interface{}) error
	HSet(ctx context.Context, key string, values interface{}, expire time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	Execer
	Pinger
	Close() error
}
