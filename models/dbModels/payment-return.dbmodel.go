package dbmodels

import "time"

type PaymentReturn struct {
	ID                 int64            `json:"id" gorm:"column:id"`
	PaymentID          int64            `json:"paymentId" gorm:"column:payment_id"`
	ReturnSalesOrderID int64            `json:"salesOrderReturnId" gorm:"column:return_sales_order_id"`
	ReturnSalesOrder   ReturnSalesOrder `json:"salesOrderReturn" gorm:"foreignkey:id; association_foreignkey:ReturnSalesOrderID; association_autoupdate:false;association_autocreate:false"`
	Total              float32          `json:"total" gorm:"column:total"`
	LastUpdateBy       string           `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate         time.Time        `json:"lastUpdate"  gorm:"column:last_update"`
}

// TableName ...
func (t *PaymentReturn) TableName() string {
	return "public.payment_return"
}
