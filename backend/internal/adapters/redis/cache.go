package redis

import (
	"context"
	"time"

	"github.com/DidierParody/youtube-downloader/backend/internal/ports"
	"github.com/redis/go-redis/v9"
)

type cache struct {
	client *redis.Client
}

// NewCache creates a new Redis-backed cache.
func NewCache(addr, password string, db int) (ports.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &cache{client: client}, nil
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *cache) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return c.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
