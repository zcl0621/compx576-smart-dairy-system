package report

import "github.com/zcl0621/compx576-smart-dairy-system/dto/common"

type ListQuery struct {
	common.PageQuery
}

type LatestQuery struct {
	CowID string `form:"cow_id" binding:"required"`
}
