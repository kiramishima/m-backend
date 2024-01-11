package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"kiramishima/m-backend/internal/core/domain"
	"time"
)

// RedisCache struct
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Get(key string, value interface{}) error {
	data, err := c.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("cache miss for key %q", key)
	} else if err != nil {
		return fmt.Errorf("failed to get value for key %q: %v", key, err)
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return fmt.Errorf("failed to unmarshal cache value for key %q: %v", key, err)
	}

	return nil
}

func (c *RedisCache) Set(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %q: %v", key, err)
	}

	if err := c.client.Set(context.Background(), key, data, duration).Err(); err != nil {
		return fmt.Errorf("failed to set value for key %q: %v", key, err)
	}

	return nil
}

func (c *RedisCache) Delete(key string) error {
	if err := c.client.Del(context.Background(), key).Err(); err != nil {
		return fmt.Errorf("failed to delete value for key %q: %v", key, err)
	}
	return nil
}

var Module = fx.Module("cache",
	fx.Provide(func(cfg *domain.Configuration) *RedisCache {
		cache, _ := NewRedisCache(cfg.Addr, cfg.Password)
		return cache
	}),
)
