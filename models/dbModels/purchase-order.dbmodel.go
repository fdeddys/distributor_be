package dbmodels

import "time"

// PurchaseOrder ...
type PurchaseOrder struct {
	ID            int64     `json:"id" gorm:"column:id"`
	PurchaserNo   string    `json:"purchaseOrderNo" gorm:"column:po_no"`
	PurchaserDate time.Time `json:"purchaseOrderDate" gorm:"column:po_date"`

	SupplierID int64    `json:"supplierId" gorm:"column:supplier_id"`
	Supplier   Supplier `json:"supplier" gorm:"foreignkey:id;association_foreignkey:SupplierID;association_autoupdate:false;association_autocreate:false"`

	// SalesmanID int64    `json:"salesmanId" gorm:"column:sales_id"`
	// Salesman   Salesman `json:"salesman" gorm:"foreignkey:id;association_foreignkey:SalesmanID;association_autoupdate:false;association_autocreate:false"`

	Note       string  `json:"note" gorm:"column:note"`
	Tax        float32 `json:"tax" gorm:"column:tax"`
	IsTax      bool    `json:"isTax" gorm:"column:is_tax"`
	Total      float32 `json:"total" gorm:"column:total"`
	GrandTotal float32 `json:"grandTotal" gorm:"column:grand_total"`
	// status
	// 10 = new order
	// 20 = approve
	// 30 = reject
	// 40 = ditarik
	Status       int8      `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

// TableName ...
func (o *PurchaseOrder) TableName() string {
	return "public.po"
}
