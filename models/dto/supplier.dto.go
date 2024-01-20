package dto

type FilterSupplierMerchantDto struct {
	MerchantCode          string    `json:"merchant_code" gorm:"column:merchant_code"`
}

type FilterSupplierWarehouseDto struct {
	Code          string    `json:"code" gorm:"column:code"`
}

type FilterSupplierPriceDto struct {
	ProductCode          string    `json:"code" gorm:"column:code"`
}


type FilterSupplierNooChecklistDto struct {
	Name  string `json:"name" gorm:"column:name"`
}

type FilterSupplierNooDocDto struct {
	Name  string `json:"name" gorm:"column:name"`
}