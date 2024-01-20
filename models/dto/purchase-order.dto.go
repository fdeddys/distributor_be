package dto

type FilterPurchaseOrder struct {
	StartDate           string `json:"startDate"`
	EndDate             string `json:"endDate"`
	Status              int    `json:"status"`
	PurchaseOrderNumber string `json:"purchaseOrderNumber"`
	SupplierId          int64  `json:"supplierId"`
	SupplierName        string `json:"supplierName"`
}

// FilterBatchExpired ...
type FilterBatchExpired struct {
	ExpiredStart string `json:"expiredStart"`
	ExpiredEnd   string `json:"expiredEnd"`
	Batch        string `json:"batch"`
	ProductName string `json:"productName"`
}

// FilterReceiveDetail ...
type FilterPurchaseOrderDetail struct {
	PurchaseOrderNumber string `json:"purchaseOrderNumber"`
	PurchaseOrderID     int64  `json:"purchaseOrderId"`
}

// PurchaseOrderSaveResult ...
type PurchaseOrderSaveResult struct {
	ErrDesc         string `json:"errDesc"`
	ErrCode         string `json:"errCode"`
	PurchaseOrderNo string `json:"purchaseOrderNo"`
	Status          int8   `json:"status"`
	ID              int64  `json:"id"`
}

// PurchaseOrderDetailSaveResult ...
type PurchaseOrderDetailSaveResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
	ID      int64  `json:"id"`
}

type ResultLastPrice struct {
	Price int64   `json:"price"`
	Disc1 int64   `json:"disc1"`
	Disc2 int64   `json:"disc2"`
	Hpp   float32 `json:"hpp"`
}

type ResultLastPrice2 struct {
	Price float32 `json:"price"`
	// Disc1 int64   `json:"disc1"`
	// Disc2 int64   `json:"disc2"`
	Hpp float32 `json:"hpp"`
}
