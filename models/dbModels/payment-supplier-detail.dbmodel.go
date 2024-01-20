package dbmodels

import "time"

type PaymentSupplierDetail struct {
	ID                int64 `json:"id" gorm:"column:id"`
	PaymentSupplierID int64 `json:"paymentSupplierId" gorm:"column:payment_supplier_id"`

	ReceiveID int64   `json:"receiveId" gorm:"column:receiving_id"`
	Receive   Receive `json:"receive" gorm:"foreignkey:id; association_foreignkey:ReceiveID"`

	Total        float32   `json:"total" gorm:"column:total"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate"  gorm:"column:last_update"`
}

// TableName ...
func (t *PaymentSupplierDetail) TableName() string {
	return "public.payment_supplier_detail"
}
