package agent

import (
	"context"
	"math/rand"
	"time"

	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"go.uber.org/zap"
)

type milkingAgent struct {
	cowID   string
	baseURL string
	token   string
}

func StartMilkingAgent(ctx context.Context, cowID, baseURL, token string) {
	a := &milkingAgent{cowID: cowID, baseURL: baseURL, token: token}
	a.run(ctx)
}

func (a *milkingAgent) run(ctx context.Context) {
	for {
		now := time.Now()
		next := nextMilkingTime(now)
		wait := time.Until(next)

		projectlog.L().Info("milking agent waiting",
			zap.String("cow_id", a.cowID),
			zap.Time("next_milking", next),
		)

		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			a.sendMilking()
		}
	}
}

func (a *milkingAgent) sendMilking() {
	amount := 8.0 + rand.Float64()*7.0      // 8-15 liters
	duration := 300.0 + rand.Float64()*300.0 // 300-600 seconds

	postMetric(a.baseURL, a.token, a.cowID, "milking_machine", "milk_amount", amount, "liters")
	postMetric(a.baseURL, a.token, a.cowID, "milking_machine", "milking_duration", duration, "seconds")
}

// next milking at 6:00 or 18:00
func nextMilkingTime(now time.Time) time.Time {
	today6 := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
	today18 := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())

	if now.Before(today6) {
		return today6
	}
	if now.Before(today18) {
		return today18
	}
	// next day 6am
	return today6.Add(24 * time.Hour)
}
