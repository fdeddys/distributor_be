package dbmodels

// Lookup model ...
type Lookup struct {
	ID            int64       `json:"id" gorm:"column:id"`
	Status        int64       `json:"status" gorm:"column:status"`
	Code          string      `json:"code" gorm:"column:code"`
	LookupGroupID int64       `json:"lookupGroupId" gorm:"column:lookup_group_id"`
	LookupGroup   LookupGroup `json:"lookupGroup" gorm:"foreignkey:id; association_foreignkey:LookupGroupID"`
	Name          string      `json:"name" gorm:"column:name"`
	IsViewable    int8        `json:"isViewable" gorm:"column:is_viewable"`
}

// TableName ...
func (t *Lookup) TableName() string {
	return "public.lookup"
}
