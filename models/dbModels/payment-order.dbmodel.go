package dbmodels

import "time"

type PaymentOrder struct {
	ID        int64 `json:"id" gorm:"column:id"`
	PaymentID int64 `json:"paymentId" gorm:"column:payment_id"`

	SalesOrderID int64      `json:"salesOrderId" gorm:"column:sales_order_id"`
	SalesOrder   SalesOrder `json:"salesOrder" gorm:"foreignkey:id; association_foreignkey:SalesOrderID; association_autoupdate:false;association_autocreate:false"`

	Total        float32   `json:"total" gorm:"column:total"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate"  gorm:"column:last_update"`
}

// TableName ...
func (t *PaymentOrder) TableName() string {
	return "public.payment_order"
}
