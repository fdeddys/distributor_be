package excel

import (
	"fmt"
	"time"

	"distribution-system-be/database"
	dbModels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/appstatus"

	excelize "github.com/xuri/excelize/v2"
)

func ExportToExcelPo(pos []dbModels.PurchaseOrder, namaFile string) bool {

	sheet1Name := "Sheet1"
	xls := excelize.NewFile()
	xls.NewSheet(sheet1Name)

	xls.SetCellValue(sheet1Name, "A1", "PoNo")
	xls.SetCellValue(sheet1Name, "B1", "Date")
	xls.SetCellValue(sheet1Name, "C1", "Supplier")
	xls.SetCellValue(sheet1Name, "D1", "Status")
	xls.SetCellValue(sheet1Name, "E1", "Product")
	xls.SetCellValue(sheet1Name, "F1", "Price")
	xls.SetCellValue(sheet1Name, "G1", "qty")
	xls.SetCellValue(sheet1Name, "H1", "unit")

	no := 1
	for _, po := range pos {

		podetails := database.GetAllDataDetailPurchaseOrder(po.ID)
		for _, poDetail := range podetails {
			no++
			xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), po.PurchaserNo)
			xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), po.PurchaserDate.Format("2006-01-02"))
			xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), po.Supplier.Name)
			xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), appstatus.GetStatus(po.Status))

			xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), poDetail.Product.Name)
			xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), poDetail.PoPrice)
			xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), poDetail.PoQty)
			xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), poDetail.PoUOM.Name)
		}
	}

	if err := xls.SaveAs(namaFile); err != nil {
		return false
	}
	return true
}

func ExportToExcelReceive(receives []dbModels.Receive, namaFile string) bool {

	sheet1Name := "Sheet1"
	xls := excelize.NewFile()
	xls.NewSheet(sheet1Name)

	xls.SetCellValue(sheet1Name, "A1", "No")
	xls.SetCellValue(sheet1Name, "B1", "ReceiveNo")
	xls.SetCellValue(sheet1Name, "C1", "Date")
	xls.SetCellValue(sheet1Name, "D1", "Supplier")
	xls.SetCellValue(sheet1Name, "E1", "Status")
	xls.SetCellValue(sheet1Name, "F1", "Total")
	// xls.SetCellValue(sheet1Name, "E1", "Product")
	// xls.SetCellValue(sheet1Name, "F1", "Price")
	// xls.SetCellValue(sheet1Name, "G1", "qty")
	// xls.SetCellValue(sheet1Name, "H1", "unit")

	no := 1
	for _, receive := range receives {
		no++
		xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), no-1)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), receive.ReceiveNo)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), receive.ReceiveDate.Format("2006-01-02"))
		xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), receive.Supplier.Name)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), appstatus.GetStatus(receive.Status))
		xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), receive.GrandTotal)
		// podetails := database.GetAllDataDetailPurchaseOrder(po.ID)
		// for _, poDetail := range podetails {

		// 	xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), poDetail.Product.Name)
		// 	xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), poDetail.PoPrice)
		// 	xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), poDetail.PoQty)
		// 	xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), poDetail.PoUOM.Name)
		// }
	}

	if err := xls.SaveAs(namaFile); err != nil {
		return false
	}
	return true
}

func ExportToExcelReportSales(reportSales []dto.ReportSales, dateStart, dateEnd, namaFile string) (filename string, success bool) {

	filename = fmt.Sprintf("%v_%v_%v.csv", namaFile, dateStart, dateEnd)

	sheet1Name := "Sheet1"
	xls := excelize.NewFile()
	xls.NewSheet(sheet1Name)

	t1, _ := time.Parse("2006-01-02", dateStart)
	t2, _ := time.Parse("2006-01-02", dateEnd)
	fmt.Println("Waktu nya adalah", t1.Format("02-January-2006"))
	no := 1
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "REPORT SALES")

	no++
	// xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), dateStart)
	// xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), dateEnd)
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), fmt.Sprintf("%v - %v ", t1.Format("02-January-2006"), t2.Format("02-January-2006")))

	no = no + 2
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "No")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), "OrderDate")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), "SalesOrderNo")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), "Status")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), "PLU")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), "ProductName")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), "QtyOrder")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), "Uom")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), "Price")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), "Disc1")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("K%d", no), "Final Price")
	urut := 0
	for _, rs := range reportSales {
		no++
		urut++
		xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), urut)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), rs.OrderDate)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), rs.SalesOrderNo)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), rs.Status)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), rs.Plu)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), rs.ProductName)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), rs.QtyOrder)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), rs.Uom)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), rs.Price)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), rs.Disc1)
		total := rs.QtyOrder * rs.Price
		disc := total * rs.Disc1 / 100
		finalPrice := total - disc
		xls.SetCellValue(sheet1Name, fmt.Sprintf("K%d", no), finalPrice)

	}

	if err := xls.SaveAs(filename); err != nil {
		success = false
		return
	}
	success = true
	return
}

