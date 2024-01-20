package dbmodels

import "time"

// ReceiveDetail ...
type ReceiveDetail struct {
	ID        int64 `json:"id" gorm:"column:id"`
	ReceiveID int64 `json:"receiveId" gorm:"column:receive_id"`

	ProductID int64   `json:"productId" gorm:"column:product_id"`
	Product   Product `json:"product" gorm:"foreignkey:id;association_foreignkey:ProductID;association_autoupdate:false;association_autocreate:false"`
	Qty       int64   `json:"qty" gorm:"column:qty"`
	Price     float32 `json:"price" gorm:"column:price"`
	Disc1     float32 `json:"disc1" gorm:"column:disc1"`
	Disc2     float32 `json:"disc2" gorm:"column:disc2"`
	Hpp       float32 `json:"hpp" gorm:"column:hpp"`

	UomID int64  `json:"uomId" gorm:"column:uom"`
	UOM   Lookup `json:"uom" gorm:"foreignkey:id;association_foreignkey:UomID;association_autoupdate:false;association_autocreate:false"`

	BatchNo string `json:"batchNo" gorm:"column:batch_no"`
	Ed      string `json:"ed" gorm:"column:ed"`

	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate" gorm:"column:last_update"`

	QtyWh int64 `json:"qtyWh" gorm:"-" sql:"-"`
}

// TableName ...
func (o *ReceiveDetail) TableName() string {
	return "public.receive_detail"
}
