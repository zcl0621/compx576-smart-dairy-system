package alert

import (
	"context"
	"fmt"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"go.uber.org/zap"
)

// StartOfflineTicker runs EvaluateOffline every minute until ctx is cancelled.
func StartOfflineTicker(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := EvaluateOffline(); err != nil {
				projectlog.L().Error("offline evaluation failed", zap.Error(err))
			}
		}
	}
}

type offlineRow struct {
	ID       string
	LastSeen *time.Time
}

// EvaluateOffline creates or resolves offline alerts for all in-farm cows.
// Only cow_agent metrics count — milking/weight machines run on fixed schedules.
func EvaluateOffline() error {
	var rows []offlineRow
	err := pg.DB.Raw(`
		SELECT c.id, MAX(m.created_at) AS last_seen
		FROM cows c
		LEFT JOIN metrics m
		  ON m.cow_id = c.id
		  AND m.source = ?
		  AND m.deleted_at IS NULL
		WHERE c.status = ?
		  AND c.deleted_at IS NULL
		GROUP BY c.id
	`, model.MetricSourceCowAgent, model.CowStatusInFarm).Scan(&rows).Error
	if err != nil {
		return err
	}

	threshold := time.Now().Add(-time.Duration(OfflineThreshold) * time.Minute)

	for _, row := range rows {
		if row.LastSeen == nil || row.LastSeen.Before(threshold) {
			if err := CreateIfNotExists(row.ID, model.MetricTypeDevice, model.AlertSeverityOffline,
				"Device offline", fmt.Sprintf("No health data received in the last %d minutes", OfflineThreshold)); err != nil {
				projectlog.L().Error("create offline alert failed", zap.String("cow_id", row.ID), zap.Error(err))
			}
		} else {
			if err := ResolveIfExists(row.ID, model.MetricTypeDevice); err != nil {
				projectlog.L().Error("resolve offline alert failed", zap.String("cow_id", row.ID), zap.Error(err))
			}
		}
	}

	return nil
}
