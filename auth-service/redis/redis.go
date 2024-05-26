package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	Ping(ctx context.Context) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (int64, error)
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Close() error
}

type cacher struct{ *redis.Client }

func New() IRedis {
	uri := "localhost:6379"
	redisClient := redis.NewClient(&redis.Options{
		Addr: uri,
	})
	cmd, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("redis connected", cmd)

	return &cacher{redisClient}
}

func (c *cacher) Close() error {
	return c.Client.Close()
}

func (c *cacher) Ping(ctx context.Context) (string, error) {
	return c.Client.Ping(ctx).Result()
}

func (c *cacher) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *cacher) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *cacher) Del(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *cacher) Exists(ctx context.Context, key string) (int64, error) {
	return c.Client.Exists(ctx, key).Result()
}

func (c *cacher) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.SetEx(ctx, key, value, expiration).Err()
}
