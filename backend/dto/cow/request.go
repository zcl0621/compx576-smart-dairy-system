package cow

import (
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type ListQuery struct {
	common.PageQuery
	Name      string `form:"name"`
	Condition string `form:"condition"`
	Status    string `form:"status"`
	Sort      string `form:"sort,default=updated_at"`
}

type InfoQuery struct {
	ID string `form:"id" binding:"required"`
}

type MetricQuery struct {
	CowID string `form:"cow_id" binding:"required"`
	Range string `form:"range,default=24h"`
}

type CreateRequest struct {
	Name       string             `json:"name" binding:"required"`
	Tag        string             `json:"tag" binding:"required"`
	Age        int                `json:"age"`
	CanMilking bool               `json:"can_milking"`
	Status     model.CowStatus    `json:"status"`
	Condition  model.CowCondition `json:"condition"`
}

type UpdateRequest struct {
	ID         string             `json:"id" binding:"required"`
	Name       string             `json:"name" binding:"required"`
	Tag        string             `json:"tag" binding:"required"`
	Age        int                `json:"age"`
	CanMilking bool               `json:"can_milking"`
	Status     model.CowStatus    `json:"status"`
	Condition  model.CowCondition `json:"condition"`
}
