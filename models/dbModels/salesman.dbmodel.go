package dbmodels

import "time"

//Salesman model ...
type Salesman struct {
	ID           int64     `json:"id" gorm:"column:id"`
	Code         string    `json:"code" gorm:"column:code"`
	Name         string    `json:"name" gorm:"column:name"`
	Description  string    `json:"description" gorm:"column:description"`
	Status       int       `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"last_update_by" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"last_update" gorm:"column:last_update"`
}

// TableName ...
func (t *Salesman) TableName() string {
	return "public.sales"
}
