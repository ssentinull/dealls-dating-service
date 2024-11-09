package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func (c *redisc) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.Do(ctx).Get(ctx, key)
}

func (c *redisc) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.Do(ctx).Set(ctx, key, value, expiration)
}

func (c *redisc) Delete(ctx context.Context, key string) *redis.IntCmd {
	return c.Do(ctx).Del(ctx, key)
}

func (c *redisc) Lock(ctx context.Context, key string, ttl time.Duration, opt *LockOptions) (*Lock, error) {
	if opt == nil {
		opt = new(LockOptions)
	}
	opt.Context = ctx
	return c.locker.Obtain(key, ttl, opt)
}
