package dto

// FilterOrderDetail ...
type FilterOrderReturnDetail struct {
	OrderReturnNo       string `json:"orderReturnNo"`
	OrderReturnID       int64  `json:"orderReturnId"`
	StartDate           string `json:"startDate"`
	EndDate             string `json:"endDate"`
	OrderReturnDetailId int64  `json:"orderReturnDetailId"`
	QtyReturn           int64  `json:"QtyReturn"`
	SalesOrderReturnID  int64  `json:"salesOrderReturnId"`
	CustomerID          int64  `json:"customerId"`
}

// ReturnOrderSaveResult ...
type ReturnOrderSaveResult struct {
	ErrDesc  string `json:"errDesc"`
	ErrCode  string `json:"errCode"`
	ReturnNo string `json:"returnNo"`
	Status   int8   `json:"status"`
	ID       int64  `json:"id"`
}
