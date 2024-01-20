package dbmodels

// Menu ...
type Menu struct {
	// id bigint NOT NULL DEFAULT nextval('m_menus_id_seq'::regclass),
	// name character varying(30) COLLATE pg_catalog."default" NOT NULL,
	// description character varying(100) COLLATE pg_catalog."default" NOT NULL,
	// last_update_by character varying(100) COLLATE pg_catalog."default",
	// last_update timestamp without time zone,
	// link character varying(200) COLLATE pg_catalog."default",
	// parent_id bigint,
	// icon character varying(50) COLLATE pg_catalog."default",
	// status integer,

	ID          int64  `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Link        string `json:"link" gorm:"column:link"`
	ParentID    int64  `json:"parentId" gorm:"column:parent_id"`
	Icon        string `json:"icon" gorm:"column:icon"`
	Status      int    `json:"status" gorm:"column:status"`
}

// TableName ...
func (t *Menu) TableName() string {
	return "public.m_menus"
}
