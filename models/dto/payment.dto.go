package dto

type FilterPayment struct {
	PaymentNo       string `json:"paymentNo"`
	PaymentID       int64  `json:"paymentId"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	PaymentDetailId int64  `json:"paymentDetailId"`
	IsCash          bool   `json:"isCash"`
}

type SavePaymentResult struct {
	ErrDesc   string `json:"errDesc"`
	ErrCode   string `json:"errCode"`
	PaymentNo string `json:"paymentNo"`
	Status    int8   `json:"status"`
	ID        int64  `json:"id"`
}
