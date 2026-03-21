package alert

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type ListItem struct {
	ID         string              `json:"id" gorm:"column:id"`
	CreatedAt  time.Time           `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time           `json:"updated_at" gorm:"column:updated_at"`
	CowID      string              `json:"cow_id" gorm:"column:cow_id"`
	CowName    string              `json:"cow_name" gorm:"column:cow_name"`
	ReportID   *string             `json:"report_id" gorm:"column:report_id"`
	MetricKey  model.MetricType    `json:"metric_key" gorm:"column:metric_key"`
	Title      string              `json:"title" gorm:"column:title"`
	Message    string              `json:"message" gorm:"column:message"`
	Severity   model.AlertSeverity `json:"severity" gorm:"column:severity"`
	Status     model.AlertStatus   `json:"status" gorm:"column:status"`
	ResolvedAt *time.Time          `json:"resolved_at" gorm:"column:resolved_at"`
}

type ListResponse struct {
	List []ListItem `json:"list"`
	common.PageResponse
}

type SummaryResponse struct {
	Active   int64 `json:"active"`
	Warning  int64 `json:"warning"`
	Critical int64 `json:"critical"`
	Offline  int64 `json:"offline"`
}
