package db

import (
	"context"
	"encoding/json"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
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

func (r *Redis) StoreToSortedList(ctx context.Context, key string, value *models.RedisSaveData) *redis.IntCmd {
	score := float64(value.TimeStamp)

	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	cmd := r.redisClient.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: jsonData,
	})

	if value.Expiration > 0 {
		r.redisClient.Expire(ctx, key, time.Duration(value.Expiration)*time.Second)
	}

	return cmd
}

func (r *Redis) GetSortedList(ctx context.Context, key string) ([]*models.RedisSaveData, error) {
	results, err := r.redisClient.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var items []*models.RedisSaveData
	for _, str := range results {
		var item models.RedisSaveData
		if err := json.Unmarshal([]byte(str), &item); err == nil {
			items = append(items, &item)
		}
	}
	return items, nil
}
