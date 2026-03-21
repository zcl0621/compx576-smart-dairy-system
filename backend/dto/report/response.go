package report

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type ReportItem struct {
	ID          string              `json:"id" gorm:"column:id"`
	CreatedAt   time.Time           `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time           `json:"updated_at" gorm:"column:updated_at"`
	CowID       string              `json:"cow_id" gorm:"column:cow_id"`
	CowName     string              `json:"cow_name" gorm:"column:cow_name"`
	PeriodStart time.Time           `json:"period_start" gorm:"column:period_start"`
	PeriodEnd   time.Time           `json:"period_end" gorm:"column:period_end"`
	Summary     string              `json:"summary" gorm:"column:summary"`
	Score       float64             `json:"score" gorm:"column:score"`
	Details     model.ReportDetails `json:"details" gorm:"column:details;serializer:json"`
}

type ListResponse struct {
	List []ReportItem `json:"list"`
	common.PageResponse
}
