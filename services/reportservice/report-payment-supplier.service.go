package reportservice

import (
	"distribution-system-be/database"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/excel"
)

// StockOpnameService ...
type ReportPaymentSupplierService struct {
}

// Approve ...
func (o ReportPaymentSupplierService) GenerateReportPaymentSupplier(filterData dto.FilterReportDate) (filename string, success bool) {

	dateStart := filterData.StartDate + " 00:00:00"
	dateEnd := filterData.EndDate + " 23:59:59"
	datas := generateDataReportPaymentSupplier(dateStart, dateEnd)
	// filename = ExportToCSV(datas, filterData.StartDate, filterData.EndDate, "report-payment")

	filename, success = excel.ExportToExcelReportPaymentSupplier(datas, filterData.StartDate, filterData.EndDate, "report-sales")
	return
}

func generateDataReportPaymentSupplier(dateStart, dateEnd string) []dto.ReportPaymentSupplier {

	datas := database.ReportPaymentSupplierByDate(dateStart, dateEnd)

	return datas
}
