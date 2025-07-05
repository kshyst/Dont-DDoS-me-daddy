package db

import (
	"context"
	"encoding/json"
	"fmt"
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

func CreateRedis(redisClient *redis.Client) Redis {
	return Redis{redisClient}
}

func (r *Redis) StoreToSortedList(ctx context.Context, key string, value *models.RedisSaveData) (*redis.IntCmd, error) {
	// checking if the given key is a sorted list key
	keyType, err := r.redisClient.Type(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check key type: %w", err)
	}

	if keyType != "none" && keyType != "zset" {
		// Delete the key if it's not a sorted set
		if _, err := r.redisClient.Del(ctx, key).Result(); err != nil {
			return nil, fmt.Errorf("failed to delete existing key: %w", err)
		}
	}

	// put the current time as the score of the sortedlist item
	score := float64(value.TimeStamp)

	// marshal all needed data to json to keep as the member
	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	cmd := r.redisClient.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: jsonData,
	})

	if value.Expiration > 0 {
		r.redisClient.Expire(ctx, key, time.Duration(value.Expiration)*time.Second)
	}

	return cmd, nil
}

func (r *Redis) GetSortedList(ctx context.Context, key string) ([]*models.RedisSaveData, error) {
	results, err := r.redisClient.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get data using zrange: %w", err)
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
