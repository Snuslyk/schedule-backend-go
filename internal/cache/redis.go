package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(redisURL, password string) (*RedisClient, error) {
	var opts *redis.Options
	if strings.HasPrefix(redisURL, "redis://") || strings.HasPrefix(redisURL, "rediss://") {
		var err error
		opts, err = redis.ParseURL(redisURL)
		if err != nil {
			return nil, err
		}
	} else {
		opts = &redis.Options{
			Addr:     redisURL,
			Password: password,
			DB:       0,
		}
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return &RedisClient{client: client}, fmt.Errorf("redis ping warning: %w", err)
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) Client() *redis.Client {
	return r.client
}
