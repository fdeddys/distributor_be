package dbmodels

import "time"

// Role ...
type Role struct {
	// id bigint NOT NULL DEFAULT nextval('m_roles_id_seq'::regclass),
	// name character varying(50) COLLATE pg_catalog."default" NOT NULL,
	// description character varying(255) COLLATE pg_catalog."default" NOT NULL,
	// last_update_by character varying(100) COLLATE pg_catalog."default",
	// last_update timestamp without time zone,

	ID           int64     `json:"id" gorm:"column:id"`
	Name         string    `json:"name" gorm:"column:name"`
	Description  string    `json:"description" gorm:"column:description"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

// TableName ...
func (t *Role) TableName() string {
	return "public.m_roles"
}
