package reportservice

import (
	"distribution-system-be/database"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/excel"
)

// StockOpnameService ...
type ReportPaymentCashService struct {
}

// Approve ...
func (o ReportPaymentCashService) GenerateReportPaymentCash(filterData dto.FilterReportDate) (filename string, success bool) {

	dateStart := filterData.StartDate + " 00:00:00"
	dateEnd := filterData.EndDate + " 23:59:59"
	datas := generateDataReport(dateStart, dateEnd)
	// filename = ExportToCSV(datas, filterData.StartDate, filterData.EndDate, "report-payment")
	filename, success = excel.ExportToExcelReportPaymentCash(datas, filterData.StartDate, filterData.EndDate, "report-sales")
	return
}

func generateDataReport(dateStart, dateEnd string) []dto.ReportPaymentCash {

	datas := database.ReportPaymentCashByDate(dateStart, dateEnd)
	return datas
}
