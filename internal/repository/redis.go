package repository

import (
	"context"
	"fmt"
	"github.com/SiriusServiceDesk/auth-service/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepositoryImpl struct {
	client *redis.Client
}

type RedisRepository interface {
	Get(key string) (string, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

func NewRedisClient(cfg *config.Config) RedisRepository {
	client := redis.NewClient(&redis.Options{
		Password: cfg.Redis.Password,
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		DB:       0,
	})
	return RedisRepositoryImpl{client: client}
}

func (rr RedisRepositoryImpl) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := rr.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rr RedisRepositoryImpl) Set(key string, value interface{}) error {
	ctx := context.Background()
	if err := rr.client.Set(ctx, key, value, time.Minute*10); err != nil {
		return err.Err()
	}
	return nil
}

func (rr RedisRepositoryImpl) Delete(key string) error {
	ctx := context.Background()
	if err := rr.client.Del(ctx, key); err != nil {
		return err.Err()
	}
	return nil
}
