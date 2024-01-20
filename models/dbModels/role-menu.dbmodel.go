package dbmodels

import "time"

// RoleMenu ...
type RoleMenu struct {
	RoleID       int       `json:"roleId" gorm:"column:role_id"`
	MenuID       int       `json:"menuId" gorm:"column:menu_id"`
	Status       int       `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column(last_update);type(timestamp without time zone);null"`
}

// TableName ...
func (t *RoleMenu) TableName() string {
	return "public.m_role_menu"
}
