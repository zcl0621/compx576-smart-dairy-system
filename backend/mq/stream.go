package mq

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	redisdb "github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"go.uber.org/zap"
)

// Publish adds a message to the metrics stream
func Publish(fields map[string]string) error {
	args := &redis.XAddArgs{
		Stream: StreamName,
		Values: fields,
	}
	return redisdb.GetClient().XAdd(context.Background(), args).Err()
}

// StreamLen returns number of messages in the stream
func StreamLen() (int64, error) {
	return redisdb.GetClient().XLen(context.Background(), StreamName).Result()
}

// MessageHandler is called for each consumed message
// Return nil to ACK, return error to skip ACK (will redeliver)
type MessageHandler func(id string, values map[string]interface{}) error

// Consume blocks and reads messages from the stream using consumer group
// Runs until ctx is cancelled
func Consume(ctx context.Context, group, consumer string, handler MessageHandler) {
	client := redisdb.GetClient()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{StreamName, ">"},
			Count:    10,
			Block:    3 * time.Second,
		}).Result()

		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			if errors.Is(err, redis.Nil) {
				continue
			}
			// recreate stream and group if they got wiped
			if strings.Contains(err.Error(), "NOGROUP") {
				projectlog.L().Warn("stream or group missing, recreating", zap.Error(err))
				_ = Init()
				continue
			}
			projectlog.L().Error("xreadgroup failed", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				if handleErr := handler(msg.ID, msg.Values); handleErr != nil {
					projectlog.L().Error("handle message failed",
						zap.String("id", msg.ID),
						zap.Error(handleErr),
					)
					continue
				}

				if ackErr := client.XAck(ctx, StreamName, group, msg.ID).Err(); ackErr != nil {
					projectlog.L().Error("xack failed",
						zap.String("id", msg.ID),
						zap.Error(ackErr),
					)
				}
			}
		}
	}
}
