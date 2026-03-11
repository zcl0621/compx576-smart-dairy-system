package model

type Cow struct {
	BaseModel
	Name       string `json:"name" gorm:"not null"`
	Tag        string `json:"tag" gorm:"not null;uniqueIndex"`
	Age        int    `json:"age" gorm:"default:0"`
	CanMilking bool   `json:"can_milking" gorm:"default:false"`
	Status     string `json:"status" gorm:"default:in_farm"` // in_farm/sold/inactive
}
