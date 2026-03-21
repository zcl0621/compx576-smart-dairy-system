package redis

import (
	"context"
	"fmt"
	"sync"

	redislibrary "github.com/redis/go-redis/v9"
	"github.com/zcl0621/compx576-smart-dairy-system/config"
)

var (
	client *redislibrary.Client
	once   sync.Once
)

func InitRedis() error {
	var initErr error
	once.Do(func() {
		cfg := config.Get().Redis

		client = redislibrary.NewClient(&redislibrary.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		initErr = client.Ping(context.Background()).Err()
	})

	return initErr
}

func GetClient() *redislibrary.Client {
	return client
}
