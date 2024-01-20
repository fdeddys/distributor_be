package dbmodels

import "time"

// StockMutation ...
type StockMutation struct {
	ID                int64     `json:"id" gorm:"column:id"`
	StockMutationNo   string    `json:"stockMutationNo" gorm:"column:mutation_no"`
	StockMutationDate time.Time `json:"mutationDate" gorm:"column:mutation_date"`
	WarehouseSourceID int64     `json:"warehouseSourceId" gorm:"column:warehouse_source"`
	WarehouseSource   Warehouse `json:"warehouseSource" gorm:"foreignkey:id;association_foreignkey:WarehouseSourceID;association_autoupdate:false;association_autocreate:false"`
	WarehouseDestID   int64     `json:"warehouseDestId" gorm:"column:warehouse_dest"`
	WarehouseDest     Warehouse `json:"warehouseDest" gorm:"foreignkey:id;association_foreignkey:WarehouseDestID;association_autoupdate:false;association_autocreate:false"`
	Note              string    `json:"note" gorm:"column:note"`
	Total             float32   `json:"Total" gorm:"column:total"`
	Requestor         string    `json:"requestor" gorm:"column:requestor"`
	Approver          string    `json:"approver" gorm:"column:approver"`
	// status
	// 10 = new
	// 20 = approve
	// 30 = reject
	Status       int8      `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

// TableName ...
func (o *StockMutation) TableName() string {
	return "public.mutation"
}
