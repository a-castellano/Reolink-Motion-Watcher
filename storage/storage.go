package storage

import (
	"context"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

type Storage struct {
	RedisClient *goredis.Client
	TTL         int32
}

func (storage Storage) UpdateTrigger(ctx context.Context, webcamName string) (bool, error) {
	err := storage.RedisClient.Set(ctx, webcamName, "triggered", time.Duration(storage.TTL)*time.Second).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
