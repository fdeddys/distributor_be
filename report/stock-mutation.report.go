package report

import (
	"distribution-system-be/database"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/signintech/gopdf"
)

var (
	stockMutationNumber string
)

func GenerateStockMutationReport(stockMutationID int64) {

	title = "Stock Mutation"

	spaceLen = beego.AppConfig.DefaultFloat("report.space-len", 15)
	pageMargin = beego.AppConfig.DefaultFloat("report.page-margin", 12)

	curPage = 1

	spaceCustomerInfo = 300
	spaceTitik = spaceCustomerInfo + 150
	spaceValue = spaceCustomerInfo + 160

	spaceSummaryInfo = spaceCustomerInfo
	spaceTitikSumamry = spaceTitik
	spaceValueSummary = spaceValue

	tblCol1 = 25
	tblCol2 = 80
	tblCol3 = 300
	tblCol4 = 370
	tblCol5 = 430
	tblCol6 = 500

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.SetMargins(pageMargin, pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	if err := pdf.AddTTFFont("open-sans", "font/OpenSans-Regular.ttf"); err != nil {
		log.Print(err.Error())
		return
	}

	if err := pdf.AddTTFFont("open-sans-bold", "font/OpenSans-Bold.ttf"); err != nil {
		log.Print(err.Error())
		return
	}

	// untuk nomor urut barang
	number = 1

	// get Data mockup utk display ke grid
	fmt.Println("data  send to fillData Details : ", stockMutationID)
	dataDetails := fillDataDetailStockMutation(stockMutationID)

	fmt.Println("hasil fill")
	for i, detail := range dataDetails {
		fmt.Println(i, "====", detail)
	}
	fmt.Println("=============")
	// setFont(&pdf, 12)
	setHeader(&pdf, "sm")
	pdf.Br(20)

	setDetail(&pdf, dataDetails, "mt")
	setSummary(&pdf)
	setSign(&pdf, "Requestor", "Approver", "")

	pdf.WritePdf("stoc-mutation .pdf")

}

func fillDataDetailStockMutation(stockMutationID int64) []DataDetail {

	stockMutation, err := database.GetStockMutationById(stockMutationID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stockMutation)

	stockMutationNumber = stockMutation.StockMutationNo

	stockMutationDetails := database.GetAllDataStockMutationDetail(stockMutationID)

	fmt.Println("orderDetails : ", stockMutationDetails)

	go fillDataWarehouse(
		stockMutation.WarehouseSource.Name,
		stockMutation.WarehouseDest.Name,
		stockMutation.StockMutationDate.Format("02-01-2006"),
		stockMutation.StockMutationNo,
	)
	// tdk blh kosong
	// per halaman max 25 item detail
	totalRec = len(stockMutationDetails)
	res := make([]DataDetail, totalRec+1)
	var data DataDetail

	subTotal = 0
	tax = 0
	grandTotal = 0
	for i, detail := range stockMutationDetails {
		data.Item = detail.Product.Name
		data.Quantity = int64(detail.Qty)
		data.Unit = detail.UOM.Name
		data.Price = int64(detail.Hpp)
		total := data.Price * data.Quantity
		data.Total = int64(detail.Hpp) * int64(detail.Qty)
		subTotal += total
		res[i+1] = data
		fmt.Println("total sub total", subTotal)
	}
	totalRec = len(res)
	fmt.Println("Jumlah record [fill] =>", totalRec)

	tax = subTotal / 10
	grandTotal = subTotal + tax

	return res
}

func fillDataWarehouse(WhDest, WhSource, MutationDate, MutationNo string) {
	invInfo.CustCode = WhSource
	invInfo.CustName = WhDest
	invInfo.TransAt = MutationDate
	invInfo.SourceDoc = MutationNo
}
