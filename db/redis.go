package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedis(redisURI string) *redis.Client {
	ctx := context.Background()

	redisClient := redis.NewClient(
		&redis.Options{
			Addr: redisURI,
		},
	)

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err.Error())
	}

	return redisClient
}
