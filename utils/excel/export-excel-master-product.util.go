package excel

import (
	dbmodels "distribution-system-be/models/dbModels"
	"fmt"
	"time"

	excelize "github.com/xuri/excelize/v2"
)

func ExportToExcelReportMasterProduct(datas []dbmodels.Product, namaFile string) (filename string, success bool) {

	t1, _ := time.Now().ISOWeek()
	//time.Parse("2006-01-02", time.Now().String())

	filename = fmt.Sprintf("%v_%v.xls", namaFile, t1)

	sheet1Name := "Sheet1"
	xls := excelize.NewFile()
	xls.NewSheet(sheet1Name)

	no := 1
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "REPORT PRODUCT")

	no = no + 1
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "No")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), "Name")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), "Code")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), "Big UOM")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), "Small UOM")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), "Qty UOM")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), "Status")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), "Sell Price")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), "PLU")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), "Composition")
	urut := 0

	for _, rs := range datas {
		no++
		urut++
		xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), urut)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), rs.Name)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), rs.Code)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), rs.BigUom.Name)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), rs.SmallUom.Name)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), rs.QtyUom)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), rs.Status)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), rs.SellPrice)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), rs.PLU)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), rs.Composition)
	}

	if err := xls.SaveAs(filename); err != nil {
		success = false
		return
	}
	success = true
	return
}
