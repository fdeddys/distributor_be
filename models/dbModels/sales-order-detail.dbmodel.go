package dbmodels

import "time"

// SalesOrderDetail ...
type SalesOrderDetail struct {
	ID           int64 `json:"id" gorm:"column:id"`
	SalesOrderID int64 `json:"salesOrderId" gorm:"column:sales_order_id"`

	ProductID  int64   `json:"productId" gorm:"column:product_id"`
	Product    Product `json:"product" gorm:"foreignkey:id;association_foreignkey:ProductID;association_autoupdate:false;association_autocreate:false"`
	QtyOrder   int64   `json:"qtyOrder" gorm:"column:qty_order"`
	QtyPicking int64   `json:"qtyPicking" gorm:"column:qty_picking"`
	QtyReceive int64   `json:"qtyReceive" gorm:"column:qty_receive"`

	Price float32 `json:"price" gorm:"column:price"`
	Disc1 float32 `json:"disc1" gorm:"column:disc1"`
	Disc2 float32 `json:"disc2" gorm:"column:disc2"`
	Hpp   float32 `json:"hpp" gorm:"column:hpp"`

	UomID int64  `json:"uomId" gorm:"column:uom"`
	UOM   Lookup `json:"uom" gorm:"foreignkey:id;association_foreignkey:UomID;association_autoupdate:false;association_autocreate:false"`

	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

// TableName ...
func (o *SalesOrderDetail) TableName() string {
	return "public.sales_order_detail"
}
