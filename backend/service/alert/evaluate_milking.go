package alert

import (
	"fmt"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

func EvaluateMilking(cowID string, currentValue float64, currentTimestamp time.Time) error {
	avg, count, err := sevenDayAverage(cowID, model.MetricTypeMilkAmount, currentTimestamp)
	if err != nil {
		return err
	}
	if count < 3 {
		return nil
	}

	if currentValue < avg*MilkWarnFraction {
		title := "Low milk yield"
		message := fmt.Sprintf("Milk yield %.1fL is more than 30%% below the 7-day average (%.1fL)", currentValue, avg)
		return CreateIfNotExists(cowID, model.MetricTypeMilkAmount, model.AlertSeverityWarning, title, message)
	}

	if currentValue >= avg*MilkResolveFraction {
		return ResolveIfExists(cowID, model.MetricTypeMilkAmount)
	}

	return nil
}
