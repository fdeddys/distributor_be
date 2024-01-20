package dto

type FilterStockOpname struct {
	StartDate           string `json:"startDate"`
	EndDate             string `json:"endDate"`
	StockOpnameNumber   string `json:"stockOpnameNumber"`
	StockOpnameID       int64  `json:"stockOpnameId"`
	InternalStatus      string `json:"internalStatus"`
	StockOpnameDetailId int64  `json:"stockOpnameDetailId"`
	Qty                 int64  `json:"qty"`
}

// StockOpnameSaveResult ...
type StockOpnameSaveResult struct {
	ErrDesc       string `json:"errDesc"`
	ErrCode       string `json:"errCode"`
	StockOpnameNo string `json:"stockOpnameNo"`
	Status        int8   `json:"status"`
	ID            int64  `json:"id"`
}

type TemplateReportStockOpname struct {
	ProductID   int64
	ProductName string
	Qty         int64
	UomName     string
	UomID       int64
}
