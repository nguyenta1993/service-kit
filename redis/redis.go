package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	maxRetries      = 3
	minRetryBackoff = 300 * time.Millisecond
	maxRetryBackoff = 500 * time.Millisecond
	dialTimeout     = 3 * time.Second
	readTimeout     = 3 * time.Second
	writeTimeout    = 3 * time.Second
	minIdleConns    = 20
	poolTimeout     = 6 * time.Second
	idleTimeout     = 12 * time.Second
)

func NewUniversalRedisClient(cfg Config) redis.UniversalClient {
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,
		DB:       cfg.DB,

		MaxRetries:      maxRetries,
		MinRetryBackoff: minRetryBackoff,
		MaxRetryBackoff: maxRetryBackoff,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    minIdleConns,

		PoolTimeout: poolTimeout,
		IdleTimeout: idleTimeout,
	})
}
