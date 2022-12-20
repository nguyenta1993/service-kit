package cache

import (
	"context"
	"github.com/go-redis/cache/v9"
	"github.com/go-redis/redis/v9"
	"github.com/gogovan/ggx-kr-service-utils/tracing"
	"reflect"
	"time"
)

const (
	defaultCacheSize = 1000
	defaultCacheTTL  = 5 * time.Minute
)

type CacheInterface[T any] interface {
	Get(ctx context.Context, key string) (*T, error)
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type cacheImpl[T any] struct {
	cache *cache.Cache
}

func (c cacheImpl[T]) Del(ctx context.Context, key string) (err error) {
	ctx, span := tracing.StartSpanFromContext(ctx, "cacheImpl[T].Del"+key)
	defer func() {
		span.RecordError(err)
		span.End()
	}()
	return c.cache.Delete(ctx, key)
}

func (c cacheImpl[T]) Get(ctx context.Context, key string) (val *T, err error) {
	ctx, span := tracing.StartSpanFromContext(ctx, "cacheImpl.Get"+key+" "+reflect.TypeOf(val).String())
	defer func() {
		span.RecordError(err)
		span.End()
	}()
	var value T
	if err := c.cache.Get(ctx, key, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (c cacheImpl[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) (err error) {
	ctx, span := tracing.StartSpanFromContext(ctx, "cacheImpl.Set"+key+" "+reflect.TypeOf(value).String())
	defer func() {
		span.RecordError(err)
		span.End()
	}()
	item := &cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
		Do:    nil,
	}
	if err := c.cache.Set(item); err != nil {
		return err
	}
	return nil
}

func NewCache[T any](rdb redis.UniversalClient, size int, ttl time.Duration) CacheInterface[T] {
	if size == 0 {
		size = defaultCacheSize
	}
	if ttl == 0 {
		ttl = defaultCacheTTL
	}
	return &cacheImpl[T]{
		cache: cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(size, ttl),
		})}
}

func NewInmemoryCache[T any](size int, ttl time.Duration) CacheInterface[T] {
	return NewCache[T](nil, size, ttl)
}
