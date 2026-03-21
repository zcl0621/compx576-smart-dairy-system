package redis

import (
	"context"
	"errors"
	"time"

	redislibrary "github.com/redis/go-redis/v9"
)

func Set(key string, value any, expiration time.Duration) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	return client.Set(context.Background(), key, value, expiration).Err()
}

func Get(key string) (string, error) {
	client, err := getClient()
	if err != nil {
		return "", err
	}

	return client.Get(context.Background(), key).Result()
}

func Del(keys ...string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	return client.Del(context.Background(), keys...).Err()
}

func Incr(key string) (int64, error) {
	client, err := getClient()
	if err != nil {
		return 0, err
	}

	return client.Incr(context.Background(), key).Result()
}

func Expire(key string, expiration time.Duration) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	return client.Expire(context.Background(), key, expiration).Err()
}

func Exists(key string) (bool, error) {
	client, err := getClient()
	if err != nil {
		return false, err
	}

	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func FlushDB() error {
	client, err := getClient()
	if err != nil {
		return err
	}

	return client.FlushDB(context.Background()).Err()
}

func getClient() (*redislibrary.Client, error) {
	if client == nil {
		return nil, errors.New("redis isn't ready")
	}

	return client, nil
}
