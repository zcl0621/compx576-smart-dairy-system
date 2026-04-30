package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

// totalMinutesIn7Days converts coverage % to "minutes offline" in the prompt
const totalMinutesIn7Days = 7 * 24 * 60 // 10080

const systemPrompt = `You are a dairy veterinary assistant writing brief 7-day health reports for individual cows. Output must be valid JSON with exactly two fields:

- "summary": one or two sentences on overall condition, max 200 chars
- "note": two or three sentences with concrete observations and any recommended action, max 500 chars

Use plain language and matter-of-fact tone. Never invent metrics not provided. If data is missing, state it directly. No markdown, no bullet points. Do not begin with "this cow" or "the cow".`

// BuildSystemPrompt returns the system prompt. Fixed across calls so DeepSeek caching applies.
func BuildSystemPrompt() string {
	return systemPrompt
}

// BuildUserPrompt builds the per-cow prompt. Uses the no-data variant when the window has no readings.
func BuildUserPrompt(
	cow *model.Cow,
	stats *CowStats,
	alerts []model.Alert,
	score float64,
	periodStart, periodEnd time.Time,
) string {
	if isNoData(stats) {
		return buildNoDataPrompt(cow, score, periodStart, periodEnd)
	}
	return buildFullPrompt(cow, stats, alerts, score, periodStart, periodEnd)
}

func isNoData(stats *CowStats) bool {
	c := stats.Current
	return c.Temperature.Count == 0 &&
		c.HeartRate.Count == 0 &&
		c.BloodOxygen.Count == 0 &&
		c.Weight.Count == 0 &&
		c.MilkAmount.Count == 0
}

func buildNoDataPrompt(cow *model.Cow, score float64, start, end time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Cow: %s\n", cow.Name)
	fmt.Fprintf(&b, "Period: %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
	fmt.Fprintf(&b, "Health score: %.0f/100\n\n", score)
	b.WriteString("No metric data received during this period. Device appears offline for the entire window.\n\n")
	b.WriteString("Reply with JSON only.")
	return b.String()
}

func buildFullPrompt(
	cow *model.Cow,
	stats *CowStats,
	alerts []model.Alert,
	score float64,
	start, end time.Time,
) string {
	cur := stats.Current
	base := stats.Baseline

	var b strings.Builder
	fmt.Fprintf(&b, "Cow: %s\n", cow.Name)
	fmt.Fprintf(&b, "Period: %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
	fmt.Fprintf(&b, "Health score: %.0f/100\n\n", score)

	b.WriteString("7-day metric statistics:\n")
	writeHealthLine(&b, "Temperature", "°C", "38.0-39.0°C", cur.Temperature, ClassifyDeviation(model.MetricTypeTemperature, cur.Temperature.Avg, 0))
	writeHealthLine(&b, "Heart rate", " bpm", "48-84 bpm", cur.HeartRate, ClassifyDeviation(model.MetricTypeHeartRate, cur.HeartRate.Avg, 0))
	writeHealthLine(&b, "Blood oxygen", "%", "90-100%", cur.BloodOxygen, ClassifyDeviation(model.MetricTypeBloodOxygen, cur.BloodOxygen.Avg, 0))
	writeBaselineLine(&b, "Weight", "kg", cur.Weight, base.Weight)
	writeBaselineLine(&b, "Milk yield", "L/day", cur.MilkAmount, base.MilkAmount)

	offlineMinutes := (100 - cur.CoveragePct) / 100 * float64(totalMinutesIn7Days)
	fmt.Fprintf(&b, "\nData coverage: %.0f%% (%.0f minutes offline)\n", cur.CoveragePct, offlineMinutes)

	b.WriteString("\nAlerts in period (chronological, max 20):\n")
	if len(alerts) == 0 {
		b.WriteString("- none\n")
	} else {
		// already capped at 20 by the DB query
		for _, a := range alerts {
			fmt.Fprintf(&b, "- %s [%s] %s\n",
				a.CreatedAt.Format("2006-01-02 15:04"), a.Severity, a.Message)
		}
	}

	b.WriteString("\nReply with JSON only.")
	return b.String()
}

func writeHealthLine(b *strings.Builder, label, unit, normalRange string, s MetricStats, lv DeviationLevel) {
	if s.Count == 0 {
		fmt.Fprintf(b, "- %s: no readings (offline)\n", label)
		return
	}
	fmt.Fprintf(b, "- %s: avg %.1f%s, range %.1f-%.1f%s, status: %s (normal: %s, %d abnormal readings)\n",
		label, s.Avg, unit, s.Min, s.Max, unit, DeviationToStatus(lv), normalRange, s.AbnormalCount)
}

func writeBaselineLine(b *strings.Builder, label, unit string, current, baseline MetricStats) {
	if current.Count == 0 {
		fmt.Fprintf(b, "- %s: no readings (offline)\n", label)
		return
	}
	if baseline.Count < minBaselineCount {
		fmt.Fprintf(b, "- %s: avg %.1f %s (no baseline available, status: unknown)\n",
			label, current.Avg, unit)
		return
	}
	delta := current.Avg - baseline.Avg
	sign := "+"
	if delta < 0 {
		sign = ""
	}
	fmt.Fprintf(b, "- %s: avg %.1f %s, change vs previous 7d: %s%.1f %s\n",
		label, current.Avg, unit, sign, delta, unit)
}
