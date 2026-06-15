package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCacheConnection() *Cache {
	host := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	password := os.Getenv("CACHE_PASSWORD")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           0,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("cache unavailable: %v", err))
	}

	return &Cache{client: rdb}
}

func (c *Cache) IsReady() bool {
	return c != nil && c.client != nil
}

func (c *Cache) Get(key string) (string, bool, error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	if c == nil || c.client == nil {
		return nil // cache apagado, no rompe el flujo
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("cache set failed: %w", err)
	}

	return nil
}