func ExportToExcelReportPaymentCash(reportPayements []dto.ReportPaymentCash, dateStart, dateEnd, namaFile string) (filename string, success bool) {

	filename = fmt.Sprintf("%v_%v_%v.csv", namaFile, dateStart, dateEnd)

	sheet1Name := "Sheet1"
	xls := excelize.NewFile()
	xls.NewSheet(sheet1Name)

	t1, _ := time.Parse("2006-01-02", dateStart)
	t2, _ := time.Parse("2006-01-02", dateEnd)
	fmt.Println("Waktu nya adalah", t1.Format("02-January-2006"))
	no := 1
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "REPORT PAYMENT")

	no++
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), fmt.Sprintf("%v - %v ", t1.Format("02-January-2006"), t2.Format("02-January-2006")))

	no = no + 2
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "No")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), "PaymentTypeName")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), "PaymentNo")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), "PaymentDate")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), "SalesOrderNo")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), "OrderDate")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), "TotalOrder")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), "TotalPPn")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), "TotalPayment")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), "LastUpdate")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("K%d", no), "LastUpdateBy")
	urut := 0

	// sort.SliceStable(reportPayements, func(i, j int) bool {
	// 	return reportPayements[i].PaymentDate > reportPayements[j].PaymentDate
	// })

	for _, rs := range reportPayements {
		no++
		urut++
		xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), urut)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), rs.PaymentTypeName)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), rs.PaymentNo)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), rs.PaymentDate1)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), rs.SalesOrderNo)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), rs.OrderDate)
		ppnRp := int64(0)
		sebelumPPn := rs.TotalOrder
		if rs.TotalPpn > 0 {
			fmt.Println("hitung PPN ", rs.PaymentNo)
			fmt.Println("total pay ", rs.TotalPayment)
			fmt.Println("total ppn ", rs.TotalPpn)
			ppnPersen := (100 + (float32(rs.TotalPpn))) / 100
			fmt.Println("total ppn persen ", ppnPersen)
			//totalord := float32(rs.TotalPayment) / (ppnPersen)
			totalord := float32(rs.TotalOrder) / (ppnPersen)
			fmt.Println("total ord ", totalord)
			sebelumPPn = int64(totalord)
			fmt.Println("total sblm ppn ", sebelumPPn)
			ppnRp = sebelumPPn * rs.TotalPpn / 100
			fmt.Println("total  ppn ", ppnRp)
		}
		xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), sebelumPPn)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), ppnRp)
		// xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), rs.TotalPayment)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), rs.TotalOrder)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), rs.LastUpdate)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("K%d", no), rs.LastUpdateBy)
	}

	if err := xls.SaveAs(filename); err != nil {
		success = false
		return
	}
	success = true
	return
}

func ExportToExcelReportPaymentSupplier(reportPayements []dto.ReportPaymentSupplier, dateStart, dateEnd, namaFile string) (filename string, success bool) {

	filename = fmt.Sprintf("%v_%v_%v.csv", namaFile, dateStart, dateEnd)

	sheet1Name := "Sheet1"
	xls := excelize.NewFile()
	xls.NewSheet(sheet1Name)

	t1, _ := time.Parse("2006-01-02", dateStart)
	t2, _ := time.Parse("2006-01-02", dateEnd)
	fmt.Println("Waktu nya adalah", t1.Format("02-January-2006"))
	no := 1
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "REPORT PAYMENT")

	no++
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), fmt.Sprintf("%v - %v ", t1.Format("02-January-2006"), t2.Format("02-January-2006")))

	no = no + 2
	xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), "#")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), "PaymentNo")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), "PaymentDate")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), "ReceiveNo")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), "ReceiveDate")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), "Supplier")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), "PaymentReff")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), "Status")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), "PaymentMethod")
	xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), "Total")
	urut := 0

	for _, rs := range reportPayements {

		// fmt.Println(rs.ReceiveNo, "-", rs.ReceiveTgl)
		no++
		urut++
		xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", no), urut)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", no), rs.PaymentNo)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", no), rs.PaymentDate)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", no), rs.ReceiveNo)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", no), rs.ReceiveTgl)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", no), rs.Supplier)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", no), rs.PaymentReff)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", no), rs.Status)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", no), rs.PaymentType)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", no), rs.GrandTotal)
	}

	if err := xls.SaveAs(filename); err != nil {
		success = false
		return
	}
	success = true
	return
}
