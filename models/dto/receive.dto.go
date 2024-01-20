package dto

// FilterReceive ...
type FilterReceive struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	Status        int    `json:"status"`
	ReceiveNumber string `json:"receiveNumber"`
	SupplierName  string `json:"supplierName"`
	PurchaseOrderNo string `json:"purchaseOrderNo"`
}

// FilterReceiveDetail ...
type FilterReceiveDetail struct {
	ReceiveNo string `json:"receiveNo"`
	ReceiveID int64  `json:"receiveId"`
}

// ReceiveSaveResult ...
type ReceiveSaveResult struct {
	ErrDesc   string `json:"errDesc"`
	ErrCode   string `json:"errCode"`
	ReceiveNo string `json:"receiveNo"`
	Status    int8   `json:"status"`
	ID        int64  `json:"id"`
}

// ReceiveDetailSaveResult ...
type ReceiveDetailSaveResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
	ID      int64  `json:"id"`
}

// ReceiveDetailResult ...
type ReceiveDetailResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
}
