package dto

// FilterOrder ...
type FilterOrder struct {
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	InternalStatus string `json:"internalStatus"`
	OrderNumber    string `json:"orderNumber"`
	SalesNo        string `json:"salesNo"`
	MerchantPhone  string `json:"merchantPhone"`
	OrderID        int64  `json:"orderId"`
	CustomerID     int64  `json:"customerId"`
	IsCash         bool   `json:"isCash"`
}

// FilterOrderResult ...
type FilterOrderResult struct {
	ErrDesc      string      `json:"errDesc"`
	ErrCode      string      `json:"errCode"`
	Data         interface{} `json:"data"`
	Page         int         `json:"page"`
	PerPage      int         `json:"per_page"`
	TotalPages   int         `json:"total_pages"`
	TotalRecords int         `json:"total_records"`
}

// OrderSaveResult ...
type OrderSaveResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
	OrderNo string `json:"salesOrderNo"`
	Status  int8   `json:"status"`
	ID      int64  `json:"id"`
}

// FilterOrderDetail ...
type FilterOrderDetail struct {
	OrderNo       string `json:"orderNo"`
	OrderID       int64  `json:"orderId"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	OrderDetailId int64  `json:"orderDetailId"`
	QtyReceive    int64  `json:"qtyReceive"`
	QtyOrder      int64  `json:"qtyOrder"`
}

// SalesOrder ...
type SalesOrder struct {
	AccessToken     string           `json:"access_token"`
	WarehouseCode   string           `json:"warehouse_code"`
	CustomerCode    string           `json:"customer_code"`
	TransactionAt   string           `json:"transaction_at"`
	StateAction     string           `json:"state_action"`
	SalesOrderItems []SalerOrderItem `json:"sales_order_lines"`
	Code            string           `json:"code"`
	SupplierCode    string           `json:"supplier_code"`
}

// SalerOrderItem ...
type SalerOrderItem struct {
	ItemCode    string `json:"item_code"`
	Quantity    string `json:"quantity"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Uom         string `json:"unit_of_measurement_code"`
}

// OrderDetailSaveResult ...
type OrderDetailSaveResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
	ID      int64  `json:"id"`
}

// OrderDetailResult ...
type OrderDetailResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
}

type SaveResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
	OrderNo string `json:"salesOrderNo"`
	Status  int8   `json:"status"`
	ID      int64  `json:"id"`
}
