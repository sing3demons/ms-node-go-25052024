package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type IRedis interface{}

type Cacher struct{ *redis.Client }

func New() *Cacher {
	uri := "localhost:6379"
	redisClient := redis.NewClient(&redis.Options{
		Addr: uri,
	})
	cmd, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("redis connected", cmd)

	return &Cacher{redisClient}
}
