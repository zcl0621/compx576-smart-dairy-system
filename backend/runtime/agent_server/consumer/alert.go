package consumer

import (
	"context"
	"strconv"
	"time"

	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	alertservice "github.com/zcl0621/compx576-smart-dairy-system/service/alert"
	"go.uber.org/zap"
)

func StartAlertEvaluator(ctx context.Context) {
	projectlog.L().Info("alert evaluator consumer started")
	mq.Consume(ctx, mq.GroupAlertEvaluator, "evaluator_1", HandleAlertEval)
	projectlog.L().Info("alert evaluator consumer stopped")
}

func HandleAlertEval(id string, values map[string]any) error {
	cowID, _ := values["cow_id"].(string)
	metricType, _ := values["metric_type"].(string)
	metricValueStr, _ := values["metric_value"].(string)
	timestampStr, _ := values["timestamp"].(string)

	if cowID == "" || metricType == "" {
		return nil
	}

	switch model.MetricType(metricType) {
	case model.MetricTypeTemperature, model.MetricTypeHeartRate, model.MetricTypeBloodOxygen:
		if err := alertservice.EvaluateHealth(cowID, model.MetricType(metricType)); err != nil {
			projectlog.L().Error("health eval failed", zap.String("id", id), zap.Error(err))
			return err
		}

	case model.MetricTypeMilkAmount:
		currentValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			projectlog.L().Warn("skip milking eval: bad metric_value", zap.String("id", id))
			return nil
		}
		if err := alertservice.EvaluateMilking(cowID, currentValue, parseTimestamp(timestampStr)); err != nil {
			projectlog.L().Error("milking eval failed", zap.String("id", id), zap.Error(err))
			return err
		}

	case model.MetricTypeWeight:
		currentValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			projectlog.L().Warn("skip weight eval: bad metric_value", zap.String("id", id))
			return nil
		}
		if err := alertservice.EvaluateWeight(cowID, currentValue, parseTimestamp(timestampStr)); err != nil {
			projectlog.L().Error("weight eval failed", zap.String("id", id), zap.Error(err))
			return err
		}

	// latitude, longitude, milking_duration, device → not evaluated
	}

	return nil
}

func parseTimestamp(s string) time.Time {
	if s != "" && s != "0" {
		if ts, err := strconv.ParseInt(s, 10, 64); err == nil && ts > 0 {
			return time.Unix(ts, 0)
		}
	}
	return time.Now()
}
