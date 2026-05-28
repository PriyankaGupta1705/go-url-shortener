package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore() (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{
		client: client,
		ttl:    24 * time.Hour, // cache URLs for 24 hours
	}, nil
}

func (r *RedisStore) Set(ctx context.Context, code, url string) error {
	return r.client.Set(ctx, "url:"+code, url, r.ttl).Err()
}

func (r *RedisStore) Get(ctx context.Context, code string) (string, bool) {
	val, err := r.client.Get(ctx, "url:"+code).Result()
	if err != nil {
		return "", false // cache miss
	}
	return val, true // cache hit
}

func (r *RedisStore) Delete(ctx context.Context, code string) {
	r.client.Del(ctx, "url:"+code)
}
