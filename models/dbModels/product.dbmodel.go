package dbmodels

import (
	"time"
)

//Product model ...
type Product struct {
	ID   int64  `json:"id" gorm:"column:id"`
	Code string `json:"code" gorm:"column:code"`
	Name string `json:"name" gorm:"column:name"`
	PLU  string `json:"plu" gorm:"column:plu"`

	ProductGroupID int64        `json:"productGroupId" gorm:"column:product_group_id"`
	ProductGroup   ProductGroup `gorm:"foreignkey:id; association_foreignkey:ProductGroupID; association_autoupdate:false;association_autocreate:false"`

	BrandID int64 `json:"brandId" gorm:"column:brand_id"`
	Brand   Brand `gorm:"foreignkey:id; association_foreignkey:BrandID; association_autoupdate:false;association_autocreate:false"`

	BigUomID int64  `json:"bigUomId" gorm:"column:big_uom_id"`
	BigUom   Lookup `json:"bigUom"   gorm:"foreignkey:ID; association_foreignkey:BigUomID; association_autoupdate:false;association_autocreate:false"`

	SmallUomID int64  `json:"smallUomId" gorm:"column:small_uom_id"`
	SmallUom   Lookup `json:"smallUom" gorm:"foreignkey:ID; association_foreignkey:SmallUomID;association_autoupdate:false;association_autocreate:false"`

	Status       int       `json:"status" gorm:"column:status"`
	LastUpdateBy string    `json:"lastUpdateBy" gorm:"column:last_update_by"`
	LastUpdate   time.Time `json:"lastUpdate"  gorm:"column:last_update"`
	QtyUom       int16     `json:"qtyUom" gorm:"column:qty_uom"`
	// QtyStock     float32   `json:"qtyStock" gorm:"column:qty_stock"`
	Hpp           float32 `json:"hpp" gorm:"column:hpp"`
	SellPrice     float32 `json:"sellPrice" gorm:"column:sell_price"`
	SellPriceType int     `json:"sellPriceType" gorm:"column:sell_price_type"`
	// SellPricePercent float32 `json:"sellPricePercent" gorm:"column:sell_price_percent"`

	Composition string `json:"composition" gorm:"column:composition"`
}

// TableName ...
func (t *Product) TableName() string {
	return "public.product"
}

// KafkaReq ...
type KafkaReq struct {
	Topic string        `json:"topic"`
	Data  ProductVendor `json:"data"`
}

// ProductVendor ...
type ProductVendor struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Uom  string `json:"unit_of_measurement_code"`
}
