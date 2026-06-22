package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis() (*redis.Client, error) {
	host := os.Getenv("CACHE_HOST")
	if host == "" {
		host = "redis"
	}

	port := os.Getenv("CACHE_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("CACHE_PASSWORD")
	db := 0

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
