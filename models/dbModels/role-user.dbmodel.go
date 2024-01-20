package dbmodels

import "time"

// Role ...
type RoleUser struct {
	// ID           int64     `json:"id" gorm:"column:id"`
	RoleID       int64     `json:"roleId" gorm:"column:role_id"`
	UserID       int64     `json:"userId" gorm:"column:user_id"`
	Status       int       `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

// TableName ...
func (t *RoleUser) TableName() string {
	return "public.m_role_user"
}
