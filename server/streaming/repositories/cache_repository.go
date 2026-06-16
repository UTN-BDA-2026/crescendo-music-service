package repositories

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository struct {
	Client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{Client: client}
}

// GetAudioCache retrieves audio file from cache
func (cr *CacheRepository) GetAudioCache(ctx context.Context, fileID string) (io.ReadCloser, error) {
	data, err := cr.Client.Get(ctx, "audio:"+fileID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // not found
		}
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(data)), nil
}

// SetAudioCache stores audio file in cache with TTL
func (cr *CacheRepository) SetAudioCache(ctx context.Context, fileID string, data io.ReadCloser) error {
	defer data.Close()

	// Read all data from stream
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, data); err != nil {
		return err
	}

	// Cache for 1 hour by default
	ttl := 1 * time.Hour
	ttlStr := time.Now().Add(ttl).Sub(time.Now())

	return cr.Client.Set(ctx, "audio:"+fileID, buf.Bytes(), ttlStr).Err()
}

// SetStreamingState stores the current playback state
func (cr *CacheRepository) SetStreamingState(ctx context.Context, songID string, position int64) error {
	return cr.Client.Set(ctx, "playback:"+songID, position, 30*time.Minute).Err()
}

// GetStreamingState retrieves the current playback state
func (cr *CacheRepository) GetStreamingState(ctx context.Context, songID string) (int64, error) {
	val, err := cr.Client.Get(ctx, "playback:"+songID).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // not found
		}
		return 0, err
	}
	return val, nil
}

// ClearAudioCache removes audio from cache
func (cr *CacheRepository) ClearAudioCache(ctx context.Context, fileID string) error {
	return cr.Client.Del(ctx, "audio:"+fileID).Err()
}
