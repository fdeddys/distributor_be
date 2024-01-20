package dbmodels

import (
	"time"
)


// Customer ...
type Customer struct {
	ID            int64     `json:"id" gorm:"column:id"`
	Code          string    `json:"code" gorm:"column:code"`
	Name          string    `json:"name" gorm:"column:name"`
	Top           int8      `json:"top" gorm:"column:top"`
	Status        int       `json:"status" gorm:"column:status"`
	LastUpdateBy  string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate    time.Time `json:"lastUpdate" gorm:"column:last_update"`
	Address1      string    `json:"address1" gorm:"column:address1"`
	Address2      string    `json:"address2" gorm:"column:address2"`
	Address3      string    `json:"address3" gorm:"column:address3"`
	Address4      string    `json:"address4" gorm:"column:address4"`
	ContactPerson string    `json:"contactPerson" gorm:"column:contact_person"`
	PhoneNumber   string    `json:"phoneNumber" gorm:"column:phone_number"`
}

// TableName ...
func (m *Customer) TableName() string {
	return "public.customer"
}
