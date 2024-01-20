package dbmodels

import "time"

// ReturnReceive ...
type ReturnReceive struct {
	ID                int64     `json:"id" gorm:"column:id"`
	ReturnReceiveNo   string    `json:"returnNo" gorm:"column:return_receive_no"`
	ReturnReceiveDate time.Time `json:"returnDate" gorm:"column:return_date"`

	SupplierID int64    `json:"supplierId" gorm:"column:supplier_id"`
	Supplier   Supplier `json:"supplier" gorm:"foreignkey:id;association_foreignkey:SupplierID;association_autoupdate:false;association_autocreate:false"`

	WarehouseID int64     `json:"warehouseId" gorm:"column:warehouse_id"`
	Warehouse   Warehouse `json:"warehouse" gorm:"foreignkey:id;association_foreignkey:WarehouseID;association_autoupdate:false;association_autocreate:false"`

	ReasonID int64  `json:"reasonId" gorm:"column:reason_id"`
	Reason   Lookup `json:"reason" gorm:"foreignkey:id;association_foreignkey:ReasonID;association_autoupdate:false;association_autocreate:false"`

	InvoiceNo    string    `json:"invoiceNo" gorm:"column:invoice_no"`
	Note         string    `json:"note" gorm:"column:note"`
	Total        float32   `json:"total" gorm:"column:total"`
	Disc         float32   `json:"disc" gorm:"column:disc"`
	Tax          float32   `json:"tax" gorm:"column:tax"`
	GrandTotal   float32   `json:"grandTotal" gorm:"column:grand_total"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
	// status
	// 10 = new
	// 20 = approve
	// 30 = reject
	Status int8 `json:"status" gorm:"column:status"`
	IsPaid bool `json:"isPaid" gorm:"column:is_paid"`
}

// TableName ...
func (o *ReturnReceive) TableName() string {
	return "public.return_receive"
}
