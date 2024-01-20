package reportservice

import (
	"distribution-system-be/database"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/excel"
)

// StockOpnameService ...
type ReportSalesOrderService struct {
}

// Approve ...
func (o ReportSalesOrderService) GenerateReport(filterData dto.FilterReportDate) (filename string, success bool) {

	dateStart := filterData.StartDate + " 00:00:00"
	dateEnd := filterData.EndDate + " 23:59:59"
	datas := generateDataReportSales(dateStart, dateEnd)
	// filename = ExportToCSV(datas, filterData.StartDate, filterData.EndDate, "report-sales-order")
	filename, success = excel.ExportToExcelReportSales(datas, filterData.StartDate, filterData.EndDate, "report-sales")
	return
}

func generateDataReportSales(dateStart, dateEnd string) []dto.ReportSales {

	datas := database.ReportSalesByDate(dateStart, dateEnd)
	return datas
}
