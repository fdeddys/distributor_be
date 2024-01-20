package dbmodels

import "time"

//Stock model ...
type Stock struct {
	ID          int64     `json:"id" gorm:"column:id"`
	ProductID   int64     `json:"productId" gorm:"column:product_id"`
	WarehouseID int64     `json:"warehouseId" gorm:"column:warehouse_id"`
	Warehouse   Warehouse `json:"warehouse" gorm:"foreignkey:id;association_foreignkey:WarehouseID"`
	Qty         int64     `json:"qty" gorm:"column:qty"`
	// Hpp          float32   `json:"hpp" gorm:"column:hpp"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

// TableName ...
func (t *Stock) TableName() string {
	return "public.stock"
}
