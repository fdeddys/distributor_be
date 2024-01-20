package dbmodels

import "time"

// StockOpname ...
type StockOpname struct {
	ID              int64     `json:"id" gorm:"column:id"`
	StockOpnameNo   string    `json:"stockOpnameNo" gorm:"column:stock_opname_no"`
	StockOpnameDate time.Time `json:"stockOpnameDate" gorm:"column:stock_opname_date"`
	Note            string    `json:"note" gorm:"column:note"`
	Total           float32   `json:"Total" gorm:"column:total"`
	Status          int8      `json:"status" gorm:"column:status"`
	LastUpdateBy    string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate      time.Time `json:"lastUpdate" gorm:"column:last_update"`
	WarehouseID     int64     `json:"warehouseId" gorm:"column:warehouse_id"`
	Warehouse       Warehouse `json:"warehouse" gorm:"foreignkey:id;association_foreignkey:WarehouseID;"`
}

// TableName ...
func (o *StockOpname) TableName() string {
	return "public.stock_opname"
}
