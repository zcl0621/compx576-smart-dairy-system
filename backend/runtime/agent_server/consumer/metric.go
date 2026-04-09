package consumer

import (
	"context"
	"strconv"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	"go.uber.org/zap"
)

func StartMetricWriter(ctx context.Context) {
	projectlog.L().Info("metric writer consumer started")
	mq.Consume(ctx, mq.GroupMetricWriter, "writer_1", handleMetric)
	projectlog.L().Info("metric writer consumer stopped")
}

func handleMetric(id string, values map[string]interface{}) error {
	cowID, _ := values["cow_id"].(string)
	source, _ := values["source"].(string)
	metricType, _ := values["metric_type"].(string)
	metricValueStr, _ := values["metric_value"].(string)
	unit, _ := values["unit"].(string)
	timestampStr, _ := values["timestamp"].(string)

	if cowID == "" || source == "" || metricType == "" || metricValueStr == "" || unit == "" {
		projectlog.L().Error("skip message with missing fields", zap.String("id", id))
		return nil // ACK to skip bad messages
	}

	metricValue, err := strconv.ParseFloat(metricValueStr, 64)
	if err != nil {
		projectlog.L().Error("bad metric_value", zap.String("id", id), zap.String("value", metricValueStr))
		return nil
	}

	m := &model.Metric{
		CowID:       cowID,
		Source:      model.MetricSource(source),
		MetricType:  model.MetricType(metricType),
		MetricValue: metricValue,
		Unit:        model.MetricUnit(unit),
	}

	// override created_at with device timestamp if present
	if timestampStr != "" {
		ts, err := strconv.ParseInt(timestampStr, 10, 64)
		if err == nil && ts > 0 {
			m.CreatedAt = time.Unix(ts, 0)
		}
	}

	if err := pg.DB.Create(m).Error; err != nil {
		projectlog.L().Error("write metric failed", zap.String("id", id), zap.Error(err))
		return err // don't ACK, will retry
	}

	return nil
}
