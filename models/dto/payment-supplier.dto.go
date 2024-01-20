package dto

type FilterSupplierPayment struct {
	PaymentNo       string `json:"paymentNo"`
	PaymentID       int64  `json:"paymentId"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	PaymentDetailId int64  `json:"paymentDetailId"`
	PaymentStatus   int    `json:"paymentStatus"`
}
