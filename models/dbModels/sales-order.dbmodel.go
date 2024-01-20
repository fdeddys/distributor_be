package dbmodels

import "time"

// SalesOrder ...
type SalesOrder struct {
	ID           int64     `json:"id" gorm:"column:id"`
	SalesOrderNo string    `json:"salesOrderNo" gorm:"column:sales_order_no"`
	OrderDate    time.Time `json:"orderDate" gorm:"column:order_date"`
	CustomerID   int64     `json:"customerId" gorm:"column:customer_id"`
	Customer     Customer  `json:"customer" gorm:"foreignkey:id;association_foreignkey:CustomerID;association_autoupdate:false;association_autocreate:false"`
	Note         string    `json:"note" gorm:"column:note"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`
	Tax          float32   `json:"tax" gorm:"column:tax"`
	Total        float32   `json:"total" gorm:"column:total"`
	GrandTotal   float32   `json:"grandTotal" gorm:"column:grand_total"`

	// SalesmanID int64 `json:"salesmanId" gorm:"column:salesman_id"`
	// Salesman   User  `json:"salesman" gorm:"foreignkey:id;association_foreignkey:SalesmanID;association_autoupdate:false;association_autocreate:false"`

	// status
	// 10 = new order
	// 20 = approve
	// 30 = reject
	// 40 = INVOICE
	// 50 = PAID
	Status int8 `json:"status" gorm:"column:status"`
	Top    int8 `json:"top" gorm:"column:top"`
	IsCash bool `json:"isCash" gorm:"column:is_cash"`

	WarehouseID int64     `json:"warehouseId" gorm:"column:warehouse_id"`
	Warehouse   Warehouse `json:"warehouse" gorm:"foreignkey:id;association_foreignkey:WarehouseID;association_autoupdate:false;association_autocreate:false"`
	SalesmanID  int64     `json:"salesmanId" gorm:"column:sales_id"`
	Salesman    Salesman  `json:"salesman" gorm:"foreignkey:id;association_foreignkey:SalesmanID;association_autoupdate:false;association_autocreate:false"`

	PickingNo   string    `json:"pickingNo" gorm:"column:picking_no"`
	PickingDate time.Time `json:"pickingDate" gorm:"column:picking_date"`
	PickingUser string    `json:"pickingUser" gorm:"column:picking_user"`
	IsPaid      bool      `json:"isPaid" gorm:"column:is_paid"`

	DeliveryNo     string    `json:"deliveryNo" gorm:"column:delivery_no"`
	DeliveryDate   time.Time `json:"deliveryDate" gorm:"column:delivery_date"`
	DeliveryDriver string    `json:"deliveryDriver" gorm:"column:delivery_driver"`
	InvoiceNo      string    `json:"invoiceNo" gorm:"column:invoice_no"`
}

// TableName ...
func (o *SalesOrder) TableName() string {
	return "public.sales_order"
}
