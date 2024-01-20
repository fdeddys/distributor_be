package dbmodels

import "time"

// ReturnSalesOrder ...
type ReturnSalesOrder struct {
	ID                   int64     `json:"id" gorm:"column:id"`
	ReturnSalesOrderNo   string    `json:"returnNo" gorm:"column:return_no"`
	ReturnSalesOrderDate time.Time `json:"returnDate" gorm:"column:return_date"`

	CustomerID int64    `json:"customerId" gorm:"column:customer_id"`
	Customer   Customer `json:"customer" gorm:"foreignkey:id;association_foreignkey:CustomerID;association_autoupdate:false;association_autocreate:false"`

	WarehouseID int64     `json:"warehouseId" gorm:"column:warehouse_id"`
	Warehouse   Warehouse `json:"warehouse" gorm:"foreignkey:id;association_foreignkey:WarehouseID;association_autoupdate:false;association_autocreate:false"`

	SalesmanID int64    `json:"salesmanId" gorm:"column:sales_id"`
	Salesman   Salesman `json:"salesman" gorm:"foreignkey:id;association_foreignkey:SalesmanID;association_autoupdate:false;association_autocreate:false"`

	ReasonID int64  `json:"reasonId" gorm:"column:reason_id"`
	Reason   Lookup `json:"reason" gorm:"foreignkey:id;association_foreignkey:ReasonID;association_autoupdate:false;association_autocreate:false"`

	IsPaid       bool      `json:"isPaid" gorm:"column:is_paid"`
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
	// 40 = INVOICE
	// 50 = PAID
	Status int8 `json:"status" gorm:"column:status"`
}

// TableName ...
func (o *ReturnSalesOrder) TableName() string {
	return "public.return_sales_order"
}
