package agent

import (
	"context"
	"math/rand"
	"time"

	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"go.uber.org/zap"
)

type weightAgent struct {
	cowID      string
	baseURL    string
	token      string
	lastWeight float64
}

func StartWeightAgent(ctx context.Context, cowID, baseURL, token string) {
	a := &weightAgent{
		cowID:      cowID,
		baseURL:    baseURL,
		token:      token,
		lastWeight: 400.0 + rand.Float64()*300.0, // initial 400-700 kg
	}
	a.run(ctx)
}

func (a *weightAgent) run(ctx context.Context) {
	for {
		now := time.Now()
		next := nextWeighTime(now)
		wait := time.Until(next)

		projectlog.L().Info("weight agent waiting",
			zap.String("cow_id", a.cowID),
			zap.Time("next_weigh", next),
		)

		t := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			a.sendWeight()
		}
	}
}

func (a *weightAgent) sendWeight() {
	// +/- 2 kg daily fluctuation
	a.lastWeight += (rand.Float64() - 0.5) * 4.0
	if a.lastWeight < 350 {
		a.lastWeight = 350
	}
	if a.lastWeight > 750 {
		a.lastWeight = 750
	}

	postMetric(a.baseURL, a.token, a.cowID, "weight_machine", "weight", a.lastWeight, "kg")
}

// next weigh at 12:00 noon
func nextWeighTime(now time.Time) time.Time {
	today12 := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	if now.Before(today12) {
		return today12
	}
	return today12.Add(24 * time.Hour)
}
