package model

import (
	"time"

	"gorm.io/datatypes"
)

type Report struct {
	BaseModel
	CowID       string         `json:"cow_id" gorm:"not null;index"`
	PeriodStart time.Time      `json:"period_start" gorm:"not null"`
	PeriodEnd   time.Time      `json:"period_end" gorm:"not null"`
	Summary     string         `json:"summary" gorm:"not null"`
	Score       float64        `json:"score" gorm:"not null"`
	Details     datatypes.JSON `json:"details_json" gorm:"column:details_json;type:jsonb;not null"`
}
