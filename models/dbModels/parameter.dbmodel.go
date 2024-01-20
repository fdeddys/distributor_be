package dbmodels

// import "time"

type Parameter struct {
	ID           int64     `json:"id" gorm:"column:id"`
	Name         string    `json:"name" gorm:"column:name"`
	Value        string    `json:"value" gorm:"column:value"`
	IsViewable   int8      `json:"IsViewable" gorm:"column:isviewable"`
	// LastUpdateBy string    `json:"last_update_by" gorm:"column:last_update_by"`
	// LastUpdate   time.Time `json:"last_update" gorm:"column:last_update"`
}

// TableName ...
func (t *Parameter) TableName() string {
	return "public.parameter"
}
