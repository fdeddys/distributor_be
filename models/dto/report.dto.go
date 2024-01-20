package dto

type ReportPaymentCash struct {
	PaymentTypeName string
	PaymentNo       string
	PaymentDate1    string
	SalesOrderNo    string
	OrderDate       string
	TotalOrder      int64
	TotalPpn        int64
	TotalPayment    int64
	LastUpdate      string
	LastUpdateBy    string
}

type FilterReportDate struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type ReportSales struct {
	OrderDate    string
	SalesOrderNo string
	Status       string
	Plu          string
	ProductName  string
	QtyOrder     int64
	Uom          string
	Price        int64
	Disc1        int64
}

type ReportPaymentSupplier struct {
	PaymentNo   string
	PaymentDate string
	Supplier    string
	PaymentType string
	PaymentReff string
	ReceiveNo   string
	ReceiveTgl  string
	Status      string
	GrandTotal  float32
}
