package dbmodels

// LookupGroup model ...
type LookupGroup struct {
	ID   int64  `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

// TableName ...
func (t *LookupGroup) TableName() string {
	return "public.lookup_group"
}
