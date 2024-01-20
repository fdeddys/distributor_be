package dto

type FilterStockMutation struct {
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	StockMutationNumber   string `json:"stockMutationNumber"`
	StockMutationID       int64  `json:"stockMutationId"`
	InternalStatus        string `json:"internalStatus"`
	StockMutationDetailId int64  `json:"stockMutationDetailId"`
	Qty                   int64  `json:"qty"`
}

// StockMutationSaveResult ...
type StockMutationSaveResult struct {
	ErrDesc         string `json:"errDesc"`
	ErrCode         string `json:"errCode"`
	StockMutationNo string `json:"stockMutationNo"`
	Status          int8   `json:"status"`
	ID              int64  `json:"id"`
}
