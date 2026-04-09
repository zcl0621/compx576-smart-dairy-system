package cow

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type ListItem struct {
	ID         string             `json:"id" gorm:"column:id"`
	CreatedAt  time.Time          `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time          `json:"updated_at" gorm:"column:updated_at"`
	Name       string             `json:"name" gorm:"column:name"`
	Tag        string             `json:"tag" gorm:"column:tag"`
	Age        int                `json:"age" gorm:"column:age"`
	CanMilking bool               `json:"can_milking" gorm:"column:can_milking"`
	Status     model.CowStatus    `json:"status" gorm:"column:status"`
	Condition  model.CowCondition `json:"condition" gorm:"column:condition"`
}

type ListResponse struct {
	List []ListItem `json:"list"`
	common.PageResponse
}

type InfoResponse struct {
	ID          string             `json:"id" gorm:"column:id"`
	CreatedAt   time.Time          `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time          `json:"updated_at" gorm:"column:updated_at"`
	Name        string             `json:"name" gorm:"column:name"`
	Tag         string             `json:"tag" gorm:"column:tag"`
	Age         int                `json:"age" gorm:"column:age"`
	CanMilking  bool               `json:"can_milking" gorm:"column:can_milking"`
	Status      model.CowStatus    `json:"status" gorm:"column:status"`
	Condition   model.CowCondition `json:"condition" gorm:"column:condition"`
	Weight      *float64           `json:"weight" gorm:"-"`
	Temperature *float64           `json:"temperature" gorm:"-"`
	HeartRate   *float64           `json:"heart_rate" gorm:"-"`
	BloodOxygen *float64           `json:"blood_oxygen" gorm:"-"`
	MilkAmount  *float64           `json:"milk_amount" gorm:"-"`
}

type MetricPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

type MovementPoint struct {
	Time      time.Time `json:"time"`
	DistanceM float64   `json:"distance_m"`
}

type TemperatureMetricSummary struct {
	Current *float64                 `json:"current"`
	Avg     *float64                 `json:"avg"`
	Min     *float64                 `json:"min"`
	Max     *float64                 `json:"max"`
	Status  model.ReportMetricStatus `json:"status"`
}

type TemperatureMetricResponse struct {
	CowID     string                   `json:"cow_id"`
	Range     model.MetricRange        `json:"range"`
	UpdatedAt *time.Time               `json:"updated_at"`
	Summary   TemperatureMetricSummary `json:"summary"`
	Series    []MetricPoint            `json:"series"`
}

type HeartRateMetricSummary struct {
	Current *float64                 `json:"current"`
	Avg     *float64                 `json:"avg"`
	Min     *float64                 `json:"min"`
	Max     *float64                 `json:"max"`
	Status  model.ReportMetricStatus `json:"status"`
}

type HeartRateMetricResponse struct {
	CowID     string                 `json:"cow_id"`
	Range     model.MetricRange      `json:"range"`
	UpdatedAt *time.Time             `json:"updated_at"`
	Summary   HeartRateMetricSummary `json:"summary"`
	Series    []MetricPoint          `json:"series"`
}

type BloodOxygenMetricSummary struct {
	Current *float64                 `json:"current"`
	Avg     *float64                 `json:"avg"`
	Min     *float64                 `json:"min"`
	Max     *float64                 `json:"max"`
	Status  model.ReportMetricStatus `json:"status"`
}

type BloodOxygenMetricResponse struct {
	CowID     string                   `json:"cow_id"`
	Range     model.MetricRange        `json:"range"`
	UpdatedAt *time.Time               `json:"updated_at"`
	Summary   BloodOxygenMetricSummary `json:"summary"`
	Series    []MetricPoint            `json:"series"`
}

type MilkAmountMetricSummary struct {
	Total         float64 `json:"total"`
	AvgPerSession float64 `json:"avg_per_session"`
	SessionCount  int64   `json:"session_count"`
}

type MilkAmountMetricResponse struct {
	CowID     string                  `json:"cow_id"`
	Range     model.MetricRange       `json:"range"`
	UpdatedAt *time.Time              `json:"updated_at"`
	Summary   MilkAmountMetricSummary `json:"summary"`
	Series    []MetricPoint           `json:"series"`
}

type MovementMetricSummary struct {
	DistanceM  float64                  `json:"distance_m"`
	PointCount int64                    `json:"point_count"`
	Status     model.ReportMetricStatus `json:"status"`
}

type MovementMetricResponse struct {
	CowID     string                `json:"cow_id"`
	Range     model.MetricRange     `json:"range"`
	UpdatedAt *time.Time            `json:"updated_at"`
	Summary   MovementMetricSummary `json:"summary"`
	Series    []MovementPoint       `json:"series"`
}

type WeightMetricSummary struct {
	Current *float64 `json:"current"`
	Avg     *float64 `json:"avg"`
	Min     *float64 `json:"min"`
	Max     *float64 `json:"max"`
}

type WeightMetricResponse struct {
	CowID     string              `json:"cow_id"`
	Range     model.MetricRange   `json:"range"`
	UpdatedAt *time.Time          `json:"updated_at"`
	Summary   WeightMetricSummary `json:"summary"`
	Series    []MetricPoint       `json:"series"`
}

type MovementPathPoint struct {
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Time        int64   `json:"time"`
	StaySeconds int64   `json:"stay_seconds"`
}

type MovementPathResponse struct {
	CowID  string              `json:"cow_id"`
	Range  model.MetricRange   `json:"range"`
	Points []MovementPathPoint `json:"points"`
}
