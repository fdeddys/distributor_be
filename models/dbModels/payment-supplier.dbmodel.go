package dbmodels

import "time"

//Payment model ...
type PaymentSupplier struct {
	ID          int64  `json:"id" gorm:"column:id"`
	PaymentNo   string `json:"paymentNo" gorm:"column:payment_no"`
	PaymentDate string `json:"paymentDate" gorm:"column:payment_date"`

	SupplierID int64    `json:"supplierId" gorm:"column:supplier_id"`
	Supplier   Supplier `json:"supplier" gorm:"foreignkey:id; association_foreignkey:SupplierID"`

	PaymentTypeID int64  `json:"paymentTypeId" gorm:"column:payment_type_id"`
	PaymentType   Lookup `json:"paymentType" gorm:"foreignkey:id; association_foreignkey:PaymentTypeID"`

	PaymentReff  string    `json:"paymentReff" gorm:"column:payment_reff"`
	Note         string    `json:"note" gorm:"column:note"`
	Total        float32   `json:"total" gorm:"column:total"`
	Status       int8      `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate"  gorm:"column:last_update"`
}

// TableName ...
func (t *PaymentSupplier) TableName() string {
	return "public.payment_supplier"
}
