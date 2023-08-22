package storage

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Storage struct {
	RedisClient *goredis.Client
	TTL         uint8
}

func (storage Storage) UpdateTrigger(ctx context.Context, webcamName string) (bool, error) {
	err := storage.RedisClient.Set(ctx, webcamName, "triggered", time.Duration(storage.TTL)*time.Second).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (storage Storage) CheckTrigger(ctx context.Context, webcamName string) (bool, error) {
	_, err := storage.RedisClient.Get(ctx, webcamName).Result()
	if err == goredis.Nil { // Motion has not been triggered
		return false, nil
	} else {
		if err != nil {
			return false, err
		} else { // Key exists
			return true, nil
		}
	}
}
