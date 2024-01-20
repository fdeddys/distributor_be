package dbmodels

import "time"

type PaymentDetail struct {
	ID        int64 `json:"id" gorm:"column:id"`
	PaymentID int64 `json:"paymentId" gorm:"column:payment_id"`

	PaymentTypeID int64  `json:"paymentTypeId" gorm:"column:payment_type_id"`
	PaymentType   Lookup `json:"paymentType" gorm:"foreignkey:id; association_foreignkey:PaymentTypeID; association_autoupdate:false;association_autocreate:false"`

	PaymentReff  string    `json:"paymentReff" gorm:"column:payment_reff"`
	Total        float32   `json:"total" gorm:"column:total"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate"  gorm:"column:last_update"`
}

// TableName ...
func (t *PaymentDetail) TableName() string {
	return "public.payment_detail"
}
