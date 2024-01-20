package dbmodels

import "time"

// ReturnSalesOrderDetail ...
type ReturnSalesOrderDetail struct {
	ID                 int64 `json:"id" gorm:"column:id"`
	ReturnSalesOrderID int64 `json:"returnSalesOrderId" gorm:"column:return_sales_order_id"`

	ProductID int64   `json:"productId" gorm:"column:product_id"`
	Product   Product `json:"product" gorm:"foreignkey:id;association_foreignkey:ProductID;association_autoupdate:false;association_autocreate:false"`

	Qty   int64   `json:"qty" gorm:"column:qty"`
	Price float32 `json:"price" gorm:"column:price"`
	Disc1 float32 `json:"disc1" gorm:"column:disc1"`
	Disc2 float32 `json:"disc2" gorm:"column:disc2"`

	UomID int64  `json:"uomId" gorm:"column:uom"`
	UOM   Lookup `json:"uom" gorm:"foreignkey:id;association_foreignkey:UomID;association_autoupdate:false;association_autocreate:false"`

	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

// TableName ...
func (o *ReturnSalesOrderDetail) TableName() string {
	return "public.return_sales_order_detail"
}
