package dbmodels

import (
	"time"
)

//ProductGroup model ...
type ProductGroup struct {
	ID           int64     `json:"id" gorm:"column:id"`
	Name         string    `json:"name" gorm:"column:name"`
	Status       int       `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"last_update_by" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"last_update"`
	Code         string    `json:"code" gorm:"column:code"`
}

// TableName ...
func (t *ProductGroup) TableName() string {
	return "public.product_group"
}
