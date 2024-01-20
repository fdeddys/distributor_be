package dto

type FilterReturnReceive struct {
	ReturnNo       string `json:"returnNo"`
	ReturnID       int64  `json:"returnId"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	ReturnDetailId int64  `json:"returnDetailId"`
	Qty            int64  `json:"qty"`
}
