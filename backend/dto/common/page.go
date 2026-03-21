package common

const DefaultPageLimit = 20
const MinPageLimit = 1
const MaxPageLimit = 50

type PageQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type PageResponse struct {
	Page       int   `json:"page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (p *PageQuery) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *PageQuery) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PageQuery) GetLimit() int {
	if p.PageSize < MinPageLimit {
		return DefaultPageLimit
	}
	if p.PageSize > MaxPageLimit {
		return MaxPageLimit
	}
	return p.PageSize
}

func (p *PageQuery) GetTotalPages(total int64) int {
	if total <= 0 {
		return 0
	}

	limit := int64(p.GetLimit())
	return int((total + limit - 1) / limit)
}
