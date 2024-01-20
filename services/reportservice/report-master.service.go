package reportservice

import (
	"distribution-system-be/database"
	"distribution-system-be/utils/excel"
)

// StockOpnameService ...
type ReportMasterService struct {
}

// Approve ...
func (o ReportMasterService) GenerateReportMasterProduct() (filename string, success bool) {

	datas := database.ProductList()
	filename, success = excel.ExportToExcelReportMasterProduct(datas, "report-product")
	return
}
