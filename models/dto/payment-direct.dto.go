package dto

import "time"

type FilterPaymentDirect struct {
	PaymentNo     string `json:"paymentNo"`
	SalesOrderNo  string `json:"salesOrderNo"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	PaymentStatus int64  `json:"paymentStatus"`
}

type PaymentDirectModel struct {
	PaymentStatus int8      `json:"paymentStatus" `
	PaymentNo     string    `json:"paymentNo" `
	SalesOrderNo  string    `json:"salesOrderNo" `
	OrderDate     time.Time `json:"orderDate"`
	SoStatus      int64     `json:"soStatus"`
	GrandTotal    float32   `json:"salesOrderGrandTotal"`
	PaymentID     int64     `json:"paymentId"`
	SalesOrderID  int64     `json:"salesOrderId"`
}
