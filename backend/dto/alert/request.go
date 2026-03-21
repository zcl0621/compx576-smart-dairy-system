package alert

import "github.com/zcl0621/compx576-smart-dairy-system/dto/common"

type ListQuery struct {
	common.PageQuery
	CowID    string `form:"cow_id"`
	Severity string `form:"severity"`
}
