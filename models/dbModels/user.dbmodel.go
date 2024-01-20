package dbmodels

import (
	"time"
)

// User ...
type User struct {
	ID       int64  `json:"id" gorm:"column:id"`
	UserName string `json:"userName" gorm:"column:user_name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"-" gorm:"column:password"`
	// CurPass  string `json:"curpass" `
	Status int `json:"status" gorm:"column:status"`
	// SupplierCode string    `json:"supplierCode" gorm:"column:suppliercode"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
	FirstName    string    `json:"firstName" gorm:"column:first_name"`
	LastName     string    `json:"lastName" gorm:"column:last_name"`
	IsAdmin      int       `json:"isAdmin" gorm:"column:isadmin"`
	RoleID       int64     `json:"roleId" gorm:"column:role_id"`
	Role         Role      `json:"role" gorm:"foreignkey:RoleID; association_foreignkey:ID"`
}

// TableName ...
func (t *User) TableName() string {
	return "public.m_users"
}
