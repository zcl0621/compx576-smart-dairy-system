package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"go.uber.org/zap"
)

// base lat/lng for a farm near Hamilton NZ
const farmLat = -37.7870
const farmLng = 175.2793

type healthAgent struct {
	cowID   string
	baseURL string
	token   string
	status  string
}

func StartHealthAgent(ctx context.Context, cowID, baseURL, token, status string) {
	a := &healthAgent{cowID: cowID, baseURL: baseURL, token: token, status: status}
	a.run(ctx)
}

func (a *healthAgent) run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// send first batch right away
	a.sendAll()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.sendAll()
		}
	}
}

func (a *healthAgent) sendAll() {
	abnormal := shouldBeAbnormal(a.status)

	temp := normalFloat(38.5, 0.5)
	hr := normalFloat(70, 10)
	bo := normalFloat(97.5, 1.0)
	if abnormal {
		temp = normalFloat(40.5, 0.5)
		hr = normalFloat(110, 10)
		bo = normalFloat(88, 3)
	}

	lat := farmLat + (rand.Float64()-0.5)*0.002
	lng := farmLng + (rand.Float64()-0.5)*0.002

	a.send("temperature", temp, "celsius")
	a.send("heart_rate", hr, "bpm")
	a.send("blood_oxygen", bo, "percent")
	a.send("latitude", lat, "degrees")
	a.send("longitude", lng, "degrees")
}

func (a *healthAgent) send(metricType string, value float64, unit string) {
	postMetric(a.baseURL, a.token, a.cowID, "cow_agent", metricType, value, unit)
}

func shouldBeAbnormal(status string) bool {
	switch status {
	case "ill":
		return rand.Float64() < 0.6
	case "unhealth":
		return rand.Float64() < 0.2
	default:
		return rand.Float64() < 0.05
	}
}

func normalFloat(mean, stddev float64) float64 {
	return mean + rand.NormFloat64()*stddev
}

var agentHTTPClient = &http.Client{Timeout: 10 * time.Second}

func postMetric(baseURL, token, cowID, source, metricType string, value float64, unit string) {
	body, _ := json.Marshal(map[string]interface{}{
		"cow_id":       cowID,
		"source":       source,
		"metric_type":  metricType,
		"metric_value": value,
		"unit":         unit,
	})

	req, err := http.NewRequest("POST", baseURL+"/api/metric", bytes.NewReader(body))
	if err != nil {
		projectlog.L().Error("build request failed", zap.Error(err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := agentHTTPClient.Do(req)
	if err != nil {
		projectlog.L().Error("send metric failed", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		projectlog.L().Warn("metric rejected",
			zap.String("cow_id", cowID),
			zap.String("type", metricType),
			zap.Int("status", resp.StatusCode),
		)
	}
}

// FetchToken calls agent server to get a cow JWT
func FetchToken(baseURL, cowID string) (string, error) {
	resp, err := http.Post(fmt.Sprintf("%s/api/token?cow_id=%s", baseURL, url.QueryEscape(cowID)), "", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	token, ok := result["token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("empty token in response")
	}

	return token, nil
}
