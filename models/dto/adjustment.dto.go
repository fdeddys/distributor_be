package dto

// FilterAdjustment ...
type FilterAdjustment struct {
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	Status           string `json:"status"`
	AdjustmentNumber string `json:"adjustmentNumber"`
}

// FilterAdjustmentDetail ...
type FilterAdjustmentDetail struct {
	AdjustmentNo string `json:"adjustmentNo"`
	AdjustmentID int64  `json:"adjustmentId"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
}

// AdjustmentSaveResult ...
type AdjustmentSaveResult struct {
	ErrDesc      string `json:"errDesc"`
	ErrCode      string `json:"errCode"`
	AdjustmentNo string `json:"adjustmentNo"`
	Status       int8   `json:"status"`
	ID           int64  `json:"id"`
}

// AdjustmentDetailSaveResult ...
type AdjustmentDetailSaveResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
	ID      int64  `json:"id"`
}

// AdjustmentDetailResult ...
type AdjustmentDetailResult struct {
	ErrDesc string `json:"errDesc"`
	ErrCode string `json:"errCode"`
}
