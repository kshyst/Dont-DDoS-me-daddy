package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"main/internal/models"
)

type Redis struct {
	redisClient *redis.Client
}

func NewRedis() (Redis, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return Redis{redisClient}, nil
}

func (r *Redis) StoreToSortedList(ctx context.Context, key string, value models.RedisSaveData) *redis.IntCmd {
	return r.redisClient.ZAdd(ctx, key, redis.Z{})
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.redisClient.Get(ctx, key).Result()
}
