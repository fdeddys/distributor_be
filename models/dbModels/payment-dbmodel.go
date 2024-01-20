package dbmodels

import "time"

//Payment model ...
type Payment struct {
	ID          int64  `json:"id" gorm:"column:id"`
	PaymentNo   string `json:"paymentNo" gorm:"column:payment_no"`
	PaymentDate string `json:"paymentDate" gorm:"column:payment_date"`

	CustomerID int64    `json:"customerId" gorm:"column:customer_id"`
	Customer   Customer `gorm:"foreignkey:id; association_foreignkey:CustomerID; association_autoupdate:false;association_autocreate:false"`

	IsCash       bool      `json:"isCash" gorm:"column:is_cash"`
	Note         string    `json:"note" gorm:"column:note"`
	TotalOrder   float32   `json:"totalOrder" gorm:"column:total_order"`
	TotalReturn  float32   `json:"totalReturn" gorm:"column:total_return"`
	TotalPayment float32   `json:"totalPayment" gorm:"column:total_payment"`
	Status       int8      `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate"  gorm:"column:last_update"`
}

// TableName ...
func (t *Payment) TableName() string {
	return "public.payment"
}
