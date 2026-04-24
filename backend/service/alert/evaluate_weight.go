package alert

import (
	"fmt"
	"math"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

func EvaluateWeight(cowID string, currentValue float64, currentTimestamp time.Time) error {
	avg, count, err := sevenDayAverage(cowID, model.MetricTypeWeight, currentTimestamp)
	if err != nil {
		return err
	}
	if count < 3 {
		return nil
	}

	deviation := math.Abs(currentValue-avg) / avg

	if deviation > WeightWarnDeviation {
		title := "Abnormal weight"
		message := fmt.Sprintf("Weight %.1fkg deviates %.0f%% from the 7-day average (%.1fkg)", currentValue, deviation*100, avg)
		return CreateIfNotExists(cowID, model.MetricTypeWeight, model.AlertSeverityWarning, title, message)
	}

	if deviation < WeightResolveDeviation {
		return ResolveIfExists(cowID, model.MetricTypeWeight)
	}

	return nil
}
