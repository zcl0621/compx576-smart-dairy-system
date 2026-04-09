package mq

import (
	"context"
	"strings"

	redisdb "github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"go.uber.org/zap"
)

const StreamName = "metrics_stream"
const GroupMetricWriter = "metric_writer"

func Init() error {
	client := redisdb.GetClient()
	ctx := context.Background()

	// create consumer group, MKSTREAM creates the stream if missing
	err := client.XGroupCreateMkStream(ctx, StreamName, GroupMetricWriter, "0").Err()
	if err != nil && !strings.HasPrefix(err.Error(), "BUSYGROUP") {
		return err
	}

	projectlog.L().Info("mq init done", zap.String("stream", StreamName))
	return nil
}
