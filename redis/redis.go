package redis

import (
	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"time"
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
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:                 cfg.Addrs,
		DB:                    cfg.DB,
		MaxRetries:            maxRetries,
		MinRetryBackoff:       minRetryBackoff,
		MaxRetryBackoff:       maxRetryBackoff,
		DialTimeout:           dialTimeout,
		ReadTimeout:           readTimeout,
		WriteTimeout:          writeTimeout,
		ContextTimeoutEnabled: false,
		PoolFIFO:              false,
		PoolSize:              cfg.PoolSize,
		PoolTimeout:           poolTimeout,
		MinIdleConns:          minIdleConns,
		ConnMaxIdleTime:       idleTimeout,
	})
	if err := redisotel.InstrumentTracing(rdb, redisotel.WithDBStatement(false)); err != nil {
		return nil
	}
	return rdb
}
