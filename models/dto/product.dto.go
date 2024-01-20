package dto

import (
	dbmodels "distribution-system-be/models/dbModels"
	"time"
)

//Product model ...
type ProductSearch struct {
	ID   int64  `json:"id" `
	Code string `json:"code" `
	Name string `json:"name" `
	PLU  string `json:"plu"`

	BigUomID int64           `json:"bigUomId" `
	BigUom   dbmodels.Lookup `json:"bigUom"`

	SmallUomID int64           `json:"smallUomId" `
	SmallUom   dbmodels.Lookup `json:"smallUom" `

	Status        int       `json:"status" `
	LastUpdateBy  string    `json:"lastUpdateBy" `
	LastUpdate    time.Time `json:"lastUpdate"`
	QtyUom        int16     `json:"qtyUom"`
	Hpp           float32   `json:"hpp"`
	SellPrice     float32   `json:"sellPrice"`
	SellPriceType int       `json:"sellPriceType"`
	QtyOnHand     int64     `json:"qtyOnHand"`
}
