package repositories

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository struct {
	Client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{Client: client}
}

func (cr *CacheRepository) SetStreamingState(ctx context.Context, songID string, position int64) error {
	return cr.Client.Set(ctx, "playback:"+songID, position, 30*time.Minute).Err()
}

func (cr *CacheRepository) GetStreamingState(ctx context.Context, songID string) (int64, error) {
	val, err := cr.Client.Get(ctx, "playback:"+songID).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}
